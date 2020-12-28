package command

import (
	"bytes"
	"testing"

	"github.com/brettbuddin/partner/internal/manifest"
	"github.com/stretchr/testify/require"
)

func TestManifestFetchAdd(t *testing.T) {
	cmd := New(newWorkspace(t))
	f := fetcher{
		coauthor: manifest.Coauthor{
			ID:    "brettbuddin",
			Name:  "Brett Buddin",
			Email: "6059+brettbuddin@users.noreply.github.com",
			Type:  manifest.CoauthorTypeGitHub,
		},
	}
	err := cmd.ManifestFetchAdd(f, "brettbuddin")
	require.NoError(t, err)
	err = cmd.TemplateSet("brettbuddin")
	require.NoError(t, err)

	out := bytes.NewBuffer(nil)
	err = cmd.TemplateStatus(out)
	require.NoError(t, err)
	require.Equal(t, listExample(`
ID           NAME          EMAIL                                      TYPE
brettbuddin  Brett Buddin  6059+brettbuddin@users.noreply.github.com  github
`), out.String())
}

type fetcher struct {
	coauthor manifest.Coauthor
	err      error
}

func (f fetcher) Fetch(username string) (manifest.Coauthor, error) {
	return f.coauthor, f.err
}
