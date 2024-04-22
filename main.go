package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	command "github.com/copito/goscaffold/cmd"
)

func main() {
	// Force Stop at any point (with termination signals)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	go func(sig chan os.Signal) {
		<-sig
		fmt.Println("Forced exited the cli...")
		os.Exit(1)
	}(sig)

	// Start CLI
	command.Execute()
}
