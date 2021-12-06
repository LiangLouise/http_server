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

	go func() {
		s.ListenRequest()
	}()

	<-done
	log.Print("Server Stopped")
	s.ShutDown()
	log.Print("Server Exited Properly")
}
