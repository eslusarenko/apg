package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// charClass is a set of symbol classes, as in apg's -M mode.
type charClass uint8

const (
	classNumber charClass = 1 << iota
	classSpecial
	classCapital
	classSmall
)

var classSymbols = map[charClass]string{
	classSmall:   "abcdefghijklmnopqrstuvwxyz",
	classCapital: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	classNumber:  "0123456789",
	classSpecial: "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
}

// classOrder fixes the order classes are listed in messages.
var classOrder = []charClass{classSpecial, classNumber, classCapital, classSmall}

var classNames = map[charClass]string{
	classSpecial: "special symbols",
	classNumber:  "numbers",
	classCapital: "capital letters",
	classSmall:   "small letters",
}

// mode is a parsed -M value: classes to use and classes every password
// must contain.
type mode struct {
	use, must charClass
}

// parseMode parses an apg -M string. S, N, C and L mean the password must
// contain special symbols, numbers, capital or small letters; the lowercase
// letters mean it should (the class is used but not required).
func parseMode(s string) (mode, error) {
	if s == "" {
		return mode{}, fmt.Errorf("mode must not be empty")
	}
	var m mode
	for _, r := range s {
		var c charClass
		switch r {
		case 'S', 's':
			c = classSpecial
		case 'N', 'n':
			c = classNumber
		case 'C', 'c':
			c = classCapital
		case 'L', 'l':
			c = classSmall
		default:
			return mode{}, fmt.Errorf("invalid mode %q: use the letters S, N, C and L (upper case: must, lower case: should)", s)
		}
		m.use |= c
		if r >= 'A' && r <= 'Z' {
			m.must |= c
		}
	}
	return m, nil
}

// charset is the set of symbols a password may be built from.
type charset struct {
	mode    mode
	exclude string
	symbols map[charClass]string // allowed symbols of each used class
}

func newCharset(m mode, exclude string) *charset {
	cs := &charset{mode: m, exclude: exclude, symbols: map[charClass]string{}}
	for c, all := range classSymbols {
		if m.use&c == 0 {
			continue
		}
		var b strings.Builder
		for i := 0; i < len(all); i++ {
			if !strings.ContainsRune(exclude, rune(all[i])) {
				b.WriteByte(all[i])
			}
		}
		cs.symbols[c] = b.String()
	}
	return cs
}

func (cs *charset) uses(c charClass) bool { return cs.mode.use&c != 0 }

func (cs *charset) allowed(b byte) bool { return !strings.ContainsRune(cs.exclude, rune(b)) }

// alphabet returns every allowed symbol of the used classes.
func (cs *charset) alphabet() string {
	var b strings.Builder
	for _, c := range []charClass{classSmall, classCapital, classNumber, classSpecial} {
		b.WriteString(cs.symbols[c])
	}
	return b.String()
}

// randomOf returns a random allowed symbol of class c.
func (cs *charset) randomOf(c charClass) (byte, bool) {
	s := cs.symbols[c]
	if s == "" {
		return 0, false
	}
	return s[randInt(len(s))], true
}

// satisfies reports whether p contains a symbol of every required class.
func (cs *charset) satisfies(p string) bool {
	for _, c := range classOrder {
		if cs.mode.must&c != 0 && !strings.ContainsAny(p, classSymbols[c]) {
			return false
		}
	}
	return true
}

// randInt returns a uniform random integer in [0, n).
func randInt(n int) int {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		// crypto/rand.Reader never fails since Go 1.24.
		panic(err)
	}
	return int(v.Int64())
}
