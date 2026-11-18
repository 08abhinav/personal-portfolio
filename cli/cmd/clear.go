package cmd

import (
	"fmt"

	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/spf13/cobra"
)

func newClearCmd(e *env.Env) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear the screen",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStdout(), "\x1b[2J\x1b[3J\x1b[H")
		},
	}
}
