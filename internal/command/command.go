package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/atrox/homedir"
	"github.com/brettbuddin/partner/internal/manifest"
	"github.com/brettbuddin/partner/internal/repository"
)

// Command holds actions for commands
type Command struct {
	Paths Paths
}

// New returns a new Command
func New(paths Paths) *Command {
	return &Command{Paths: paths}
}

// Paths contains paths necessary for partner to do its work
type Paths struct {
	ManifestFile string
	Repository   RepositoryPaths
}

// Repository specific paths
type RepositoryPaths struct {
	Root         string
	TemplateFile string
}

// DefaultPaths returns calculated Git repository root, commit template and
// manifest paths relative to the current working directory.
func DefaultPaths(pwd string) (Paths, error) {
	root, err := repository.Root(pwd)
	if err != nil {
		return Paths{}, err
	}

	manifestPath := os.Getenv("PARTNER_MANIFEST")
	if manifestPath == "" {
		manifestPath = manifest.DefaultPath
	}
	manifestPath, err = homedir.Expand(manifestPath)
	if err != nil {
		return Paths{}, err
	}

	return Paths{
		ManifestFile: os.ExpandEnv(manifestPath),
		Repository: RepositoryPaths{
			Root:         root,
			TemplateFile: filepath.Join(root, ".git/gitmessage.txt"),
		},
	}, nil
}

func writeList(w io.Writer, coauthors ...manifest.Coauthor) error {
	if len(coauthors) == 0 {
		return nil
	}

	sort.Slice(coauthors, func(i, j int) bool {
		return strings.ToLower(coauthors[i].ID) < strings.ToLower(coauthors[j].ID)
	})

	tabw := tabwriter.NewWriter(w, 5, 2, 2, ' ', 0)
	fmt.Fprintln(tabw, "ID\tNAME\tEMAIL\tTYPE")
	for _, ca := range coauthors {
		fmt.Fprintf(tabw, "%s\t%s\t%s\t%s\n", ca.ID, ca.Name, ca.Email, ca.Type)
	}
	return tabw.Flush()
}
