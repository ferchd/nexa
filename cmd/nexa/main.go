package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferchd/nexa/internal/checker"
	"github.com/ferchd/nexa/internal/config"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		log.Printf("Nexa version %s, commit %s, built at %s", version, commit, date)
		os.Exit(0)
	}

	nexa, err := checker.NewNexa(cfg)
	if err != nil {
		log.Fatalf("Error creating checker: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		nexa.Shutdown()
	}()

	result := nexa.Run()

	if cfg.StdoutJSON {
		result.PrintJSON()
	} else {
		result.PrintHuman()
	}

	os.Exit(result.ExitCode())
}