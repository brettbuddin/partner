package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/brettbuddin/partner/internal/manifest"
)

type GitHubFetcher struct {
	Client  *http.Client
	BaseURL string
}

func (f *GitHubFetcher) Fetch(username string) (manifest.Coauthor, error) {
	parsed, err := url.Parse(fmt.Sprintf("%s/users/%s", f.BaseURL, username))
	if err != nil {
		return manifest.Coauthor{}, err
	}
	r, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		return manifest.Coauthor{}, err
	}
	resp, err := f.Client.Do(r)
	if err != nil {
		return manifest.Coauthor{}, err
	}

	if resp.StatusCode != http.StatusOK {
		var ghError struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&ghError); err != nil {
			return manifest.Coauthor{}, err
		}
		return manifest.Coauthor{}, fmt.Errorf("error fetching %q from GitHub: %s", username, ghError.Message)
	}

	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return manifest.Coauthor{}, err
	}
	return manifest.Coauthor{
		Email: fmt.Sprintf("%d+%s@users.noreply.github.com", user.ID, user.Login),
		ID:    user.Login,
		Name:  user.Name,
		Type:  manifest.CoauthorTypeGitHub,
	}, nil
}
