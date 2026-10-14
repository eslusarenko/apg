package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const configFileName = ".apg.conf"

// config holds the defaults that command-line flags start from.
type config struct {
	length   int
	strategy strategy
}

func defaultConfig() config {
	return config{length: defaultLength, strategy: strategyFull}
}

// loadConfig reads ~/.apg.conf on top of the built-in defaults.
// A missing file is not an error.
func loadConfig() (config, error) {
	cfg := defaultConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}
	path := filepath.Join(home, configFileName)
	f, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	if err := parseConfig(f, &cfg); err != nil {
		return cfg, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

// parseConfig applies key=value lines from r to cfg. Blank lines and lines
// starting with # are ignored.
func parseConfig(r io.Reader, cfg *config) error {
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("line %d: expected key=value", n)
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)

		switch key {
		case "default_length":
			l, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("line %d: default_length must be a number, got %q", n, value)
			}
			cfg.length = l
		case "default_strategy":
			st, err := parseStrategy(value)
			if err != nil {
				return fmt.Errorf("line %d: default_strategy: %w", n, err)
			}
			cfg.strategy = st
		default:
			return fmt.Errorf("line %d: unknown key %q", n, key)
		}
	}
	return sc.Err()
}
