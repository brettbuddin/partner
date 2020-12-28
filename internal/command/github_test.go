package command

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/brettbuddin/partner/internal/manifest"
	"github.com/stretchr/testify/require"
)

func TestGitHubAdd(t *testing.T) {
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

func TestGitHubFetcher(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/brettbuddin" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		f, err := os.Open("testdata/user.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	defer server.Close()

	f := GitHubFetcher{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseURL: server.URL,
	}
	ca, err := f.Fetch("brettbuddin")
	require.NoError(t, err)
	require.Equal(t, manifest.Coauthor{
		ID:    "brettbuddin",
		Name:  "Brett Buddin",
		Email: "6059+brettbuddin@users.noreply.github.com",
		Type:  manifest.CoauthorTypeGitHub,
	}, ca)
}

func TestGitHubFetcher_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"message": "Not Found"}`)
	}))
	defer server.Close()

	f := GitHubFetcher{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseURL: server.URL,
	}
	_, err := f.Fetch("brettbuddin-doesntexist")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Not Found")
}
