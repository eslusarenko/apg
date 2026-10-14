package main

import (
	"strings"
	"testing"
)

func TestGenerateLength(t *testing.T) {
	for _, n := range []int{4, 16, 32, 128} {
		p, err := Generate(n, true)
		if err != nil {
			t.Fatalf("Generate(%d): %v", n, err)
		}
		if len(p) != n {
			t.Errorf("Generate(%d) returned %d chars", n, len(p))
		}
	}
}

func TestGenerateContainsAllGroups(t *testing.T) {
	for range 1000 {
		p, err := Generate(4, true)
		if err != nil {
			t.Fatal(err)
		}
		for _, g := range []string{lowerChars, upperChars, digitChars, specialChars} {
			if !strings.ContainsAny(p, g) {
				t.Fatalf("%q has no characters from %q", p, g)
			}
		}
	}
}

func TestGenerateNoSpecial(t *testing.T) {
	allowed := lowerChars + upperChars + digitChars
	for range 1000 {
		p, err := Generate(32, false)
		if err != nil {
			t.Fatal(err)
		}
		if strings.ContainsAny(p, specialChars) {
			t.Fatalf("%q contains special characters", p)
		}
		for _, c := range p {
			if !strings.ContainsRune(allowed, c) {
				t.Fatalf("%q contains unexpected %q", p, c)
			}
		}
	}
}

func TestGenerateTooShort(t *testing.T) {
	cases := []struct {
		length     int
		useSpecial bool
	}{
		{3, true},
		{2, false},
		{0, true},
		{-1, false},
	}
	for _, c := range cases {
		if _, err := Generate(c.length, c.useSpecial); err == nil {
			t.Errorf("Generate(%d, %v): expected error", c.length, c.useSpecial)
		}
	}
}

func TestGenerateIsRandom(t *testing.T) {
	seen := make(map[string]bool)
	for range 100 {
		p, err := Generate(32, true)
		if err != nil {
			t.Fatal(err)
		}
		if seen[p] {
			t.Fatalf("duplicate password %q", p)
		}
		seen[p] = true
	}
}
