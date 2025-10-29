# hypercmd

A [spf13/cobra](https://github.com/spf13/cobra)-compatible wrapper to make magical
hyperbinaries: a binary containing other commands. The hyperbinary behaves like tools
such as `busybox`, where multiple related commands can be bundled into a single binary
for convenience, and then installed as separate symlinked commands.

## Usage

Create a hyperbinary by initializing a `HyperCommand` and adding your cobra commands to it:

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ripta/hypercmd/pkg/hypercmd"
)

func main() {
	// Create a new hypercommand
	hc := hypercmd.New("mytool")

	// Add commands to the hypercommand
	addCmd := &cobra.Command{
		Use:   "add [numbers...]",
		Short: "Add numbers together",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Adding numbers:", args)
		},
	}

	multiplyCmd := &cobra.Command{
		Use:   "multiply [numbers...]",
		Short: "Multiply numbers together",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Multiplying numbers:", args)
		},
	}

	hc.AddCommand(addCmd)
	hc.AddCommand(multiplyCmd)

	// Resolve and execute the appropriate command
	cmd, err := hc.Resolve(os.Args, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

After building the binary, you can use it in multiple ways.

### As a traditional command with subcommands

```bash
# Run as the main command
❯ mytool
Run a command in mytool

Usage: mytool [command]

Available Commands:
  add         Add numbers together
  help        Help about any command
  install     Install hyperbinary commands as symlinks
  multiply    Multiply numbers together

❯ mytool add 1 2 3
Adding numbers: [1 2 3]

❯ mytool multiply 4 5 6
Multiplying numbers: [4 5 6]
```

### As symlinked binaries (after running `install`)

You'd first install the symlinks for each subcommand. Only the first level of
subcommands are actually installed.

```bash
❯ mytool install -y
Installing 2 symlinks to /home/foo/bin/mytool in /home/foo/bin
Installed symlink for add at /home/foo/bin/add
Installed symlink for multiply at /home/foo/bin/multiply
```

Once installed, each subcommand becomes available as its own binary:

```bash
❯ add 1 2 3
Adding numbers: [1 2 3]

❯ multiply 4 5 6
Multiplying numbers: [4 5 6]
```
