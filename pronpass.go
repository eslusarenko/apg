// Pronounceable password generation (FIPS-181), ported to Go from
// pronpass.c of apg (Automated Password Generator). apg's module in turn
// uses code from the NIST implementation of FIPS-181 and the algorithm
// originally worked out by Morrie Gasser in 1975.
//
// Original copyright notice and license:
//
// Copyright (c) 1999, 2000, 2001, 2002, 2003
// Adel I. Mirzazhanov. All rights reserved
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//     1.Redistributions of source code must retain the above copyright notice,
//       this list of conditions and the following disclaimer.
//     2.Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     3.The name of the author may not be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR  ``AS IS'' AND ANY EXPRESS
// OR IMPLIED WARRANTIES, INCLUDING,  BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN  NO  EVENT  SHALL THE AUTHOR BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO,  PROCUREMENT OF SUBSTITUTE
// GOODS OR SERVICES;  LOSS OF USE,  DATA,  OR  PROFITS;  OR BUSINESS
// INTERRUPTION)  HOWEVER  CAUSED  AND  ON  ANY  THEORY OF LIABILITY,
// WHETHER  IN  CONTRACT,   STRICT   LIABILITY,  OR  TORT  (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import "strings"

// numUnits is the number of units (letters and letter pairs) words are
// built from; see rules in pronpass_tables.go.
const numUnits = 34

// Unit flags.
const (
	notBeginSyllable = 0o10
	noFinalSplit     = 0o4
	vowel            = 0o2
	alternateVowel   = 0o1
	noSpecialRule    = 0
)

// Digram (unit pair) flags.
const (
	dBegin       = 0o200
	dNotBegin    = 0o100
	dBreak       = 0o40
	dPrefix      = 0o20
	dIllegalPair = 0o10
	dSuffix      = 0o4
	dEnd         = 0o2
	dNotEnd      = 0o1
	dAny         = 0
)

type unitRule struct {
	code  string
	flags uint8
}

func isVowel(u int) bool             { return rules[u].flags&vowel != 0 }
func isAltVowel(u int) bool          { return rules[u].flags&alternateVowel != 0 }
func isFinalE(u int) bool            { return rules[u].flags&noFinalSplit != 0 }
func pairHas(a, b int, f uint8) bool { return digram[a][b]&f != 0 }

// pronGen holds the units carried over from one syllable to the next.
type pronGen struct {
	savedUnit int
	savedPair [2]int
}

// pronounceable returns a pronounceable word of exactly length characters
// and its hyphenated form (syllables separated by "-"). Characters of the
// classes in use may replace or capitalize letters as in apg.
func pronounceable(length int, cs *charset) (word, hyphenated string) {
	if length <= 0 {
		return "", ""
	}
	g := &pronGen{}
	maxRetries := 4*length + numUnits
	syllableUnits := make([]int, length+1)

	for {
		var w, h strings.Builder
		var wordUnits []int
		wordLength := 0

		for tries := 0; wordLength < length && tries <= maxRetries; tries++ {
			syl, last := g.syllable(length-wordLength, syllableUnits)
			units := syllableUnits[:last+1]
			candidate := append(wordUnits, units...)

			if improperWord(candidate) ||
				(wordLength == 0 && haveInitialY(units)) ||
				(wordLength+len(syl) == length && haveFinalSplit(units)) {
				continue
			}
			wordUnits = candidate

			syl = decorateSyllable(syl, cs)
			if wordLength > 0 {
				h.WriteByte('-')
			}
			w.WriteString(syl)
			if len(syl) == 1 {
				h.WriteString(symbolName(syl[0]))
			} else {
				h.WriteString(syl)
			}
			wordLength += len(syl)
		}
		if wordLength == length {
			return w.String(), h.String()
		}
		// Too many rejected syllables: start the word from scratch.
	}
}

