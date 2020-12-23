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
	err := cmd.ManifestManualAdd("brett", "Brett Buddin", "brett@buddin.org")
	require.NoError(t, err)

	out := bytes.NewBuffer(nil)
	err = cmd.ManifestList(out)
	require.NoError(t, err)

	require.Equal(t, trimLeft(`
ID     NAME          EMAIL             TYPE
brett  Brett Buddin  brett@buddin.org  manual
`), out.String())

	err = cmd.ManifestRemove("brett")
	require.NoError(t, err)

	out.Truncate(0)
	err = cmd.ManifestList(out)
	require.NoError(t, err)
	require.Equal(t, "", out.String())
}

func TestActivationWorkflow(t *testing.T) {
	cmd := New(newWorkspace(t))
	err := cmd.ManifestManualAdd("brett", "Brett Buddin", "brett@buddin.org")
	require.NoError(t, err)
	err = cmd.TemplateSet("brett")
	require.NoError(t, err)

	tmplb, err := ioutil.ReadFile(cmd.Paths.TemplateFile)
	require.NoError(t, err)
	require.Equal(t, prefixLine(`
# Managed by partner
#
# partner-id: brett
Co-Authored-By: "Brett Buddin" <brett@buddin.org>
`), string(tmplb))

	err = cmd.ManifestManualAdd("persona", "Person A", "a@buddin.org")
	require.NoError(t, err)
	err = cmd.TemplateSet("persona")
	require.NoError(t, err)

	out := bytes.NewBuffer(nil)
	err = cmd.TemplateStatus(out)
	require.NoError(t, err)
	require.Equal(t, trimLeft(`
ID       NAME          EMAIL             TYPE
brett    Brett Buddin  brett@buddin.org  manual
persona  Person A      a@buddin.org      manual
`), out.String())

	tmplb, err = ioutil.ReadFile(cmd.Paths.TemplateFile)
	require.NoError(t, err)
	require.Equal(t, prefixLine(`
# Managed by partner
#
# partner-id: brett
Co-Authored-By: "Brett Buddin" <brett@buddin.org>
# partner-id: persona
Co-Authored-By: "Person A" <a@buddin.org>
`), string(tmplb))

	err = cmd.TemplateClear()
	require.NoError(t, err)

	out.Truncate(0)
	err = cmd.TemplateStatus(out)
	require.NoError(t, err)
	require.Equal(t, "", out.String())

	_, err = os.Stat(cmd.Paths.TemplateFile)
	require.Error(t, err)
}

func newWorkspace(t *testing.T) Paths {
	t.Helper()

	tmp, err := ioutil.TempDir("", "partner_test")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmp)
	})

	cmd := exec.Command("git", "init")
	cmd.Dir = tmp
	err = cmd.Run()
	require.NoError(t, err)

	return Paths{
		WorkingDir:   tmp,
		ManifestFile: filepath.Join(tmp, "manifest.json"),
		TemplateFile: filepath.Join(tmp, ".git/gitmessage.txt"),
	}
}

func TestDefaultPaths(t *testing.T) {
	wsPaths := newWorkspace(t)

	os.Setenv("PARTNER_MANIFEST", filepath.Join(wsPaths.WorkingDir, "manifest.json"))
	paths, err := DefaultPaths(wsPaths.WorkingDir)

	// `git rev-parse --show-toplevel` returns a different path than I set.
	// Probably a macOS thing. If this test starts to fail, this is probably
	// why.
	paths.TemplateFile = strings.TrimPrefix(paths.TemplateFile, "/private")

	require.NoError(t, err)
	require.Equal(t, Paths{
		WorkingDir:   wsPaths.WorkingDir,
		TemplateFile: wsPaths.TemplateFile,
		ManifestFile: wsPaths.ManifestFile,
	}, paths)
}

func trimLeft(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

func prefixLine(s string) string {
	return "\n" + s
}
