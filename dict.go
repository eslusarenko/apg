package main

import (
	"bufio"
	"os"
	"strings"
)

// loadDict reads a dictionary file for -r: words separated by whitespace,
// usually one per line.
func loadDict(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	words := make(map[string]struct{})
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		for _, w := range strings.Fields(sc.Text()) {
			words[w] = struct{}{}
		}
	}
	return words, sc.Err()
}
