// Command apg (Advanced Password Generator) generates strong passwords. Its
// options follow the classic apg (Automated Password Generator).
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultLength = 32

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

const usage = `apg - Advanced Password Generator

Usage:
  apg [options]

Options (compatible with the classic apg):
  -a 0|1         algorithm: 0 pronounceable, 1 random (default 1)
  -M mode        symbol classes: S special symbols, N numbers, C capital
                 letters, L small letters; upper case means every password
                 must contain the class, lower case means it may (default %s)
  -E chars       exclude these characters
  -n N           number of passwords to generate (default 1)
  -m N           minimum password length (default %d)
  -x N           maximum password length (default %d)
  -r file        reject passwords found in a dictionary file
  -t             print the pronunciation of pronounceable passwords
  -l             spell out generated passwords
  -d             do not print a newline after each password
  -q             quiet mode: do not print warnings
  -s, -c seed    random seed options; accepted for compatibility and ignored
  -v, --version  print version and exit
  -h, --help     show this help

Extensions:
  --length N        exact password length, same as -m N -x N
  --strategy NAME   full (-M SNCL), alphanum (-M NCL) or digits (-M N)
  --alnum           same as --strategy alphanum
  --digits, --pin   same as --strategy digits
  --no-newline      same as -d
  -C, --copy        copy the generated passwords to the clipboard

Defaults can be changed in ~/.apg.conf, one key=value per line:
  default_length=24
  default_strategy=alphanum

Values for default_strategy are the strategy names above.
Command-line flags take precedence over the file.
`

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fail(2, err)
	}

	var (
		alg         int
		modeStr     = cfg.strategy.modeString()
		exclude     string
		count       int
		minLen      int
		maxLen      int
		dictFile    string
		showPron    bool
		spellOut    bool
		noNewline   bool
		quiet       bool
		seedIgnored bool
		toClipboard bool
		showVersion bool
	)

	flag.IntVar(&alg, "a", int(algRandom), "")
	flag.Func("M", "", func(v string) error {
		if _, err := parseMode(v); err != nil {
			return err
		}
		modeStr = v
		return nil
	})
	flag.StringVar(&exclude, "E", "", "")
	flag.IntVar(&count, "n", 1, "")
	flag.IntVar(&minLen, "m", cfg.length, "")
	flag.IntVar(&maxLen, "x", cfg.length, "")
	flag.StringVar(&dictFile, "r", "", "")
	flag.BoolVar(&showPron, "t", false, "")
	flag.BoolVar(&spellOut, "l", false, "")
	flag.BoolVar(&noNewline, "d", false, "")
	flag.BoolVar(&noNewline, "no-newline", false, "")
	flag.BoolVar(&quiet, "q", false, "")
	flag.BoolFunc("s", "", func(string) error { seedIgnored = true; return nil })
	flag.Func("c", "", func(string) error { seedIgnored = true; return nil })
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")

	flag.Func("length", "", func(v string) error {
		n, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("must be a number")
		}
		minLen, maxLen = n, n
		return nil
	})
	flag.Func("strategy", "", func(v string) error {
		st, err := parseStrategy(v)
		if err != nil {
			return err
		}
		modeStr = st.modeString()
		return nil
	})
	setStrategy := func(st strategy) func(string) error {
		return func(v string) error {
			if v == "true" {
				modeStr = st.modeString()
			}
			return nil
		}
	}
	flag.BoolFunc("alnum", "", setStrategy(strategyAlphanum))
	flag.BoolFunc("digits", "", setStrategy(strategyDigits))
	flag.BoolFunc("pin", "", setStrategy(strategyDigits))
	flag.BoolVar(&toClipboard, "C", false, "")
	flag.BoolVar(&toClipboard, "copy", false, "")

	// Options of the classic apg that are not supported.
	unsupported := func(what string) func(string) error {
		return func(string) error { return fmt.Errorf("%s is not supported", what) }
	}
	flag.Func("b", "", unsupported("Bloom filter checking"))
	flag.Func("p", "", unsupported("Bloom filter checking"))
	flag.BoolFunc("y", "", unsupported("printing crypt(3) hashes"))
	flag.BoolFunc("k", "", unsupported("cracklib checking"))

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, cfg.strategy.modeString(), cfg.length, cfg.length)
	}
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "apg: unexpected argument %q\n\n", flag.Arg(0))
		flag.Usage()
		os.Exit(2)
	}
	if showVersion {
		fmt.Println("apg (Advanced Password Generator)", version)
		return
	}
	if count < 1 {
		fail(2, errors.New("number of passwords must be at least 1"))
	}
	if minLen > maxLen {
		maxLen = minLen
	}

	warn := func(msg string) {
		if !quiet {
			fmt.Fprintln(os.Stderr, "apg: warning:", msg)
		}
	}
	if seedIgnored {
		warn("-s and -c are ignored: passwords always use the system's cryptographic random generator")
	}
	if algorithm(alg) == algPronounceable && strings.ContainsAny(exclude, classSymbols[classSmall]) {
		warn("small letters cannot be excluded from pronounceable passwords")
	}

	m, _ := parseMode(modeStr)
	opts := options{
		algorithm: algorithm(alg),
		minLen:    minLen,
		maxLen:    maxLen,
		mode:      m,
		exclude:   exclude,
	}
	if dictFile != "" {
		if opts.dict, err = loadDict(dictFile); err != nil {
			fail(2, err)
		}
	}
	gen, err := newGenerator(opts)
	if err != nil {
		fail(2, err)
	}

	generated := make([]string, 0, count)
	for range count {
		p, err := gen.next()
		if err != nil {
			fail(1, err)
		}
		line := p.text
		if showPron && opts.algorithm == algPronounceable {
			line += " (" + p.hyphenated + ")"
		}
		if spellOut {
			line += " " + spell(p.text)
		}
		if !noNewline {
			line += "\n"
		}
		fmt.Print(line)
		generated = append(generated, p.text)
	}

	if toClipboard {
		if err := copyToClipboard(strings.Join(generated, "\n")); err != nil {
			fail(1, fmt.Errorf("copy to clipboard: %w", err))
		}
		if !quiet {
			fmt.Fprintln(os.Stderr, "Copied to clipboard.")
		}
	}
}

func fail(code int, err error) {
	fmt.Fprintln(os.Stderr, "apg:", err)
	os.Exit(code)
}
