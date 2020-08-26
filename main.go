package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MattIzSpooky/tf2.rest/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGHUP,
	)

	httpServer := server.NewServer()

	go func() {
		if err := httpServer.Serve(); err != nil {
			panic(err)
		}
	}()

	<-sigChan

	if err := httpServer.GracefulShutdown(); err != nil {
		panic(err)
	}
}
