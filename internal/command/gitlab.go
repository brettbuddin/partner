package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/brettbuddin/partner/internal/manifest"
)

type GitLabFetcher struct {
	Client  *http.Client
	BaseURL string
}

func (f *GitLabFetcher) Fetch(username string) (manifest.Coauthor, error) {
	parsed, err := url.Parse(fmt.Sprintf("%s/api/v4/users", f.BaseURL))
	if err != nil {
		return manifest.Coauthor{}, err
	}
	query := parsed.Query()
	query.Set("username", username)
	parsed.RawQuery = query.Encode()

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
		return manifest.Coauthor{}, fmt.Errorf("error fetching %q from GitLab: %s", username, ghError.Message)
	}

	var users []*struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return manifest.Coauthor{}, err
	}
	if len(users) == 0 {
		return manifest.Coauthor{}, fmt.Errorf("error fetching %q from GitLab: username not found", username)
	}
	user := users[0]

	return manifest.Coauthor{
		Email: fmt.Sprintf("%d-%s@users.noreply.gitlab.com", user.ID, user.Username),
		ID:    user.Username,
		Name:  user.Name,
		Type:  manifest.CoauthorTypeGitLab,
	}, nil
}
