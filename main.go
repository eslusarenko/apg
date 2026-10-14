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
  -a, --alnum        letters and digits only
  -n, --no-newline   do not print the trailing newline
  -c, --copy         copy the password to the clipboard
  -v, --version      print version and exit
  -h, --help         show this help
`

func main() {
	var (
		length      int
		alnum       bool
		noNewline   bool
		toClipboard bool
		showVersion bool
	)
	flag.IntVar(&length, "l", defaultLength, "")
	flag.IntVar(&length, "length", defaultLength, "")
	flag.BoolVar(&alnum, "a", false, "")
	flag.BoolVar(&alnum, "alnum", false, "")
	flag.BoolVar(&noNewline, "n", false, "")
	flag.BoolVar(&noNewline, "no-newline", false, "")
	flag.BoolVar(&toClipboard, "c", false, "")
	flag.BoolVar(&toClipboard, "copy", false, "")
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

	password, err := Generate(length, !alnum)
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
