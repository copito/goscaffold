package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	command "github.com/copito/goscaffold/cmd"
	"github.com/copito/goscaffold/setup"
)

func main() {
	// Setup Logging
	logger := setup.SetupLogging()

	// Setup all configuration for this application
	setup.SetupConfig(logger)

	// Force Stop at any point (with termination signals)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func(sig chan os.Signal) {
		<-sig
		fmt.Println("Forced exited the cli...")
		os.Exit(1)
	}(sig)

	// Start CLI
	command.Execute()
}
