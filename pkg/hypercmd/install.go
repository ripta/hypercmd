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
	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "Continue with installation")
	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Overwrite existing files (requires --yes to take effect)")

	hc.Root().AddCommand(cmd)
	return cmd
}

type installOptions struct {
	hc    *HyperCommand
	yes   bool
	force bool
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

		info, statErr := os.Lstat(ln)
		exists := statErr == nil

		if exists && !opts.force {
			fmt.Fprintf(os.Stderr, "Skip: symlink for %s already exists at %s (use -f to overwrite)\n", cmd.Name(), ln)
			continue
		}

		if !opts.yes {
			if exists {
				fmt.Fprintf(os.Stderr, "Dry-run: would have overwritten symlink for %s at %s\n", cmd.Name(), ln)
			} else {
				fmt.Fprintf(os.Stderr, "Dry-run: would have installed symlink for %s at %s\n", cmd.Name(), ln)
			}
			continue
		}

		// opts.yes is true here; if exists, opts.force is also true.
		if exists {
			if info.IsDir() {
				aggErr = append(aggErr, fmt.Errorf("refusing to overwrite directory at %s", ln))
				continue
			}
			if err := os.Remove(ln); err != nil {
				aggErr = append(aggErr, err)
				continue
			}
		}

		if err := os.Symlink(target, ln); err != nil {
			aggErr = append(aggErr, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "Installed symlink for %s at %s\n", cmd.Name(), ln)
	}

	if err := errors.Join(aggErr...); err != nil {
		return fmt.Errorf("could not install one or more symlinks: %w", err)
	}
	if !opts.yes {
		if opts.force {
			fmt.Fprintf(os.Stderr, "Dry-run: use -f -y to overwrite the above symlinks\n")
		} else {
			fmt.Fprintf(os.Stderr, "Dry-run: use -y or --yes to install the above symlinks\n")
		}
	}
	return nil
}
