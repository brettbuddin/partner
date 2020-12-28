package command

import (
	"io"

	"github.com/brettbuddin/partner/internal/manifest"
)

// ManifestList lists all coauthors
func (c *Command) ManifestList(w io.Writer) error {
	m, err := manifest.Load(c.Paths.ManifestFile)
	if err != nil {
		return err
	}
	if len(m.Coauthors) == 0 {
		return nil
	}
	return writeList(w, m.Slice()...)
}

// ManifestRemove removes a coauthor from the Manifest
func (c *Command) ManifestRemove(ids ...string) error {
	m, err := manifest.Load(c.Paths.ManifestFile)
	if err != nil {
		return err
	}
	if err := m.Remove(ids...); err != nil {
		return err
	}
	return manifest.WriteFile(c.Paths.ManifestFile, m)
}

// UserFetcher fetches coauthor information from somewhere else
type UserFetcher interface {
	Fetch(username string) (manifest.Coauthor, error)
}

// ManifestFetchAdd adds a coauthor by looking up their information remotely
func (c *Command) ManifestFetchAdd(fetcher UserFetcher, usernames ...string) error {
	var coauthors []manifest.Coauthor
	for _, username := range usernames {
		coauthor, err := fetcher.Fetch(username)
		if err != nil {
			return err
		}
		coauthors = append(coauthors, coauthor)
	}

	m, err := manifest.Load(c.Paths.ManifestFile)
	if err != nil {
		return err
	}
	if err := m.Add(coauthors...); err != nil {
		return err
	}
	return manifest.WriteFile(c.Paths.ManifestFile, m)
}

// ManifestAdd adds a coauthor using manually entered information
func (c *Command) ManifestAdd(id, name, email string) error {
	m, err := manifest.Load(c.Paths.ManifestFile)
	if err != nil {
		return err
	}
	err = m.Add(manifest.Coauthor{
		ID:    id,
		Name:  name,
		Email: email,
		Type:  manifest.CoauthorTypeManual,
	})
	if err != nil {
		return err
	}
	return manifest.WriteFile(c.Paths.ManifestFile, m)
}
