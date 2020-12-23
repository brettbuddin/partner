package template

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/brettbuddin/partner/internal/manifest"
	"github.com/stretchr/testify/require"
)

func TestActiveCoauthors(t *testing.T) {
	ids, err := Active("testdata/template")
	require.NoError(t, err)
	require.Equal(t, []string{"persona", "personb"}, ids)
}

func TestSave(t *testing.T) {
	dir, err := ioutil.TempDir("", "template")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	tmpl := Template{
		Coauthors: []manifest.Coauthor{
			{
				ID:    "persona",
				Name:  "Person A",
				Type:  manifest.CoauthorTypeManual,
				Email: "a@buddin.org",
			},
			{
				ID:    "personb",
				Name:  "Person B",
				Type:  manifest.CoauthorTypeManual,
				Email: "b@buddin.org",
			},
		},
	}

	path := filepath.Join(dir, "gitmessage.txt")
	err = tmpl.Save(path)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("testdata/template")
	require.NoError(t, err)
	actual, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
