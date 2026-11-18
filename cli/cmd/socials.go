package cmd

import (
	"fmt"
	"strings"

	"github.com/08abhinav/personal-portfolio/cli/internal/browser"
	"github.com/08abhinav/personal-portfolio/cli/internal/data"
	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/spf13/cobra"
)

func newSocialsCmd(e *env.Env) *cobra.Command {
	var open string
	names := strings.Join(data.SocialNames(), ", ")

	cmd := &cobra.Command{
		Use:   "socials",
		Short: "My links (socials --open linkedin)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			if open != "" {
				s, ok := data.FindSocial(open)
				if !ok {
					return fmt.Errorf("unknown link %q (try: %s)", open, names)
				}
				
				if e.Mode == env.ModeWeb {
					fmt.Fprintf(out, "open: %s\n", s.URL)
					return nil
				}
				fmt.Fprintf(out, "opening %s ...\n", s.URL)
				return browser.Open(s.URL)
			}

			for _, s := range data.Socials {
				fmt.Fprintf(out, "  %-9s %s\n", s.Name, s.URL)
			}
			fmt.Fprintln(out, "\ntip: socials --open linkedin")
			return nil
		},
	}

	cmd.Flags().StringVarP(&open, "open", "o", "", "open a link in your browser ("+names+")")
	return cmd
}
