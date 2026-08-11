package main

import "testing"

func TestNormalizeArgs(t *testing.T) {
	// Build NFD (decomposed) inputs explicitly via code points, independent of
	// this file's own encoding.
	//   "사랑" as conjoining jamo: ᄉ ᅡ ᄅ ᅡ ᆼ
	loveNFD := string([]rune{0x1109, 0x1161, 0x1105, 0x1161, 0x11BC})
	//   "がっこう": か + combining dakuten, then っこう
	gakkoNFD := string([]rune{0x304B, 0x3099, 0x3063, 0x3053, 0x3046})

	args := []string{"search", loveNFD, gakkoNFD, "ascii"}
	normalizeArgs(args)

	if want := "사랑"; args[1] != want {
		t.Errorf("Hangul not composed to NFC: got % x, want % x", []byte(args[1]), []byte(want))
	}
	if want := "がっこう"; args[2] != want {
		t.Errorf("Japanese dakuten not composed to NFC: got % x, want % x", []byte(args[2]), []byte(want))
	}
	// Handler names and ASCII are untouched.
	if args[0] != "search" || args[3] != "ascii" {
		t.Errorf("ASCII args changed: %v", args)
	}
}
