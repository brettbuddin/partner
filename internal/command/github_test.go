package command

import (
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

func TestGitHubFetcher(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/brettbuddin" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		f, err := os.Open("testdata/github_user.json")
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
