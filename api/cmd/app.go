package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
	"github.com/stevenxie/rgv/api/internal/info"
	"github.com/stevenxie/rgv/api/pkg/api"
	ess "github.com/unixpickle/essentials"
)

// opts are a set of program options.
var opts struct {
	ShowVersion bool
	ShowHelp    bool
	Port        int
}

// Define CLI flags, initialize program.
func init() {
	pflag.BoolVarP(&opts.ShowHelp, "help", "h", false,
		"Show program help (usage).")
	pflag.BoolVarP(&opts.ShowVersion, "version", "v", false,
		"Show program version.")
	pflag.IntVarP(&opts.Port, "port", "p", 3000, "Port to listen on.")

	loadEnv() // load .env variables
	pflag.Parse()
}

// Exec is the entrypoint to command rgv.
func Exec() {
	if opts.ShowHelp {
		pflag.Usage()
		os.Exit(0)
	}
	if opts.ShowVersion {
		fmt.Println(info.Version)
		os.Exit(0)
	}

	// Create program logger.
	logger, err := buildLogger()
	if err != nil {
		ess.Die("Error while building zap.SugaredLogger:", err)
	}

	// Create and run server.
	server, err := api.NewServer(logger)
	if err != nil {
		ess.Die("Error while building server:", err)
	}

	server.Addr = fmt.Sprintf(":%d", opts.Port)
	fmt.Printf("Listening on address '%s'...\n", server.Addr)
	go shutdownUponInterrupt(server)
	err = server.ListenAndServe()
	if (err != nil) && (err != http.ErrServerClosed) {
		ess.Die("Error while starting server:", err)
	}
}

func shutdownUponInterrupt(s *api.Server) {
	const timeout = 1 * time.Second

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)

	<-ch // wait for a signal
	fmt.Printf("Shutting down server gracefully (timeout: %s)...\n",
		timeout.String())
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		ess.Die("Error during server shutdown:", err)
	}
}
