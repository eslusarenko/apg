package main

import "fmt"

// strategy is a named preset for the -M mode.
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

// modeString returns the -M mode the strategy stands for.
func (s strategy) modeString() string {
	switch s {
	case strategyAlphanum:
		return "NCL"
	case strategyDigits:
		return "N"
	default:
		return "SNCL"
	}
}
