package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const DefaultPath = "~/.config/partner/manifest.json"

// Manifest contains all coauthors
type Manifest struct {
	Coauthors map[string]Coauthor `json:"coauthors"`
}

func (m Manifest) Slice() []Coauthor {
	var l []Coauthor
	for _, ca := range m.Coauthors {
		l = append(l, ca)
	}
	return l
}

// Coauthor is an identity to be referenced as a commit coauthor
type Coauthor struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Coauthor types
const (
	CoauthorTypeGitHub = "github"
	CoauthorTypeManual = "manual"
)

// Load reads a Manifest
func Load(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Manifest{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var m Manifest
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	if m.Coauthors == nil {
		m.Coauthors = map[string]Coauthor{}
	}
	return &m, nil
}

// Find looks up a list of coauthors by their IDs
func (m *Manifest) Find(ids ...string) ([]Coauthor, error) {
	var coauthors []Coauthor
	for _, id := range ids {
		ca, ok := m.Coauthors[strings.ToLower(id)]
		if !ok {
			return nil, fmt.Errorf("unknown coauthor %q", id)
		}
		coauthors = append(coauthors, ca)
	}
	sort.Slice(coauthors, func(i, j int) bool {
		return strings.ToLower(coauthors[i].ID) < strings.ToLower(coauthors[j].ID)
	})
	return coauthors, nil
}

// Remove removes coauthors by their IDs
func (m *Manifest) Remove(ids ...string) error {
	if m.Coauthors == nil {
		m.Coauthors = map[string]Coauthor{}
	}
	for _, id := range ids {
		key := strings.ToLower(id)
		if _, ok := m.Coauthors[key]; !ok {
			return fmt.Errorf("unknown coauthor %q", id)
		}
		delete(m.Coauthors, key)
	}
	return nil
}

// Add adds coauthors to the Manifest
func (m *Manifest) Add(coauthors ...Coauthor) error {
	if m.Coauthors == nil {
		m.Coauthors = map[string]Coauthor{}
	}
	for _, newCA := range coauthors {
		key := strings.ToLower(newCA.ID)
		if _, ok := m.Coauthors[key]; ok {
			return fmt.Errorf("coauthor with ID %q already exists", newCA.ID)
		}
		m.Coauthors[key] = newCA
	}
	return nil
}

// WriteFile saves the manifest to a JSON file
func WriteFile(path string, m *Manifest) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if m.Coauthors == nil {
		m.Coauthors = map[string]Coauthor{}
	}

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	if err := e.Encode(m); err != nil {
		return err
	}
	return nil
}
