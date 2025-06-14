package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/services"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	localDir := "./local-filesystem"
	remoteDir := "./remote-filesystem"

	synchronizerService := services.NewSynchronizerService(localDir, remoteDir)

	synchronizerService.Start()

	// --- Graceful Shutdown ---
	// Listen for OS signals to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down...", sig)

	synchronizerService.Stop()
}
