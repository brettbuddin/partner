package command

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brettbuddin/partner/internal/manifest"
)

type GitHubFetcher struct {
	Client  *http.Client
	BaseURL string
}

func (gh *GitHubFetcher) Fetch(username string) (manifest.Coauthor, error) {
	url := fmt.Sprintf("%s/users/%s", gh.BaseURL, username)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return manifest.Coauthor{}, err
	}
	resp, err := gh.Client.Do(r)
	if err != nil {
		return manifest.Coauthor{}, err
	}

	var user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return manifest.Coauthor{}, err
	}
	if user.Login == "" {
		return manifest.Coauthor{}, fmt.Errorf("no GitHub user found for %q", username)
	}
	return manifest.Coauthor{
		Email: fmt.Sprintf("%d+%s@users.noreply.github.com", user.ID, user.Login),
		ID:    user.Login,
		Name:  user.Name,
		Type:  manifest.CoauthorTypeGitHub,
	}, nil
}
