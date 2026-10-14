package main

import (
	"errors"
	"fmt"
)

type algorithm int

const (
	algPronounceable algorithm = 0
	algRandom        algorithm = 1
)

// maxAttempts bounds how many passwords may be rejected (for a missing
// required symbol class or a dictionary match) before giving up.
const maxAttempts = 100000

type options struct {
	algorithm      algorithm
	minLen, maxLen int
	mode           mode
	exclude        string
	dict           map[string]struct{}
}

type password struct {
	text       string
	hyphenated string // pronunciation; only for pronounceable passwords
}

// generator produces passwords for a fixed set of options.
type generator struct {
	opts options
	cs   *charset
}

func newGenerator(o options) (*generator, error) {
	if o.algorithm != algPronounceable && o.algorithm != algRandom {
		return nil, fmt.Errorf("unknown algorithm %d (want 0 or 1)", o.algorithm)
	}
	if o.minLen < 1 || o.maxLen < o.minLen {
		return nil, errors.New("password length must be at least 1")
	}

	cs := newCharset(o.mode, o.exclude)
	must := 0
	for _, c := range classOrder {
		if o.mode.must&c == 0 {
			continue
		}
		must++
		// Small letters are always available to pronounceable passwords.
		if cs.symbols[c] == "" && !(o.algorithm == algPronounceable && c == classSmall) {
			return nil, fmt.Errorf("mode requires %s, but all of them are excluded", classNames[c])
		}
	}
	if o.algorithm == algRandom {
		if cs.alphabet() == "" {
			return nil, errors.New("no characters left to build a password from")
		}
		if o.maxLen < must {
			return nil, fmt.Errorf("length must be at least %d to include every required symbol class", must)
		}
	}
	return &generator{opts: o, cs: cs}, nil
}

// next returns a new password that satisfies the mode and is not in the
// dictionary.
func (g *generator) next() (password, error) {
	for range maxAttempts {
		length := g.opts.minLen + randInt(g.opts.maxLen-g.opts.minLen+1)
		var p password
		if g.opts.algorithm == algPronounceable {
			p.text, p.hyphenated = pronounceable(length, g.cs)
		} else {
			p.text = g.random(length)
		}
		if !g.cs.satisfies(p.text) {
			continue
		}
		if _, found := g.opts.dict[p.text]; found {
			continue
		}
		return p, nil
	}
	return password{}, errors.New("could not generate a password that meets the requirements; try a longer length or a less strict -M mode")
}

// random returns length symbols chosen uniformly from the charset.
func (g *generator) random(length int) string {
	alphabet := g.cs.alphabet()
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet[randInt(len(alphabet))]
	}
	return string(b)
}
