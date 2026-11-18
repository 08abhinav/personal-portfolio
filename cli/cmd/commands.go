package cmd

import (
	"fmt"

	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/spf13/cobra"
)

func newCommandsCmd(e *env.Env) *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "List everything you can run",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			var list []*cobra.Command
			width := 0
			for _, c := range cmd.Root().Commands() {
				if !c.IsAvailableCommand() { // skips hidden commands and Cobra's built-in help
					continue
				}
				list = append(list, c)
				if len(c.Name()) > width {
					width = len(c.Name())
				}
			}

			for _, c := range list {
				fmt.Fprintf(out, "  %-*s  %s\n", width, c.Name(), c.Short)
			}
			fmt.Fprintln(out, "\ntip: add -h to any command to see its options")
		},
	}
}
