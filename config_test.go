package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	input := `
# comment
default_length = 24
default_strategy=alphanum
`
	cfg := defaultConfig()
	if err := parseConfig(strings.NewReader(input), &cfg); err != nil {
		t.Fatal(err)
	}
	want := config{length: 24, strategy: strategyAlphanum}
	if cfg != want {
		t.Errorf("got %+v, want %+v", cfg, want)
	}
}

func TestParseConfigStrategies(t *testing.T) {
	for _, name := range []string{"full", "alphanum", "digits"} {
		cfg := config{length: 10, strategy: strategyFull}
		if err := parseConfig(strings.NewReader("default_strategy="+name), &cfg); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if cfg.strategy != strategy(name) {
			t.Errorf("%s: got strategy %q", name, cfg.strategy)
		}
		if cfg.length != 10 {
			t.Errorf("%s: length changed to %d", name, cfg.length)
		}
	}
}

func TestParseConfigErrors(t *testing.T) {
	cases := map[string]string{
		"no equals":    "default_length",
		"bad length":   "default_length=abc",
		"bad strategy": "default_strategy=symbols",
		"unknown key":  "length=10",
	}
	for name, input := range cases {
		cfg := defaultConfig()
		err := parseConfig(strings.NewReader(input), &cfg)
		if err == nil {
			t.Errorf("%s: expected error", name)
		}
	}

	cfg := defaultConfig()
	err := parseConfig(strings.NewReader("\n\ndefault_length=x"), &cfg)
	if err == nil || !strings.Contains(err.Error(), "line 3") {
		t.Errorf("expected error to mention line 3, got %v", err)
	}
}

func TestLoadConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("missing file: %v", err)
	}
	if cfg != defaultConfig() {
		t.Errorf("missing file: got %+v, want defaults", cfg)
	}

	path := filepath.Join(home, configFileName)
	if err := os.WriteFile(path, []byte("default_length=16\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err = loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.length != 16 || cfg.strategy != strategyFull {
		t.Errorf("got %+v", cfg)
	}

	if err := os.WriteFile(path, []byte("oops\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig(); err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("expected error mentioning %s, got %v", path, err)
	}
}
