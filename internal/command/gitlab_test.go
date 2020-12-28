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

func TestGitLabFetcher(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v4/users" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.FormValue("username") != "brettbuddin" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		f, err := os.Open("testdata/gitlab_user.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	defer server.Close()

	f := GitLabFetcher{
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
		Email: "4453516-brettbuddin@users.noreply.gitlab.com",
		Type:  manifest.CoauthorTypeGitLab,
	}, ca)
}

func TestGitLabFetcher_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[]`)
	}))
	defer server.Close()

	f := GitLabFetcher{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseURL: server.URL,
	}
	_, err := f.Fetch("brettbuddin-doesntexist")
	require.Error(t, err)
	require.Contains(t, err.Error(), "username not found")
}
