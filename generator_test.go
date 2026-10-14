package main

import (
	"strings"
	"testing"
)

func TestGenerateLength(t *testing.T) {
	for _, n := range []int{4, 16, 32, 128} {
		p, err := Generate(n, strategyFull)
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
		p, err := Generate(4, strategyFull)
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
		p, err := Generate(32, strategyAlphanum)
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

func TestGenerateDigits(t *testing.T) {
	for _, n := range []int{1, 4, 6, 32} {
		p, err := Generate(n, strategyDigits)
		if err != nil {
			t.Fatalf("Generate(%d, digits): %v", n, err)
		}
		if len(p) != n {
			t.Errorf("Generate(%d, digits) returned %d chars", n, len(p))
		}
		if strings.Trim(p, digitChars) != "" {
			t.Errorf("%q contains non-digits", p)
		}
	}
}

func TestParseStrategy(t *testing.T) {
	for _, name := range []string{"full", "alphanum", "digits"} {
		if s, err := parseStrategy(name); err != nil || s != strategy(name) {
			t.Errorf("parseStrategy(%q) = %q, %v", name, s, err)
		}
	}
	for _, name := range []string{"", "alnum", "Digits", "pin"} {
		if _, err := parseStrategy(name); err == nil {
			t.Errorf("parseStrategy(%q): expected error", name)
		}
	}
}

func TestGenerateTooShort(t *testing.T) {
	cases := []struct {
		length int
		s      strategy
	}{
		{3, strategyFull},
		{2, strategyAlphanum},
		{0, strategyDigits},
		{-1, strategyFull},
	}
	for _, c := range cases {
		if _, err := Generate(c.length, c.s); err == nil {
			t.Errorf("Generate(%d, %s): expected error", c.length, c.s)
		}
	}
}

func TestGenerateIsRandom(t *testing.T) {
	seen := make(map[string]bool)
	for range 100 {
		p, err := Generate(32, strategyFull)
		if err != nil {
			t.Fatal(err)
		}
		if seen[p] {
			t.Fatalf("duplicate password %q", p)
		}
		seen[p] = true
	}
}
