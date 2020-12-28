package repository

import (
	"fmt"
	"os/exec"
)

func SetCommitTemplate(dir string, templatePath string) error {
	cmd := exec.Command("git", "config", "commit.template", templatePath)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set commit template: %w", err)
	}
	return nil
}

func UnsetCommitTemplate(dir string) error {
	cmd := exec.Command("git", "config", "--unset", "commit.template")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unset commit template: %w", err)
	}
	return nil
}
