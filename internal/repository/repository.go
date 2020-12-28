package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Root(pwd string) (string, error) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
	)
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = pwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}
