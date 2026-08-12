package cache

import (
	"os"
	"testing"
	"time"
)

func TestKey(t *testing.T) {
	// Deterministic md5-based key with a stable prefix.
	if got := Key("opendict", "hello"); got != "opendict_5d41402abc4b2a76b9719d911017c592" {
		t.Fatalf("Key = %q", got)
	}
	// A word containing '/' still yields a filesystem-safe key.
	if got := Key("opendict", "a/b"); len(got) != len("opendict_")+32 {
		t.Fatalf("unexpected key length for slash word: %q", got)
	}
}

func TestReadWriteFreshness(t *testing.T) {
	t.Setenv("alfred_workflow_cache", t.TempDir())

	if err := Write("k", []byte(`{"v":1}`)); err != nil {
		t.Fatal(err)
	}

	// Within maxAge -> hit.
	if data, ok := Read("k", time.Hour); !ok || string(data) != `{"v":1}` {
		t.Fatalf("fresh read failed: %q ok=%v", data, ok)
	}

	// maxAge 0 -> never expires, always a hit when present.
	if _, ok := Read("k", 0); !ok {
		t.Fatal("maxAge=0 should always hit when file exists")
	}

	// Stale (maxAge shorter than age) -> miss.
	old := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(path("k"), old, old)
	if _, ok := Read("k", time.Hour); ok {
		t.Fatal("stale entry should miss")
	}

	// Missing key -> miss.
	if _, ok := Read("nope", time.Hour); ok {
		t.Fatal("missing key should miss")
	}
}

func TestCachedUsesLoaderOnce(t *testing.T) {
	t.Setenv("alfred_workflow_cache", t.TempDir())

	calls := 0
	loader := func() ([]byte, error) {
		calls++
		return []byte("payload"), nil
	}

	for i := 0; i < 3; i++ {
		data, err := Cached("k", time.Hour, loader)
		if err != nil || string(data) != "payload" {
			t.Fatalf("Cached returned %q err=%v", data, err)
		}
	}
	if calls != 1 {
		t.Fatalf("loader called %d times, want 1", calls)
	}

	// Nil loader on a miss returns no data without error.
	if data, err := Cached("missing", time.Hour, nil); err != nil || data != nil {
		t.Fatalf("nil loader miss = %q err=%v", data, err)
	}
}

func TestCachedPrunesStaleSiblings(t *testing.T) {
	t.Setenv("alfred_workflow_cache", t.TempDir())

	load := func(v string) func() ([]byte, error) {
		return func() ([]byte, error) { return []byte(v), nil }
	}

	// Two distinct search words populate two per-word files.
	oldKey := Key("opendict", "old")
	if _, err := Cached(oldKey, time.Hour, load("a")); err != nil {
		t.Fatal(err)
	}
	// Age the first entry beyond the freshness window.
	past := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(path(oldKey), past, past); err != nil {
		t.Fatal(err)
	}

	// A different word must not be able to serve the stale entry from a hit
	// (different key), and writing it must evict the stale sibling.
	newKey := Key("opendict", "new")
	if _, err := Cached(newKey, time.Hour, load("b")); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path(oldKey)); !os.IsNotExist(err) {
		t.Fatalf("stale sibling should have been pruned, stat err=%v", err)
	}
	if _, err := os.Stat(path(newKey)); err != nil {
		t.Fatalf("fresh entry should remain: %v", err)
	}
}

func TestCachedPruneSkipsOtherPrefixes(t *testing.T) {
	t.Setenv("alfred_workflow_cache", t.TempDir())

	// A stale entry under a different prefix must survive pruning.
	other := Key("suggest", "keep")
	if err := Write(other, []byte("keep")); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(path(other), past, past); err != nil {
		t.Fatal(err)
	}

	load := func() ([]byte, error) { return []byte("x"), nil }
	if _, err := Cached(Key("opendict", "word"), time.Hour, load); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path(other)); err != nil {
		t.Fatalf("entry under a different prefix must not be pruned: %v", err)
	}
}
