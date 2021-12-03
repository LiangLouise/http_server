package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/liangLouise/http_server/pkg/fsService"
	"github.com/liangLouise/http_server/pkg/httpProto"
	"github.com/liangLouise/http_server/pkg/server"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	args := os.Args
	if len(args) != 3 {
		fmt.Println("Please give host and port to bind")
		return
	}

	fs, err := fsService.MakeFsService()
	if err != nil {
		log.Fatalf("file system: %s\n", err)
	}

	s, err := server.MakeServer(args[1], args[2], httpProto.HTTP_1_1, fs)
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
