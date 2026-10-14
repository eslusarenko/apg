package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{};:,.<>?/|~"
)

// Generate returns a cryptographically random password of the given length.
// The result always contains at least one character from every group of the
// strategy.
func Generate(length int, s strategy) (string, error) {
	groups := s.groups()
	if length < len(groups) {
		return "", fmt.Errorf("length must be at least %d", len(groups))
	}

	buf := make([]byte, 0, length)
	alphabet := ""
	for _, g := range groups {
		c, err := pick(g)
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
		alphabet += g
	}
	for len(buf) < length {
		c, err := pick(alphabet)
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}

	// Fisher–Yates shuffle so the guaranteed characters are not at fixed positions.
	for i := len(buf) - 1; i > 0; i-- {
		j, err := randInt(i + 1)
		if err != nil {
			return "", err
		}
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf), nil
}

func pick(s string) (byte, error) {
	i, err := randInt(len(s))
	if err != nil {
		return 0, err
	}
	return s[i], nil
}

// randInt returns a uniform random integer in [0, n).
func randInt(n int) (int, error) {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}
