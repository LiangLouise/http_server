package main

import (
	"fmt"
	"os"

	"github.com/liangLouise/http_server/pkg/server"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please give host and port to bind")
		return
	}

	s, err := server.MakeServer("", args[1], server.HTTP_1_1)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	s.ListenRequest()

}
