package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please give host and port to bind")
		return
	}

	port := ":" + args[1]
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(netData) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Fprintf(c, "HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n"+
			"Content-Length: 20\r\n"+
			"\r\n"+
			"<h1>hello world</h1>")
	}
}
