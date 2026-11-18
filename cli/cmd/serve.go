package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/08abhinav/personal-portfolio/cli/internal/env"
	"github.com/08abhinav/personal-portfolio/cli/internal/geo"
	"github.com/spf13/cobra"
)

type runRequest struct {
	Line     string `json:"line"`
	Timezone string `json:"tz"`
}

type runResponse struct {
	Output string `json:"output"`
	Error  bool   `json:"error,omitempty"`
}

func newServeCmd() *cobra.Command {
	var addr, cityDB, asnDB string

	cmd := &cobra.Command{
		Use:    "serve",
		Short:  "Run the HTTP API behind the web terminal",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var resolver *geo.Resolver
			if cityDB != "" {
				r, err := geo.Open(cityDB, asnDB)
				if err != nil {
					return err
				}
				defer r.Close()
				resolver = r
			} else {
				log.Println("no --geo-city given: whereami will skip the location lookup")
			}

			mux := http.NewServeMux()
			mux.HandleFunc("/api/run", runHandler(resolver))
			mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok\n"))
			})

			srv := &http.Server{
				Addr:              addr,
				Handler:           mux,
				ReadHeaderTimeout: 3 * time.Second,
				ReadTimeout:       5 * time.Second,
				WriteTimeout:      10 * time.Second,
				IdleTimeout:       60 * time.Second,
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			go func() {
				<-ctx.Done()
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				srv.Shutdown(shutdownCtx)
			}()

			log.Printf("listening on %s", addr)
			if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&addr, "addr", "127.0.0.1:8080", "address to listen on (keep it on localhost behind Nginx)")
	cmd.Flags().StringVar(&cityDB, "geo-city", "", "path to GeoLite2-City.mmdb")
	cmd.Flags().StringVar(&asnDB, "geo-asn", "", "path to GeoLite2-ASN.mmdb")
	return cmd
}


func runHandler(resolver *geo.Resolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req runRequest
		switch r.Method {
		case http.MethodGet:
			req.Line = r.URL.Query().Get("cmd")
			req.Timezone = r.URL.Query().Get("tz")
		case http.MethodPost:
			r.Body = http.MaxBytesReader(w, r.Body, 1024)
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				reply(w, r, http.StatusBadRequest, "bad request", true)
				return
			}
		default:
			reply(w, r, http.StatusMethodNotAllowed, "method not allowed", true)
			return
		}

		line := strings.TrimSpace(req.Line)
		args := strings.Fields(line) // split on spaces only; nothing here ever touches a shell
		if len(args) == 0 || len(line) > 100 || len(args) > 6 {
			reply(w, r, http.StatusBadRequest, "type a command, e.g. ping (max 100 characters)", true)
			return
		}

		e := &env.Env{
			Mode:      env.ModeWeb,
			StartedAt: time.Now(),
			Geo:       resolver,
			Client:    env.Client{IP: clientIP(r), UserAgent: r.UserAgent(), Timezone: req.Timezone},
		}

		var buf bytes.Buffer
		root := NewRootCmd(e) // a fresh command tree per request
		root.SetOut(&buf)
		root.SetErr(&buf)
		root.SetArgs(args)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		out, failed := "", false
		if err := root.ExecuteContext(ctx); err != nil {
			failed = true
			if strings.HasPrefix(err.Error(), "unknown command") {
				out = "command not found: " + args[0] + "\ntype `commands` to see what you can run"
			} else {
				out = buf.String() + "error: " + err.Error()
			}
		} else {
			out = buf.String()
		}
		reply(w, r, http.StatusOK, out, failed)
	}
}

func reply(w http.ResponseWriter, r *http.Request, status int, out string, failed bool) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)
		w.Write([]byte(strings.TrimRight(out, "\n") + "\n"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(runResponse{Output: out, Error: failed})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		if x := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(x) != nil {
			return x
		}
	}
	return host
}