// decorateSyllable applies apg's numeric, special and capital substitutions
// to a syllable that has already passed the word checks.
func decorateSyllable(syl string, cs *charset) string {
	changed := false
	dsd := randInt(2)
	if len(syl) == 1 && dsd == 0 && cs.uses(classNumber) {
		if c, ok := cs.randomOf(classNumber); ok {
			syl, changed = string(c), true
		}
	}
	if len(syl) == 1 && dsd == 1 && cs.uses(classSpecial) {
		if c, ok := cs.randomOf(classSpecial); ok {
			syl, changed = string(c), true
		}
	}
	// Unlike the original apg, whose lookup table skips y, any letter may be
	// capitalized.
	if !changed && cs.uses(classCapital) && randInt(2) == 1 {
		if up := syl[0] - 'a' + 'A'; syl[0] >= 'a' && syl[0] <= 'z' && cs.allowed(up) {
			syl = string(up) + syl[1:]
		}
	}
	return syl
}

// syllable generates the next syllable of at most length characters. It
// fills units with the syllable's unit indexes and returns the syllable and
// the index of its last unit.
func (g *pronGen) syllable(length int, units []int) (string, int) {
	maxRetries := 4*length + numUnits
	holdSavedUnit := g.savedUnit

	var (
		syl            []byte
		unit, lastUnit int
		syllableLength int
	)
	for {
		var (
			ruleBroken      bool
			tries           int
			vowelCount      int
			currentUnit     int
			lengthLeft      = length
			wantAnotherUnit = true
		)
		g.savedUnit = holdSavedUnit
		syl = syl[:0]

		for {
			wantVowel := false
			for {
				if g.savedUnit != 0 {
					// Units saved by the previous syllable are used first.
					if g.savedUnit == 2 {
						units[0] = g.savedPair[1]
						if isVowel(g.savedPair[1]) {
							vowelCount++
						}
						currentUnit++
						syl = append(syl[:0], rules[g.savedPair[1]].code...)
						lengthLeft -= len(syl)
					}
					unit = g.savedPair[0]
					g.savedUnit = 0
				} else if wantVowel {
					unit = int(vowelNumbers[randInt(len(vowelNumbers))])
				} else {
					unit = int(unitNumbers[randInt(len(unitNumbers))])
				}
				lengthLeft -= len(rules[unit].code)

				// Prevent having a word longer than expected.
				ruleBroken = lengthLeft < 0

				if currentUnit == 0 {
					// First unit of a syllable: no digram tests yet.
					if rules[unit].flags&notBeginSyllable != 0 {
						ruleBroken = true
					} else if lengthLeft == 0 {
						// A one-unit last syllable must be a vowel.
						if isVowel(unit) {
							wantAnotherUnit = false
						} else {
							ruleBroken = true
						}
					}
				} else {
					prev := units[currentUnit-1]
					allowed := func(f uint8) bool { return pairHas(prev, unit, f) }

					if allowed(dIllegalPair) ||
						(allowed(dBreak) && vowelCount == 0) ||
						(allowed(dEnd) && vowelCount == 0 && !isVowel(unit)) {
						ruleBroken = true
					}

					if currentUnit == 1 {
						if allowed(dNotBegin) {
							ruleBroken = true
						}
					} else {
						lastUnit = units[currentUnit-1]
						beforeLast := units[currentUnit-2]

						if (currentUnit == 2 && allowed(dBegin) && isAltVowel(units[0])) ||
							(allowed(dNotEnd) && lengthLeft == 0) ||
							(allowed(dBreak) && pairHas(beforeLast, lastUnit, dNotEnd)) ||
							(allowed(dPrefix) && !isVowel(beforeLast)) {
							ruleBroken = true
						}

						if !ruleBroken && isVowel(unit) && (lengthLeft > 0 || !isFinalE(lastUnit)) {
							if vowelCount > 1 && isVowel(lastUnit) {
								// No three consecutive vowels in a syllable.
								ruleBroken = true
							} else if vowelCount != 0 && !isVowel(lastUnit) {
								// Vowel-consonant-vowel: save the vowel for
								// the next syllable if this one can end here.
								if pairHas(beforeLast, lastUnit, dNotEnd) {
									ruleBroken = true
								} else {
									g.savedUnit = 1
									g.savedPair[0] = unit
									wantAnotherUnit = false
								}
							}
						}
					}

					// Decide whether the syllable can end here.
					if !ruleBroken && wantAnotherUnit {
						if (vowelCount != 0 && isFinalE(unit) && lengthLeft == 0 && !isVowel(lastUnit)) ||
							allowed(dEnd) || lengthLeft == 0 {
							wantAnotherUnit = false
						} else if vowelCount != 0 && lengthLeft > 0 {
							if allowed(dBegin) && currentUnit > 1 && !(vowelCount == 1 && isVowel(lastUnit)) {
								g.savedUnit = 2
								g.savedPair[0] = unit
								g.savedPair[1] = lastUnit
								wantAnotherUnit = false
							} else if allowed(dBreak) {
								g.savedUnit = 1
								g.savedPair[0] = unit
								wantAnotherUnit = false
							}
						} else if allowed(dSuffix) {
							wantVowel = true
						}
					}
				}

				tries++
				if ruleBroken {
					lengthLeft += len(rules[unit].code)
				}
				if !ruleBroken || tries > maxRetries {
					break
				}
			}

			if tries <= maxRetries {
				// An initial y of a syllable counts as a consonant.
				if isVowel(unit) && (currentUnit > 0 || !isAltVowel(unit)) {
					vowelCount++
				}
				switch g.savedUnit {
				case 0:
					units[currentUnit] = unit
					syl = append(syl, rules[unit].code...)
				case 1:
					currentUnit--
				case 2:
					syl = syl[:len(syl)-len(rules[lastUnit].code)]
					lengthLeft += len(rules[lastUnit].code)
					currentUnit -= 2
				}
			} else {
				// Too many tries: start another syllable.
				ruleBroken = true
			}

			syllableLength = currentUnit
			currentUnit++
			if tries > maxRetries || !wantAnotherUnit {
				break
			}
		}

		if !ruleBroken && !illegalPlacement(units, syllableLength) {
			return string(syl), syllableLength
		}
	}
}

