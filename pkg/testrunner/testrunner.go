package testrunner

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ripta/hypercmd/pkg/hypercmd"
	"github.com/spf13/cobra"
)

func RunCode() int {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
		return 1
	}

	return 0
}

func Run() error {
	root := hypercmd.New("testrunner")
	root.Root().SilenceErrors = true
	root.Root().SilenceUsage = true

	add := &cobra.Command{
		Use:   "add",
		Short: "Add numbers together",

		RunE: func(cmd *cobra.Command, args []string) error {
			nums := []int{}
			for _, arg := range args {
				num, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}

				nums = append(nums, num)
			}

			sum := 0
			for _, n := range nums {
				sum += n
			}

			cmd.Println("Sum:", sum)
			return nil
		},
	}
	root.AddCommand(add)

	multiply := &cobra.Command{
		Use:   "multiply",
		Short: "Multiply numbers together",
		RunE: func(cmd *cobra.Command, args []string) error {
			nums := []int{}
			for _, arg := range args {
				num, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}

				nums = append(nums, num)
			}

			product := 1
			for _, n := range nums {
				product *= n
			}

			cmd.Println("Product:", product)
			return nil
		},
	}
	root.AddCommand(multiply)

	version := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of testrunner",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("testrunner version 1.0")
			return nil
		},
	}
	root.Root().AddCommand(version)
	for _, cmd := range root.Commands() {
		cmd.AddCommand(version)
	}

	cmd, err := root.Resolve(os.Args, true)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
