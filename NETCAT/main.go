package main

import (
	"fmt"
	"os"
	"os/signal"

	"net-cat/internal/chat"
	"net-cat/internal/helpers"
)

func main() {
	port := helpers.HandleArgs()
	if port == "" {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	server := chat.NewServer()

	// Create a channel to receive OS signals, to gracefully shut down the server
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt)

	go server.Start(port, shutdownSignal)
	if <-shutdownSignal != os.Kill {
		server.Stop()
	}
}
