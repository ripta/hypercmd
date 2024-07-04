package hypercmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func InjectInstaller(hc *HyperCommand) *cobra.Command {
	opts := &installOptions{
		hc: hc,
	}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install hyperbinary commands as symlinks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run()
		},
	}

	hc.Root().AddCommand(cmd)
	return cmd
}

type installOptions struct {
	hc  *HyperCommand
	yes bool
}

func (opts *installOptions) run() error {
	target, err := os.Executable()
	if err != nil {
		return err
	}

	dir := path.Dir(target)
	fmt.Fprintf(os.Stderr, "Installing %d symlinks to %s in %s\n", len(opts.hc.cmds), target, dir)

	aggErr := []error{}
	for _, cmd := range opts.hc.cmds {
		ln := path.Join(dir, cmd.Name())
		if _, err := os.Stat(ln); err == nil {
			fmt.Fprintf(os.Stderr, "Skip: symlink for %s already exists at %s\n", cmd.Name(), ln)
		}

		if !opts.yes {
			fmt.Fprintf(os.Stderr, "Dry-run: would have installed symlink for %s at %s\n", cmd.Name(), ln)
			continue
		}

		if err := os.Symlink(target, ln); err != nil {
			aggErr = append(aggErr, err)
		}
		fmt.Fprintf(os.Stderr, "Installed symlink for %s at %s\n", cmd.Name(), ln)
	}

	if err := errors.Join(aggErr...); err != nil {
		return fmt.Errorf("could not install one or more symlinks: %w", err)
	}
	return nil
}
