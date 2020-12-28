package command

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

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
	var coauthors []manifest.Coauthor
	for _, ca := range m.Coauthors {
		coauthors = append(coauthors, ca)
	}
	return printCoauthors(w, coauthors...)
}

func printCoauthors(w io.Writer, coauthors ...manifest.Coauthor) error {
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
