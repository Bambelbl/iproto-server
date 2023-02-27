package main

import (
	"github.com/Bambelbl/iproto-server/server"
	"log"
	"os"
	"os/signal"
	"runtime"
)

const (
	MAX_CLIENTS = 100
	SCALE_RPS   = 1000
	LIMIT_RPS   = 100
	ADDR        = ":8080"
	PROCS_COUNT = 4
)

func main() {
	logger := log.New(os.Stdout, "iproto: ", log.LstdFlags)
	runtime.GOMAXPROCS(PROCS_COUNT)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	iprotoServer := server.NewIprotoServer(ADDR, logger, MAX_CLIENTS, SCALE_RPS, LIMIT_RPS)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		if err := iprotoServer.Stop(); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %s", err.Error())
		}
		close(done)
	}()

	iprotoServer.Serve()

	<-done
	logger.Println("Server stopped")
}
