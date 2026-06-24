package hypercmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newInstallTarget(t *testing.T) (*HyperCommand, string, string) {
	t.Helper()

	h := New("installer_test")
	h.AddCommand(&cobra.Command{Use: "add"})
	h.AddCommand(&cobra.Command{Use: "multiply"})

	dir := t.TempDir()
	target := filepath.Join(dir, "installer_test")
	return h, target, dir
}

func readlink(t *testing.T, dir, name string) string {
	t.Helper()

	dest, err := os.Readlink(filepath.Join(dir, name))
	require.NoError(t, err)
	return dest
}

type installTest struct {
	label  string
	yes    bool
	force  bool
	setup  func(t *testing.T, dir, target string)
	expErr string
	expOut []string
	check  func(t *testing.T, dir, target string)
}

var installTests = []installTest{
	{
		label: "dry-run on clean directory installs nothing",
		expOut: []string{
			"Dry-run: would have installed symlink for add",
			"Dry-run: would have installed symlink for multiply",
			"Dry-run: use -y or --yes to install",
		},
		check: func(t *testing.T, dir, target string) {
			_, err := os.Lstat(filepath.Join(dir, "add"))
			assert.True(t, os.IsNotExist(err), "no symlink should be created on a dry-run")
		},
	},
	{
		label: "yes on clean directory installs symlinks",
		yes:   true,
		expOut: []string{
			"Installed symlink for add",
			"Installed symlink for multiply",
		},
		check: func(t *testing.T, dir, target string) {
			assert.Equal(t, target, readlink(t, dir, "add"))
			assert.Equal(t, target, readlink(t, dir, "multiply"))
		},
	},
	{
		label: "existing symlink is skipped without force",
		yes:   true,
		setup: func(t *testing.T, dir, target string) {
			require.NoError(t, os.Symlink("/some/other/target", filepath.Join(dir, "add")))
		},
		expOut: []string{
			"Skip: symlink for add already exists",
			"Installed symlink for multiply",
		},
		check: func(t *testing.T, dir, target string) {
			assert.Equal(t, "/some/other/target", readlink(t, dir, "add"))
			assert.Equal(t, target, readlink(t, dir, "multiply"))
		},
	},
	{
		label: "force and yes overwrite an existing file",
		yes:   true,
		force: true,
		setup: func(t *testing.T, dir, target string) {
			require.NoError(t, os.WriteFile(filepath.Join(dir, "add"), []byte("stale"), 0o644))
		},
		expOut: []string{
			"Installed symlink for add",
			"Installed symlink for multiply",
		},
		check: func(t *testing.T, dir, target string) {
			assert.Equal(t, target, readlink(t, dir, "add"))
			assert.Equal(t, target, readlink(t, dir, "multiply"))
		},
	},
	{
		label: "force without yes reports an overwrite dry-run",
		force: true,
		setup: func(t *testing.T, dir, target string) {
			require.NoError(t, os.WriteFile(filepath.Join(dir, "add"), []byte("stale"), 0o644))
		},
		expOut: []string{
			"Dry-run: would have overwritten symlink for add",
			"Dry-run: would have installed symlink for multiply",
			"Dry-run: use -f -y to overwrite",
		},
		check: func(t *testing.T, dir, target string) {
			info, err := os.Lstat(filepath.Join(dir, "add"))
			require.NoError(t, err)
			assert.Zero(t, info.Mode()&os.ModeSymlink, "file should not be replaced on a dry-run")
		},
	},
	{
		label: "refuses to overwrite a directory",
		yes:   true,
		force: true,
		setup: func(t *testing.T, dir, target string) {
			require.NoError(t, os.Mkdir(filepath.Join(dir, "add"), 0o755))
		},
		expErr: "refusing to overwrite directory",
		expOut: []string{
			"Installed symlink for multiply",
		},
		check: func(t *testing.T, dir, target string) {
			info, err := os.Lstat(filepath.Join(dir, "add"))
			require.NoError(t, err)
			assert.True(t, info.IsDir())
			assert.Equal(t, target, readlink(t, dir, "multiply"))
		},
	},
}

func TestInstall(t *testing.T) {
	for _, tt := range installTests {
		t.Run(tt.label, func(t *testing.T) {
			h, target, dir := newInstallTarget(t)
			if tt.setup != nil {
				tt.setup(t, dir, target)
			}

			opts := &installOptions{
				hc:    h,
				yes:   tt.yes,
				force: tt.force,
			}

			buf := bytes.Buffer{}
			err := opts.install(&buf, target, dir)

			if tt.expErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expErr)
			} else {
				require.NoError(t, err)
			}

			for _, want := range tt.expOut {
				assert.Contains(t, buf.String(), want)
			}

			if tt.check != nil {
				tt.check(t, dir, target)
			}
		})
	}
}
