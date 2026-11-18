package cmd

import (
	"fmt"
	"time"

	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/spf13/cobra"
)

func newPingCmd(e *env.Env) *cobra.Command{
	return &cobra.Command{
		Use: "ping",
		Short: "Say ping, get a pong",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			if e.Mode == env.ModeWeb && e.Client.Timezone != "" {
				if loc, err := time.LoadLocation(e.Client.Timezone); err == nil {
					now = now.In(loc)
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "pong %s\n", now.Format("2006-01-02 15:04:05 MST"))
		},
	}
}
