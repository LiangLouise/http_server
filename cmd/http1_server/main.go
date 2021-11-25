package main

import (
	"fmt"
	"os"

	p "github.com/liangLouise/http_server/pkg/protocols"
	"github.com/liangLouise/http_server/pkg/server"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please give host and port to bind")
		return
	}

	s, err := server.MakeServer("", args[1], p.HTTP_1_0)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	s.ListenRequest()

}
