// Package cache is a small on-disk cache replacing wf.cached_data / cache_data.
// Entries are JSON files under Alfred's per-workflow cache directory, with
// freshness judged by file mtime.
package cache

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// dir returns Alfred's workflow cache directory, or a temp fallback when the
// workflow env var is absent (e.g. local runs and tests).
func dir() string {
	if d := os.Getenv("alfred_workflow_cache"); d != "" {
		return d
	}
	return filepath.Join(os.TempDir(), "alfkoreansearch-cache")
}

func path(key string) string {
	return filepath.Join(dir(), key+".json")
}

// Key builds a filesystem-safe cache key: prefix + md5(word). The word may
// contain '/' or other path separators, so it is always hashed.
func Key(prefix, word string) string {
	sum := md5.Sum([]byte(word))
	return prefix + "_" + hex.EncodeToString(sum[:])
}

// Read returns the cached bytes for key when present and still fresh. A maxAge
// of 0 means "never expires": any existing entry is returned.
func Read(key string, maxAge time.Duration) ([]byte, bool) {
	p := path(key)
	fi, err := os.Stat(p)
	if err != nil {
		return nil, false
	}
	if maxAge > 0 && time.Since(fi.ModTime()) > maxAge {
		return nil, false
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, false
	}
	return data, true
}

// Write stores data under key, creating the cache directory as needed.
func Write(key string, data []byte) error {
	if err := os.MkdirAll(dir(), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path(key), data, 0o644)
}

// Cached returns a fresh cached entry for key, or invokes loader, stores its
// result and returns it. When loader is nil and the entry is missing/stale,
// it returns (nil, nil) — mirroring cached_data(key, None, ...).
func Cached(key string, maxAge time.Duration, loader func() ([]byte, error)) ([]byte, error) {
	if data, ok := Read(key, maxAge); ok {
		return data, nil
	}
	if loader == nil {
		return nil, nil
	}
	data, err := loader()
	if err != nil {
		return nil, err
	}
	if err := Write(key, data); err != nil {
		return nil, err
	}
	// Search-style keys (prefix + md5(word)) create one file per distinct
	// word, so the cache would grow without bound. Evict same-prefix siblings
	// that are already past maxAge — they can no longer produce a hit and only
	// accumulate on disk. maxAge 0 ("never expires") entries are left
	// untouched, and the sweep itself runs at most once a day (see maybePrune)
	// so it does not add a directory scan to every cache write.
	if maxAge > 0 {
		maybePrune(key, maxAge)
	}
	return data, nil
}

// pruneMarker is the fixed key whose mtime records when the last prune sweep
// ran. Its name shares no search prefix, so pruneStale never matches it.
const pruneMarker = "__prune_check"

// pruneInterval throttles pruneStale: the sweep runs at most once per interval.
const pruneInterval = 24 * time.Hour

// maybePrune runs pruneStale at most once per pruneInterval, gated by the
// mtime of a marker file, so eviction happens roughly daily rather than on
// every cache write.
func maybePrune(key string, maxAge time.Duration) {
	if _, fresh := Read(pruneMarker, pruneInterval); fresh {
		return
	}
	_ = Write(pruneMarker, []byte(time.Now().Format(time.RFC3339)))
	pruneStale(key, maxAge)
}

// pruneStale removes cache files sharing key's prefix whose mtime is older
// than maxAge. The prefix is the part of key before its final '_' (see Key),
// so pruning stays scoped to one logical cache and never removes unrelated
// fixed-key entries. Errors are ignored: pruning is best-effort housekeeping.
func pruneStale(key string, maxAge time.Duration) {
	i := strings.LastIndex(key, "_")
	if i <= 0 {
		return
	}
	matches, err := filepath.Glob(filepath.Join(dir(), key[:i]+"_*.json"))
	if err != nil {
		return
	}
	for _, p := range matches {
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		if time.Since(fi.ModTime()) > maxAge {
			_ = os.Remove(p)
		}
	}
}
