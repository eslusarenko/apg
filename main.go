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
  -l, --length N     password length (default %d)
  -n, --no-special   exclude special characters
  -v, --version      print version and exit
  -h, --help         show this help
`

func main() {
	var (
		length      int
		noSpecial   bool
		showVersion bool
	)
	flag.IntVar(&length, "l", defaultLength, "")
	flag.IntVar(&length, "length", defaultLength, "")
	flag.BoolVar(&noSpecial, "n", false, "")
	flag.BoolVar(&noSpecial, "no-special", false, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, defaultLength)
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

	password, err := Generate(length, !noSpecial)
	if err != nil {
		fmt.Fprintln(os.Stderr, "apg:", err)
		os.Exit(2)
	}
	fmt.Println(password)
}
