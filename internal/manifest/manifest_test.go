package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	m, err := Load("testdata/manifest.json")
	require.NoError(t, err)
	require.Equal(t, map[string]Coauthor{
		"georgemac": {
			ID:    "GeorgeMac",
			Name:  "George",
			Email: "1253326+GeorgeMac@users.noreply.github.com",
			Type:  "github",
		},
		"gavincabbage": {
			ID:    "gavincabbage",
			Name:  "Gavin Cabbage",
			Email: "5225414+gavincabbage@users.noreply.github.com",
			Type:  "github",
		},
		"stuartcarnie": {
			ID:    "stuartcarnie",
			Name:  "Stuart Carnie",
			Email: "52852+stuartcarnie@users.noreply.github.com",
			Type:  "github",
		},
	}, m.Coauthors)
}

func TestFind(t *testing.T) {
	m, err := Load("testdata/manifest.json")
	require.NoError(t, err)

	coauthors, err := m.Find("GeorgeMac", "gavincabbage")
	require.NoError(t, err)

	require.Equal(t, []Coauthor{
		{
			ID:    "gavincabbage",
			Name:  "Gavin Cabbage",
			Email: "5225414+gavincabbage@users.noreply.github.com",
			Type:  "github",
		},
		{
			ID:    "GeorgeMac",
			Name:  "George",
			Email: "1253326+GeorgeMac@users.noreply.github.com",
			Type:  "github",
		},
	}, coauthors)
}

func TestRemove(t *testing.T) {
	m, err := Load("testdata/manifest.json")
	require.NoError(t, err)

	m.Remove("gavincabbage")

	require.Equal(t, map[string]Coauthor{
		"georgemac": {
			ID:    "GeorgeMac",
			Name:  "George",
			Email: "1253326+GeorgeMac@users.noreply.github.com",
			Type:  "github",
		},
		"stuartcarnie": {
			ID:    "stuartcarnie",
			Name:  "Stuart Carnie",
			Email: "52852+stuartcarnie@users.noreply.github.com",
			Type:  "github",
		},
	}, m.Coauthors)
}

func TestAdd(t *testing.T) {
	m, err := Load("testdata/manifest.json")
	require.NoError(t, err)

	brett := Coauthor{
		ID:    "brettbuddin",
		Name:  "Brett Buddin",
		Email: "brett@buddin.org",
		Type:  "manual",
	}

	m.Add(brett)

	require.Equal(t, map[string]Coauthor{
		"georgemac": {
			ID:    "GeorgeMac",
			Name:  "George",
			Email: "1253326+GeorgeMac@users.noreply.github.com",
			Type:  "github",
		},
		"gavincabbage": {
			ID:    "gavincabbage",
			Name:  "Gavin Cabbage",
			Email: "5225414+gavincabbage@users.noreply.github.com",
			Type:  "github",
		},
		"stuartcarnie": {
			ID:    "stuartcarnie",
			Name:  "Stuart Carnie",
			Email: "52852+stuartcarnie@users.noreply.github.com",
			Type:  "github",
		},
		"brettbuddin": {
			ID:    "brettbuddin",
			Name:  "Brett Buddin",
			Email: "brett@buddin.org",
			Type:  "manual",
		},
	}, m.Coauthors)
}

func TestSave(t *testing.T) {
	dir, err := ioutil.TempDir("", "manifest")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	m, err := Load("testdata/manifest.json")
	require.NoError(t, err)

	m.Remove("stuartcarnie")

	newFile := filepath.Join(dir, "changed.json")
	err = WriteFile(newFile, m)
	require.NoError(t, err)

	cm, err := Load(newFile)
	require.NoError(t, err)

	require.Equal(t, map[string]Coauthor{
		"georgemac": {
			ID:    "GeorgeMac",
			Name:  "George",
			Email: "1253326+GeorgeMac@users.noreply.github.com",
			Type:  "github",
		},
		"gavincabbage": {
			ID:    "gavincabbage",
			Name:  "Gavin Cabbage",
			Email: "5225414+gavincabbage@users.noreply.github.com",
			Type:  "github",
		},
	}, cm.Coauthors)
}