// improperWord reports illegal combinations that may span syllables: an
// illegal unit pair, three consecutive vowels or three consecutive
// consonants.
func improperWord(units []int) bool {
	for i := range units {
		if i != 0 && pairHas(units[i-1], units[i], dIllegalPair) {
			return true
		}
		if i >= 2 {
			a, b, c := units[i-2], units[i-1], units[i]
			if (isVowel(a) && !isAltVowel(a) && isVowel(b) && isVowel(c)) ||
				(!isVowel(a) && !isVowel(b) && !isVowel(c)) {
				return true
			}
		}
	}
	return false
}

// haveInitialY reports a syllable whose only vowel is an initial y.
func haveInitialY(units []int) bool {
	vowels, normalVowels := 0, 0
	for i, u := range units {
		if isVowel(u) {
			vowels++
			if !isAltVowel(u) || i != 0 {
				normalVowels++
			}
		}
	}
	return vowels <= 1 && normalVowels == 0
}

// haveFinalSplit reports a final syllable whose only vowel is a silent e.
func haveFinalSplit(units []int) bool {
	vowels := 0
	for _, u := range units {
		if isVowel(u) {
			vowels++
		}
	}
	return vowels == 1 && isFinalE(units[len(units)-1])
}

// illegalPlacement checks a syllable (units[0..last]) for consonants
// between vowels and for three consecutive vowels or consonants.
func illegalPlacement(units []int, last int) bool {
	vowels := 0
	for i := 0; i <= last; i++ {
		u := units[i]
		if i >= 1 {
			p := units[i-1]
			if (!isVowel(p) && isVowel(u) && !(isFinalE(u) && i == last) && vowels != 0) ||
				(i >= 2 &&
					((!isVowel(units[i-2]) && !isVowel(p) && !isVowel(u)) ||
						(isVowel(units[i-2]) && !(isAltVowel(units[0]) && i == 2) && isVowel(p) && isVowel(u)))) {
				return true
			}
		}
		if isVowel(u) && !(isAltVowel(units[0]) && i == 0 && last != 0) {
			vowels++
		}
	}
	return false
}
