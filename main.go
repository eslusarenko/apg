// Command apg generates strong random passwords.
package main

import (
	"flag"
	"fmt"
	"os"
)

const defaultLength = 32

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

const usage = `apg - generate strong random passwords

Usage:
  apg [options]

Options:
  -l, --length N        password length (default %d)
  -s, --strategy NAME   character set to use (default %s, see below)
  -a, --alnum           same as --strategy alphanum
  -d, --digits          same as --strategy digits
      --pin             same as --strategy digits
  -n, --no-newline      do not print the trailing newline
  -c, --copy            copy the password to the clipboard
  -v, --version         print version and exit
  -h, --help            show this help

Strategies:
  full       letters, digits and special characters
  alphanum   letters and digits
  digits     digits only, e.g. for PIN codes

Defaults can be changed in ~/.apg.conf, one key=value per line:
  default_length=24
  default_strategy=alphanum

Values for default_strategy are the strategy names above.
Command-line flags take precedence over the file.
`

func main() {
	var (
		length      int
		strat       strategy
		noNewline   bool
		toClipboard bool
		showVersion bool
	)
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "apg:", err)
		os.Exit(2)
	}

	flag.IntVar(&length, "l", cfg.length, "")
	flag.IntVar(&length, "length", cfg.length, "")
	strat = cfg.strategy
	flag.Var(strategyFlag{&strat}, "s", "")
	flag.Var(strategyFlag{&strat}, "strategy", "")
	setStrategy := func(st strategy) func(string) error {
		return func(v string) error {
			if v == "true" {
				strat = st
			}
			return nil
		}
	}
	flag.BoolFunc("a", "", setStrategy(strategyAlphanum))
	flag.BoolFunc("alnum", "", setStrategy(strategyAlphanum))
	flag.BoolFunc("d", "", setStrategy(strategyDigits))
	flag.BoolFunc("digits", "", setStrategy(strategyDigits))
	flag.BoolFunc("pin", "", setStrategy(strategyDigits))
	flag.BoolVar(&noNewline, "n", false, "")
	flag.BoolVar(&noNewline, "no-newline", false, "")
	flag.BoolVar(&toClipboard, "c", false, "")
	flag.BoolVar(&toClipboard, "copy", false, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, cfg.length, cfg.strategy)
	}
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "apg: unexpected argument %q\n\n", flag.Arg(0))
		flag.Usage()
		os.Exit(2)
	}
	if showVersion {
		fmt.Println("apg", version)
		return
	}

	password, err := Generate(length, strat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "apg:", err)
		os.Exit(2)
	}
	if noNewline {
		fmt.Print(password)
	} else {
		fmt.Println(password)
	}

	if toClipboard {
		if err := copyToClipboard(password); err != nil {
			fmt.Fprintln(os.Stderr, "apg: copy to clipboard:", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Copied to clipboard.")
	}
}
