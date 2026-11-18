package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/08abhinav/personal-portfolio/cli/internal/data"
	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/spf13/cobra"
)

func newWhoamiCmd(e *env.Env) *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Who I am, on one screen",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			me := data.Me

			if asJSON {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(me)
			}

			fmt.Fprintf(out, "%s\n%s\n\n", me.Name, me.Title)

			rows := [][2]string{
				{"role", me.Role},
				{"focus", strings.Join(me.Focus, " · ")},
				{"status", me.Status},
				{"site", me.Site},
			}
			for _, r := range rows {
				fmt.Fprintf(out, "  %-8s %s\n", r[0], r[1])
			}

			fmt.Fprintf(out, "\nnext: try `%s`\n", me.Next)
			return nil
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "print as JSON")
	return cmd
}
