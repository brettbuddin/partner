package command

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/brettbuddin/partner/internal/manifest"
	"github.com/brettbuddin/partner/internal/repository"
	"github.com/brettbuddin/partner/internal/template"
)

// TemplateStatus lists active coauthors
func (c *Command) TemplateStatus(w io.Writer) error {
	m, err := manifest.Load(c.Paths.ManifestFile)
	if err != nil {
		return err
	}

	repoPaths, err := c.Paths.Repository()
	if err != nil {
		return err
	}

	ids, err := template.ExtractIDs(repoPaths.TemplateFile)
	if err != nil {
		return err
	}
	active, err := m.Find(ids...)
	if err != nil {
		return err
	}
	return writeList(w, active...)
}

// TemplateSet activates a coauthor in the Template
func (c *Command) TemplateSet(ids ...string) error {
	m, err := manifest.Load(c.Paths.ManifestFile)
	if err != nil {
		return err
	}

	repoPaths, err := c.Paths.Repository()
	if err != nil {
		return err
	}

	existingIDs, err := template.ExtractIDs(repoPaths.TemplateFile)
	if err != nil {
		return err
	}
	ids = uniqueStrings(append(ids, existingIDs...))

	coauthors, err := m.Find(ids...)
	if err != nil {
		return err
	}

	t := template.Template{Coauthors: coauthors}
	if err := template.WriteFile(repoPaths.TemplateFile, t); err != nil {
		return err
	}
	return repository.SetCommitTemplate(repoPaths.Root, repoPaths.TemplateFile)
}

// TemplateClear emptys the coauthors Template
func (c *Command) TemplateClear() error {
	repoPaths, err := c.Paths.Repository()
	if err != nil {
		return err
	}

	defer repository.UnsetCommitTemplate(repoPaths.Root)

	if err := os.Remove(repoPaths.TemplateFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to remove commit template: %w", err)
	}
	return nil
}

func uniqueStrings(ids []string) []string {
	var (
		uniq = make(map[string]bool)
		out  []string
	)
	for _, id := range ids {
		if _, v := uniq[id]; !v {
			uniq[id] = true
			out = append(out, id)
		}
	}
	return out
}
