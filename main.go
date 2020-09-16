package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MattIzSpooky/tf2.rest/server"
	"github.com/joho/godotenv"
)

//go:generate go run gen.go

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	//responses.Setup()

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
