package hypercmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// HyperCommand is a container of a root (the hypercommand command itself) and
// all commands under it. It is special from a regular command in that it can
// explode its commands into symlinked binaries.
type HyperCommand struct {
	root *cobra.Command
	cmds []*cobra.Command
}

// New initializes a new hypercommand to which commands can be added. An
// "install" subcommand is already automatically added.
func New(name string) *HyperCommand {
	h := &HyperCommand{}

	h.root = &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Run a command in %s", name),
	}

	_ = InjectInstaller(h)
	return h
}

// AddCommand adds a new command to the hypercommand. This makes it available
// as a subcommand and as a symlinked binary when installed.
func (h *HyperCommand) AddCommand(c *cobra.Command) {
	h.root.AddCommand(c)
	h.cmds = append(h.cmds, c)
}

// Root returns the root command of the hypercommand. This might be useful
// if you want to add a command that should not be expanded into symlinks.
func (h *HyperCommand) Root() *cobra.Command {
	return h.root
}

// Commands returns all the commands registered in the hypercommand. This
// may not include all commands if they were added directly to the root.
func (h *HyperCommand) Commands() []*cobra.Command {
	return h.cmds
}

// ImportCommands adds all subcommands of an existing command to the hypercommand.
func (h *HyperCommand) ImportCommands(c *cobra.Command) {
	for _, cmd := range c.Commands() {
		h.AddCommand(cmd)
	}
}

var ErrNoCommand = errors.New("no command found")

// Resolve is given a name of a binary and uses it to return the correct command,
// or the hypercommand otherwise.
func (h *HyperCommand) Resolve(allArgs []string, withAliases bool) (*cobra.Command, error) {
	name := filepath.Base(allArgs[0])
	if h.root.Name() == name {
		return h.root, nil
	}
	if withAliases {
		for _, alias := range h.root.Aliases {
			if alias == name {
				return h.root, nil
			}
		}
	}

	// Reinject the root command name into the arguments, because (cobra.Command).Execute
	// always traverses to the root before execution, even when we want a subcommand
	// (which we do, since we already know it's not the root).
	h.root.SetArgs(allArgs)

	for _, cmd := range h.cmds {
		if cmd.Name() == name {
			return cmd, nil
		}
		if withAliases {
			for _, alias := range cmd.Aliases {
				if alias == name {
					return cmd, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("%w: %s (root is %s)", ErrNoCommand, name, h.root.Name())
}
