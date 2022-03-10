package command

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestManifestManagement(t *testing.T) {
	cmd := New(newWorkspace(t))

	// Add a coauthor
	err := cmd.ManifestAdd("brett", "Brett Buddin", "brett@buddin.org")
	require.NoError(t, err)

	// List the coauthors
	out := bytes.NewBuffer(nil)
	err = cmd.ManifestList(out)
	require.NoError(t, err)
	require.Equal(t, listExample(`
ID     NAME          EMAIL             TYPE
brett  Brett Buddin  brett@buddin.org  manual
`), out.String())

	// Remove the coauthor
	err = cmd.ManifestRemove("brett")
	require.NoError(t, err)

	// Verify that the coauthor was removed
	out.Truncate(0)
	err = cmd.ManifestList(out)
	require.NoError(t, err)
	require.Equal(t, "", out.String())
}

func TestActivationWorkflow(t *testing.T) {
	cmd := New(newWorkspace(t))

	// Add and activate the first coauthor
	err := cmd.ManifestAdd("brett", "Brett Buddin", "brett@buddin.org")
	require.NoError(t, err)
	err = cmd.TemplateSet("brett")
	require.NoError(t, err)

	repoPaths, err := cmd.Paths.Repository()
	require.NoError(t, err)

	// Verify the template contains what we'd expect
	tmplb, err := ioutil.ReadFile(repoPaths.TemplateFile)
	require.NoError(t, err)
	require.Equal(t, templateExample(`
# Managed by partner
#
# partner-id: brett
Co-Authored-By: "Brett Buddin" <brett@buddin.org>
`), string(tmplb))

	// Add and activate the second coauthor
	err = cmd.ManifestAdd("persona", "Person A", "a@buddin.org")
	require.NoError(t, err)
	err = cmd.TemplateSet("persona")
	require.NoError(t, err)

	// List the active coauthors
	out := bytes.NewBuffer(nil)
	err = cmd.TemplateStatus(out)
	require.NoError(t, err)
	require.Equal(t, listExample(`
ID       NAME          EMAIL             TYPE
brett    Brett Buddin  brett@buddin.org  manual
persona  Person A      a@buddin.org      manual
`), out.String())

	// Verify the template contains what we'd expect
	tmplb, err = ioutil.ReadFile(repoPaths.TemplateFile)
	require.NoError(t, err)
	require.Equal(t, templateExample(`
# Managed by partner
#
# partner-id: brett
Co-Authored-By: "Brett Buddin" <brett@buddin.org>
# partner-id: persona
Co-Authored-By: "Person A" <a@buddin.org>
`), string(tmplb))

	// Unset all active coauthors, and verify that the tool reports nothing
	err = cmd.TemplateClear()
	require.NoError(t, err)
	out.Truncate(0)
	err = cmd.TemplateStatus(out)
	require.NoError(t, err)
	require.Equal(t, "", out.String())

	// Verify the template file is deleted
	_, err = os.Stat(repoPaths.TemplateFile)
	require.Error(t, err)
}

func newWorkspace(t *testing.T) Paths {
	t.Helper()

	tmp, err := ioutil.TempDir("", "partner_test")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmp)
	})

	absTmp, err := filepath.Abs(tmp)
	require.NoError(t, err)

	cmd := exec.Command("git", "init")
	cmd.Dir = absTmp
	err = cmd.Run()
	require.NoError(t, err)

	return Paths{
		WorkDir:      absTmp,
		ManifestFile: filepath.Join(absTmp, "manifest.json"),
	}
}

func TestDefaultPaths_RepositoryPathCalculation(t *testing.T) {
	wsPaths := newWorkspace(t)

	repoPaths, err := wsPaths.Repository()
	require.NoError(t, err)

	// `git rev-parse --show-toplevel` returns a different path on macOS,
	// because /tmp points at /private/tmp.
	repoPaths.Root = strings.TrimPrefix(repoPaths.Root, "/private")
	repoPaths.TemplateFile = strings.TrimPrefix(repoPaths.TemplateFile, "/private")

	require.Equal(t, wsPaths.WorkDir, repoPaths.Root)
	require.Equal(t, filepath.Join(wsPaths.WorkDir, ".git/gitmessage.txt"), repoPaths.TemplateFile)
}

func TestDefaultPath_PathExpansion(t *testing.T) {
	t.Run("default manifest path", func(t *testing.T) {
		paths, err := DefaultPaths(".")
		require.NoError(t, err)
		require.NotEqual(t, "~/.config/partner/manifest.json", paths.ManifestFile)
		require.True(t, strings.HasSuffix(paths.ManifestFile, "/.config/partner/manifest.json"), "tilde was not expanded to home directory")
	})

	t.Run("overridden manifest path", func(t *testing.T) {
		os.Setenv("PARTNER_MANIFEST", "~/other/path/manifest.json")
		defer os.Unsetenv("PARTNER_MANIFEST")

		paths, err := DefaultPaths(".")
		require.NoError(t, err)
		require.NotEqual(t, "~/other/path/manifest.json", paths.ManifestFile)
		require.True(t, strings.HasSuffix(paths.ManifestFile, "/other/path/manifest.json"), "tilde was not expanded to home directory")
	})

	t.Run("overridden manifest path with environment variable", func(t *testing.T) {
		os.Setenv("PARTNER_MANIFEST", "$HOME/.config/partner/manifest.json")
		defer os.Unsetenv("PARTNER_MANIFEST")

		paths, err := DefaultPaths(".")
		require.NoError(t, err)
		require.NotEqual(t, "$HOME/.config/partner/manifest.json", paths.ManifestFile)
		require.True(t, strings.HasSuffix(paths.ManifestFile, "/.config/partner/manifest.json"), "environment variable was not expanded to home directory")
	})
}

func listExample(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

func templateExample(s string) string {
	return "\n" + s
}
