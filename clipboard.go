package main

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// copyToClipboard writes s to the system clipboard using the platform's
// standard clipboard utility.
func copyToClipboard(s string) error {
	cmd, err := clipboardCmd()
	if err != nil {
		return err
	}
	cmd.Stdin = strings.NewReader(s)
	return cmd.Run()
}

func clipboardCmd() (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("pbcopy"), nil
	case "windows":
		return exec.Command("clip"), nil
	}

	var candidates [][]string
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		candidates = append(candidates, []string{"wl-copy"})
	}
	candidates = append(candidates,
		[]string{"xclip", "-selection", "clipboard"},
		[]string{"xsel", "--clipboard", "--input"},
	)
	for _, c := range candidates {
		if path, err := exec.LookPath(c[0]); err == nil {
			return exec.Command(path, c[1:]...), nil
		}
	}
	return nil, errors.New("no clipboard utility found (install wl-clipboard, xclip or xsel)")
}
