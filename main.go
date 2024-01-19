package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	kv := NewStore(100, "test.db", "test.idx")
	startServer(kv) // Starts a go routine

	// Create a channel to receive OS signals
	sig := make(chan os.Signal, 1)
	// Notify the `sig` channel on SIGINT or SIGTERM
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sig

	stopServer()
}
