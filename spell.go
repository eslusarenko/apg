package main

import "strings"

// symbolNames are the names apg uses for digits and special symbols in
// spelled (-l) and hyphenated (-t) output.
var symbolNames = map[byte]string{
	'1': "ONE", '2': "TWO", '3': "THREE", '4': "FOUR", '5': "FIVE",
	'6': "SIX", '7': "SEVEN", '8': "EIGHT", '9': "NINE", '0': "ZERO",
	'!': "EXCLAMATION_POINT", '"': "QUOTATION_MARK", '#': "CROSSHATCH",
	'$': "DOLLAR_SIGN", '%': "PERCENT_SIGN", '&': "AMPERSAND",
	'\'': "APOSTROPHE", '(': "LEFT_PARENTHESIS", ')': "RIGHT_PARENTHESIS",
	'*': "ASTERISK", '+': "PLUS_SIGN", ',': "COMMA", '-': "HYPHEN",
	'.': "PERIOD", '/': "SLASH", ':': "COLON", ';': "SEMICOLON",
	'<': "LESS_THAN", '=': "EQUAL_SIGN", '>': "GREATER_THAN",
	'?': "QUESTION_MARK", '@': "AT_SIGN", '[': "LEFT_BRACKET",
	'\\': "BACKSLASH", ']': "RIGHT_BRACKET", '^': "CIRCUMFLEX",
	'_': "UNDERSCORE", '`': "GRAVE", '{': "LEFT_BRACE", '|': "VERTICAL_BAR",
	'}': "RIGHT_BRACE", '~': "TILDE",
}

// natoAlphabet spells letters; capital letters are capitalized words.
var natoAlphabet = [26]string{
	"alfa", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel",
	"india", "juliett", "kilo", "lima", "mike", "november", "oscar", "papa",
	"quebec", "romeo", "sierra", "tango", "uniform", "victor", "whiskey",
	"x_ray", "yankee", "zulu",
}

// symbolName returns the name of a digit or special symbol, or the symbol
// itself for letters.
func symbolName(c byte) string {
	if n, ok := symbolNames[c]; ok {
		return n
	}
	return string(c)
}

// spell spells out a password for reading aloud, e.g. "aB1" becomes
// "alfa-Bravo-ONE".
func spell(p string) string {
	words := make([]string, 0, len(p))
	for i := 0; i < len(p); i++ {
		c := p[i]
		switch {
		case c >= 'a' && c <= 'z':
			words = append(words, natoAlphabet[c-'a'])
		case c >= 'A' && c <= 'Z':
			w := natoAlphabet[c-'A']
			words = append(words, strings.ToUpper(w[:1])+w[1:])
		default:
			if n, ok := symbolNames[c]; ok {
				words = append(words, n)
			}
		}
	}
	return strings.Join(words, "-")
}
