package router

import (
	"fmt"
	"log"
	"net"

	"github.com/liangLouise/http_server/pkg/httpParser"
)

type router struct {
}

func SimpleHandler(connection net.Conn) {
	// netData, err := bufio.NewReader(connection).ReadString('\n')
	httpParser.ParseRequest(connection)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// if strings.TrimSpace(netData) == "STOP" {
	// 	fmt.Println("Exiting TCP server!")
	// 	return
	// }

	log.Printf("Address: %s", connection.RemoteAddr().String())

	fmt.Fprintf(connection, "HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"Content-Length: 20\r\n"+
		"\r\n"+
		"<h1>hello world</h1>")

}
