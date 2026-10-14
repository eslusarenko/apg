package main

import (
	"strings"
	"testing"
	"time"
)

func TestPronounceableLengths(t *testing.T) {
	cs := newCharset(mode{use: classSmall}, "")
	deadline := time.Now().Add(20 * time.Second)
	for length := 1; length <= 64; length++ {
		for range 50 {
			if time.Now().After(deadline) {
				t.Fatal("pronounceable generation is too slow")
			}
			w, h := pronounceable(length, cs)
			if len(w) != length {
				t.Fatalf("length %d: got %q (%d chars)", length, w, len(w))
			}
			if strings.Trim(w, classSymbols[classSmall]) != "" {
				t.Fatalf("%q contains more than small letters", w)
			}
			// Without substitutions the pronunciation is the word split
			// into syllables.
			if strings.ReplaceAll(h, "-", "") != w {
				t.Fatalf("pronunciation %q does not match %q", h, w)
			}
			for _, syl := range strings.Split(h, "-") {
				if syl == "" {
					t.Fatalf("empty syllable in %q", h)
				}
			}
		}
	}
}

// Every syllable but the last has a vowel. Like the original apg, the
// algorithm occasionally ends a word with a consonant-only syllable
// (about 0.2% of 10-letter words).
func TestPronounceableHasVowels(t *testing.T) {
	cs := newCharset(mode{use: classSmall}, "")
	for range 2000 {
		_, h := pronounceable(10, cs)
		syls := strings.Split(h, "-")
		for _, syl := range syls[:len(syls)-1] {
			if !strings.ContainsAny(syl, "aeiouy") {
				t.Fatalf("syllable %q of %q has no vowel", syl, h)
			}
		}
	}
}

func TestPronounceableModes(t *testing.T) {
	g := mustGenerator(t, options{algorithm: algPronounceable, minLen: 8, maxLen: 10, mode: mustMode(t, "SNCL")})
	for range 200 {
		p := next(t, g)
		if n := len(p.text); n < 8 || n > 10 {
			t.Fatalf("%q: length %d outside 8..10", p.text, n)
		}
		for _, c := range classOrder {
			if !strings.ContainsAny(p.text, classSymbols[c]) {
				t.Fatalf("%q has no %s", p.text, classNames[c])
			}
		}
	}
}

func TestPronounceableExclude(t *testing.T) {
	g := mustGenerator(t, options{algorithm: algPronounceable, minLen: 12, maxLen: 12, mode: mustMode(t, "sncl"), exclude: "ABC0123!@#"})
	for range 500 {
		if p := next(t, g).text; strings.ContainsAny(p, "ABC0123!@#") {
			t.Fatalf("%q contains excluded characters", p)
		}
	}
}
