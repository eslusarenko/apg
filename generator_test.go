package main

import (
	"strings"
	"testing"
)

func mustMode(t *testing.T, s string) mode {
	t.Helper()
	m, err := parseMode(s)
	if err != nil {
		t.Fatalf("parseMode(%q): %v", s, err)
	}
	return m
}

func mustGenerator(t *testing.T, o options) *generator {
	t.Helper()
	g, err := newGenerator(o)
	if err != nil {
		t.Fatalf("newGenerator(%+v): %v", o, err)
	}
	return g
}

func next(t *testing.T, g *generator) password {
	t.Helper()
	p, err := g.next()
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseMode(t *testing.T) {
	m := mustMode(t, "SNcl")
	if m.use != classSpecial|classNumber|classCapital|classSmall {
		t.Errorf("use = %b", m.use)
	}
	if m.must != classSpecial|classNumber {
		t.Errorf("must = %b", m.must)
	}
	for _, s := range []string{"", "X", "SNCLR", "sn c"} {
		if _, err := parseMode(s); err == nil {
			t.Errorf("parseMode(%q): expected error", s)
		}
	}
}

func TestRandomLength(t *testing.T) {
	g := mustGenerator(t, options{algorithm: algRandom, minLen: 8, maxLen: 12, mode: mustMode(t, "SNCL")})
	seen := map[int]bool{}
	for range 2000 {
		n := len(next(t, g).text)
		if n < 8 || n > 12 {
			t.Fatalf("length %d outside 8..12", n)
		}
		seen[n] = true
	}
	if len(seen) != 5 {
		t.Errorf("lengths seen: %v", seen)
	}
}

func TestRandomMustClasses(t *testing.T) {
	g := mustGenerator(t, options{algorithm: algRandom, minLen: 4, maxLen: 4, mode: mustMode(t, "SNCL")})
	for range 1000 {
		p := next(t, g).text
		for _, c := range classOrder {
			if !strings.ContainsAny(p, classSymbols[c]) {
				t.Fatalf("%q has no %s", p, classNames[c])
			}
		}
	}
}

func TestRandomOnlyUsedClasses(t *testing.T) {
	cases := map[string]string{
		"NCL": classSymbols[classNumber] + classSymbols[classCapital] + classSymbols[classSmall],
		"N":   classSymbols[classNumber],
		"cl":  classSymbols[classCapital] + classSymbols[classSmall],
	}
	for ms, allowed := range cases {
		g := mustGenerator(t, options{algorithm: algRandom, minLen: 32, maxLen: 32, mode: mustMode(t, ms)})
		for range 200 {
			p := next(t, g).text
			if strings.Trim(p, allowed) != "" {
				t.Fatalf("-M %s: %q has unexpected characters", ms, p)
			}
		}
	}
}

func TestRandomExclude(t *testing.T) {
	g := mustGenerator(t, options{algorithm: algRandom, minLen: 64, maxLen: 64, mode: mustMode(t, "N"), exclude: "0123"})
	for range 200 {
		if p := next(t, g).text; strings.ContainsAny(p, "0123") {
			t.Fatalf("%q contains excluded characters", p)
		}
	}
}

func TestDictionaryRejects(t *testing.T) {
	dict := map[string]struct{}{}
	for _, d := range "012345678" {
		dict[string(d)] = struct{}{}
	}
	g := mustGenerator(t, options{algorithm: algRandom, minLen: 1, maxLen: 1, mode: mustMode(t, "N"), dict: dict})
	for range 100 {
		if p := next(t, g).text; p != "9" {
			t.Fatalf("got %q, want 9", p)
		}
	}

	dict["9"] = struct{}{}
	g = mustGenerator(t, options{algorithm: algRandom, minLen: 1, maxLen: 1, mode: mustMode(t, "N"), dict: dict})
	if _, err := g.next(); err == nil {
		t.Error("expected an error when every password is in the dictionary")
	}
}

func TestNewGeneratorErrors(t *testing.T) {
	cases := map[string]options{
		"bad algorithm":   {algorithm: 2, minLen: 8, maxLen: 8, mode: mustMode(t, "SNCL")},
		"zero length":     {algorithm: algRandom, minLen: 0, maxLen: 0, mode: mustMode(t, "SNCL")},
		"too short":       {algorithm: algRandom, minLen: 3, maxLen: 3, mode: mustMode(t, "SNCL")},
		"must excluded":   {algorithm: algRandom, minLen: 8, maxLen: 8, mode: mustMode(t, "NL"), exclude: "0123456789"},
		"nothing left":    {algorithm: algRandom, minLen: 8, maxLen: 8, mode: mustMode(t, "n"), exclude: "0123456789"},
		"pron no capital": {algorithm: algPronounceable, minLen: 8, maxLen: 8, mode: mustMode(t, "Cl"), exclude: classSymbols[classCapital]},
	}
	for name, o := range cases {
		if _, err := newGenerator(o); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestRandomIsRandom(t *testing.T) {
	g := mustGenerator(t, options{algorithm: algRandom, minLen: 32, maxLen: 32, mode: mustMode(t, "SNCL")})
	seen := map[string]bool{}
	for range 100 {
		p := next(t, g).text
		if seen[p] {
			t.Fatalf("duplicate password %q", p)
		}
		seen[p] = true
	}
}

func TestStrategies(t *testing.T) {
	want := map[string]string{"full": "SNCL", "alphanum": "NCL", "digits": "N"}
	for name, ms := range want {
		st, err := parseStrategy(name)
		if err != nil {
			t.Fatalf("parseStrategy(%q): %v", name, err)
		}
		if st.modeString() != ms {
			t.Errorf("%s: mode %q, want %q", name, st.modeString(), ms)
		}
	}
	for _, name := range []string{"", "alnum", "Digits", "pin"} {
		if _, err := parseStrategy(name); err == nil {
			t.Errorf("parseStrategy(%q): expected error", name)
		}
	}
}

func TestSpell(t *testing.T) {
	if got, want := spell("aB1!"), "alfa-Bravo-ONE-EXCLAMATION_POINT"; got != want {
		t.Errorf("spell = %q, want %q", got, want)
	}
	for _, s := range classSymbols {
		for i := 0; i < len(s); i++ {
			if spell(s[i:i+1]) == "" {
				t.Errorf("no spelling for %q", s[i])
			}
		}
	}
}
