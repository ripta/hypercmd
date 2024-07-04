package hypercmd

import (
	"fmt"
	"os"
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
func New() *HyperCommand {
	h := &HyperCommand{}

	binary := os.Args[0]
	h.root = &cobra.Command{
		Use:   binary,
		Short: fmt.Sprintf("Run a command in %s", binary),
	}

	_ = InjectInstaller(h)
	return h
}

// AddCommand adds a new command to the hypercommand.
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

// Resolve is given a name of a binary and uses it to return the correct command,
// or the hypercommand otherwise.
func (h *HyperCommand) Resolve(name string, withAliases bool) *cobra.Command {
	name = filepath.Base(name)
	for _, cmd := range h.cmds {
		if cmd.Name() == name {
			return cmd
		}
		if withAliases {
			for _, alias := range cmd.Aliases {
				if alias == name {
					return cmd
				}
			}
		}
	}
	return h.root
}
