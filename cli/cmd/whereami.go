package cmd

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/08abhinav/personal-portfolio/cli/internal/geo"
	"github.com/mileusna/useragent"
	"github.com/spf13/cobra"
)

func newWhereamiCmd(e *env.Env) *cobra.Command {
	return &cobra.Command{
		Use:   "whereami",
		Short: "What the internet knows about you",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			var rows [][2]string
			var note string
			if e.Mode == env.ModeWeb {
				rows, note = webRows(e)
			} else {
				rows, note = localRows()
			}

			for _, r := range rows {
				fmt.Fprintf(out, "  %-9s %s\n", r[0], r[1])
			}
			fmt.Fprintf(out, "\n%s\n", note)
			return nil
		},
	}
}

func webRows(e *env.Env) ([][2]string, string) {
	c := e.Client
	rows := [][2]string{{"ip", orDash(c.IP)}}
	note := "nothing here is stored."
	tz := c.Timezone

	switch ip := net.ParseIP(c.IP); {
	case ip == nil:
		rows = append(rows, [2]string{"location", "unknown"})
	case ip.IsLoopback() || ip.IsPrivate():
		rows = append(rows, [2]string{"location", "private network address, nothing to look up"})
	case e.Geo == nil:
		rows = append(rows, [2]string{"location", "lookup unavailable right now"})
	default:
		loc, err := e.Geo.Lookup(c.IP)
		if err != nil {
			rows = append(rows, [2]string{"location", "lookup failed"})
			break
		}
		rows = append(rows, locationRows(loc)...)
		if tz == "" {
			tz = loc.Timezone
		}
		note = "location is a guess from your IP and can be off by a city or more (VPNs change it). nothing here is stored."
	}

	if tz != "" {
		rows = append(rows, [2]string{"timezone", tz})
	}

	ua := useragent.Parse(c.UserAgent)
	rows = append(rows,
		[2]string{"browser", strings.TrimSpace(ua.Name + " " + majorVersion(ua.Version))},
		[2]string{"os", strings.TrimSpace(ua.OS + " " + ua.OSVersion)},
		[2]string{"device", deviceKind(ua)},
	)
	return rows, note
}

func locationRows(l *geo.Location) [][2]string {
	var rows [][2]string
	var parts []string
	for _, p := range []string{l.City, l.Region, l.Country} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) > 0 {
		rows = append(rows, [2]string{"location", strings.Join(parts, ", ")})
	}
	if l.HasCoords {
		rows = append(rows, [2]string{"coords", fmt.Sprintf("%.2f, %.2f (approximate)", l.Lat, l.Lon)})
	}
	if l.ASN != 0 {
		rows = append(rows, [2]string{"isp", fmt.Sprintf("AS%d %s", l.ASN, l.Org)})
	}
	return rows
}

func localRows() ([][2]string, string) {
	host, _ := os.Hostname()
	rows := [][2]string{
		{"os", runtime.GOOS},
		{"arch", runtime.GOARCH},
		{"cpus", strconv.Itoa(runtime.NumCPU())},
		{"host", orDash(host)},
		{"timezone", time.Now().Format("MST -07:00")},
		{"go", runtime.Version()},
	}
	return rows, "this is your own machine. open the terminal on abhinavnegi.site for IP and location."
}

func deviceKind(ua useragent.UserAgent) string {
	kind := "unknown"
	switch {
	case ua.Bot:
		kind = "bot"
	case ua.Tablet:
		kind = "tablet"
	case ua.Mobile:
		kind = "mobile"
	case ua.Desktop:
		kind = "desktop"
	}
	if ua.Device != "" {
		kind += " (" + ua.Device + ")"
	}
	return kind
}

func majorVersion(v string) string {
	major, _, _ := strings.Cut(v, ".")
	return major
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
