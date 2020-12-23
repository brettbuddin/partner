package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atrox/homedir"
	"github.com/brettbuddin/partner/internal/manifest"
)

// Command holds actions for commands
type Command struct {
	Paths Paths
}

// New returns a new Command
func New(paths Paths) *Command {
	return &Command{Paths: paths}
}

// Paths contains paths to directories
type Paths struct {
	WorkingDir   string
	TemplateFile string
	ManifestFile string
}

// DefaultPaths returns calculated Paths based on the environment
func DefaultPaths(pwd string) (Paths, error) {
	manifestPath := os.Getenv("PARTNER_MANIFEST")
	if manifestPath == "" {
		manifestPath = manifest.DefaultPath
	}
	manifestPath, err := homedir.Expand(manifestPath)
	if err != nil {
		return Paths{}, err
	}

	templatePath, err := templatePath(pwd)
	if err != nil {
		return Paths{}, err
	}

	return Paths{
		WorkingDir:   pwd,
		TemplateFile: templatePath,
		ManifestFile: os.ExpandEnv(manifestPath),
	}, nil
}

func templatePath(pwd string) (string, error) {
	path, err := repositoryRoot(pwd)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, ".git/gitmessage.txt"), nil
}

func repositoryRoot(pwd string) (string, error) {
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
