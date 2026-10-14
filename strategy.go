package main

import "fmt"

// strategy selects which character groups a password is built from.
type strategy string

const (
	strategyFull     strategy = "full"
	strategyAlphanum strategy = "alphanum"
	strategyDigits   strategy = "digits"
)

func parseStrategy(s string) (strategy, error) {
	switch st := strategy(s); st {
	case strategyFull, strategyAlphanum, strategyDigits:
		return st, nil
	}
	return "", fmt.Errorf("unknown strategy %q (want full, alphanum or digits)", s)
}

// groups returns the character groups the strategy draws from.
func (s strategy) groups() []string {
	switch s {
	case strategyAlphanum:
		return []string{lowerChars, upperChars, digitChars}
	case strategyDigits:
		return []string{digitChars}
	default:
		return []string{lowerChars, upperChars, digitChars, specialChars}
	}
}

// strategyFlag lets a strategy be set from the command line.
type strategyFlag struct{ s *strategy }

func (f strategyFlag) String() string {
	if f.s == nil {
		return ""
	}
	return string(*f.s)
}

func (f strategyFlag) Set(v string) error {
	st, err := parseStrategy(v)
	if err != nil {
		return err
	}
	*f.s = st
	return nil
}
