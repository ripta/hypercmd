package hypercmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type resolveTest struct {
	label  string
	args   []string
	expCmd string
	expErr error
	expOut string
}

var resolveTests = []resolveTest{
	{
		// Calling the binary name when in PATH should succeed.
		label:  "invoked from PATH",
		args:   []string{"command_test"},
		expCmd: "command_test",
	},
	{
		// Calling the binary name when in current directory should succeed.
		label:  "invoked in current directory",
		args:   []string{"./command_test"},
		expCmd: "command_test",
	},
	{
		// Calling the binary name in a relative directory should succeed.
		label:  "invoked with relative path",
		args:   []string{"../bin/command_test"},
		expCmd: "command_test",
	},
	{
		// Calling the binary name with absolute path should succeed.
		label:  "invoked with absolute path",
		args:   []string{"/usr/local/bin/command_test"},
		expCmd: "command_test",
	},

	{
		// Calling a subcommand as a subcommand of the main binary should succeed.
		label:  "invoked from PATH with subcommand 'add'",
		args:   []string{"command_test", "add"},
		expCmd: "command_test",
	},
	{
		// Calling a subcommand as a subcommand of the main binary should succeed,
		// even when invoked from current directory.
		label:  "invoked in current directory with subcommand 'multiply'",
		args:   []string{"./command_test", "multiply"},
		expCmd: "command_test",
	},
	{
		// Calling a subcommand as a subcommand of the main binary should succeed,
		// even when invoked from a relative path.
		label:  "invoked with relative path with subcommand 'add'",
		args:   []string{"../bin/command_test", "add"},
		expCmd: "command_test",
	},
	{
		// Calling a subcommand as a subcommand of the main binary should succeed,
		// even when invoked from an absolute path.
		label:  "invoked with absolute path with subcommand 'multiply'",
		args:   []string{"/usr/local/bin/command_test", "multiply"},
		expCmd: "command_test",
	},

	{
		// This is an error because the binary name does not match any known
		// command or subcommand.
		label:  "invoked with unknown command",
		args:   []string{"unknown_command"},
		expErr: ErrNoCommand,
	},

	{
		// Invoking an unknown subcommand should still resolve to the main command,
		// because it doesn't cause a resolution error, but rather a runtime error.
		label:  "invoked with unknown subcommand",
		args:   []string{"command_test", "unknown_subcommand"},
		expCmd: "command_test",
		expOut: "", // No output expected since it's a runtime error
	},

	{
		// Invoking add as a symlinked binary should resolve to the add command.
		label:  "invoked as 'add' symlink",
		args:   []string{"add"},
		expCmd: "add",
		expOut: "adding numbers",
	},
	{
		// Invoking multiply as a symlinked binary should resolve to the multiply command.
		label:  "invoked as 'multiply' symlink",
		args:   []string{"multiply"},
		expCmd: "multiply",
		expOut: "multiplying numbers",
	},
}

func TestResolve(t *testing.T) {
	root := New("command_test")
	root.Root().SilenceErrors = true
	root.Root().SilenceUsage = true

	add := &cobra.Command{
		Use:   "add",
		Short: "Add numbers together",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println("adding numbers")
			return nil
		},
	}
	root.AddCommand(add)

	multiply := &cobra.Command{
		Use:   "multiply",
		Short: "Multiply numbers together",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println("multiplying numbers")
			return nil
		},
	}
	root.AddCommand(multiply)

	version := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of testrunner",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("command_test version 1.0")
			return nil
		},
	}
	root.Root().AddCommand(version)
	for _, cmd := range root.Commands() {
		cmd.AddCommand(version)
	}

	for _, tt := range resolveTests {
		t.Run(tt.label, func(t *testing.T) {
			buf := bytes.Buffer{}
			root.Root().SetOut(&buf)

			cmd, err := root.Resolve(tt.args, true)
			if err != nil {
				assert.ErrorIs(t, err, tt.expErr)
				return
			}

			assert.Equal(t, tt.expCmd, cmd.Name())

			if tt.expOut == "" {
				return
			}

			if err := cmd.Execute(); !assert.NoError(t, err) {
				return
			}

			assert.Contains(t, buf.String(), tt.expOut)
		})
	}
}
