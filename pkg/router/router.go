package router

import (
	"fmt"
	"log"
	"math"
	"net"
	"regexp"
	"time"

	"github.com/liangLouise/http_server/pkg/httpParser"
)

type router struct {
}

func SimpleHandler(close chan interface{}, connection net.Conn) {
	defer connection.Close()

	for {
		select {
		case <-close:
			return
		default:
		}

		reqs := httpParser.ParseRequest(connection)
		maxReq := math.Min(float64(len(reqs)), 5)
		reqs = reqs[:int(maxReq)]
		for _, req := range reqs {
			log.Printf("Address: %s", connection.RemoteAddr())

			var res httpParser.Response
			res.ConstructRes()

			// fmt.Fprintf(connection, "HTTP/1.1 200 OK\r\n"+
			// 	"Content-Type: text/html; charset=utf-8\r\n"+
			// 	"Content-Length: 20\r\n"+
			// 	"\r\n"+
			// 	"<h1>hello world</h1>")

			// HTTP/1.1 keep connection alive unless specified or timeouted
			regex := regexp.MustCompile("(?i)keep-alive")
			match := regex.Match([]byte(req.GetConnection()))
			if !match {
				fmt.Fprintf(connection, "%s", res.ParseResponse())
				log.Printf("closing the connection %s", connection.RemoteAddr())
				return
			} else {
				res.AddHeader("Keep-Alive", "timeout=5")
				res.AddHeader("Keep-Alive", "max=5")
				// timeout := time.Duration(5) * (time.Second)
				// err := connection.SetDeadline(time.Now().Add(timeout))
				// if err != nil {
				// 	fmt.Println(err)
				// 	return
				// }
				res.SetHeader("Last-Modified", time.Now().Format("01-02-2006 15:04:05"))
				fmt.Fprintf(connection, "%s", res.ParseResponse())
			}
		}
	}

}
