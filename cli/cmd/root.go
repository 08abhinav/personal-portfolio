package cmd

import (
	"fmt"
	"os"

	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableCommandSorting = false
}

func NewRootCmd(e *env.Env) *cobra.Command {
	root := &cobra.Command{
		Use:           "abhinav",
		Short:         "Abhinav Negi's portfolio, in your terminal",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "Hi, I'm Abhinav. Try: abhinav ping, or abhinav commands to see everything.")
		},
	}

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(
		newPingCmd(e),
		newWhoamiCmd(e),
		newSocialsCmd(e),
		newWhereamiCmd(e),
		newCommandsCmd(e),
		newClearCmd(e),
	)
	return root
}

func Execute() {
	root := NewRootCmd(env.Local())
	root.AddCommand(newServeCmd())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
