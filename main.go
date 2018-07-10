package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	service, err := NewService(os.Getenv("MESG_ENDPOINT"), os.Getenv("MESG_TOKEN"))
	if err != nil {
		panic(err)
	}

	// Start the HTTP server
	go func() {
		if err := NewServer(os.Getenv("HTTP_SERVER"), service).ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// Start service task listener
	go func() {
		if err := service.ListenTask(); err != nil {
			panic(err)
		}
	}()

	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	<-abort
}
