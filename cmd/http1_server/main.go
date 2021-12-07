package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/liangLouise/http_server/pkg/config"
	"github.com/liangLouise/http_server/pkg/fsService"
	"github.com/liangLouise/http_server/pkg/server"
)

func main() {

	// Channel to capture the ctrl+C, the stop signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please give Path to config yaml file")
		return
	}

	fs, err := fsService.MakeFsService()
	if err != nil {
		log.Fatalf("file system: %s\n", err)
	}

	config, err := config.LoadConfig(args[1])
	if err != nil {
		log.Fatalf("loading configs: %s\n", err)
	}

	s, err := server.MakeServer(config, fs)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	log.Printf("Server now listen at: %s:%s\n", s.Address, s.Port)

	// Start to listen new requests
	go s.ListenRequest()

	<-done
	log.Print("Start to stop server")
	s.ShutDown()
	log.Print("Server Exited Properly")
}
