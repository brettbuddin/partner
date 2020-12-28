package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/brettbuddin/partner/internal/manifest"
)

const coAuthoredBy = "Co-Authored-By"

var extractPattern = regexp.MustCompile("# partner-id: (.+)")

// ExtractIDs returns the IDs of coauthors referenced in the git commit template
// file.
func ExtractIDs(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var usernames []string
	matches := extractPattern.FindAllSubmatch(b, -1)
	for _, m := range matches {
		usernames = append(usernames, string(m[1]))
	}
	return usernames, nil
}

// Template is a git commit template containing a list of coauthors
type Template struct {
	Coauthors []manifest.Coauthor
}

func (t Template) trailers() string {
	var b strings.Builder
	b.WriteString("\n\n# Managed by partner\n#\n")
	for _, ca := range t.Coauthors {
		fmt.Fprintf(&b, "# partner-id: %s\n%s: %q <%s>\n", ca.ID, coAuthoredBy, ca.Name, ca.Email)
	}
	return b.String()
}

// WriteFile saves and registers the git commit template
func WriteFile(path string, t Template) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open commit template: %w", err)
	}
	defer f.Close()

	if len(t.Coauthors) == 0 {
		f.Truncate(0)
		return nil
	}

	if _, err := f.WriteString(t.trailers()); err != nil {
		return fmt.Errorf("failed to write to commit template file: %w", err)
	}
	return nil
}
