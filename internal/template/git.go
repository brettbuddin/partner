package template

import (
	"fmt"
	"os/exec"
)

func Set(workDir string, templatePath string) error {
	cmd := exec.Command("git", "config", "commit.template", templatePath)
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set commit template: %w", err)
	}
	return nil
}

func Unset(workDir string) error {
	cmd := exec.Command("git", "config", "--unset", "commit.template")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set commit template: %w", err)
	}
	return nil
}
