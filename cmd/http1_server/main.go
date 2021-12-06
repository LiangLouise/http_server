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
	if len(args) != 4 {
		fmt.Println("Please give host, port to bind along with the protocol version (1.0/1.1)")
		return
	}

	fs, err := fsService.MakeFsService()
	if err != nil {
		log.Fatalf("file system: %s\n", err)
	}

	var protocol httpProto.HTTP_PROTOCOL_VERSION
	if args[3] == "1.0" {
		protocol = httpProto.HTTP_1
	} else if args[3] == "1.1" {
		protocol = httpProto.HTTP_1_1
	} else {
		fmt.Println("Please enter correct protocol version (1.0 or 1.1)")
		return
	}
	s, err := server.MakeServer(args[1], args[2], protocol, fs)
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
