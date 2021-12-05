package router

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"regexp"
	"time"

	"strconv"

	"github.com/liangLouise/http_server/pkg/fsService"
	"github.com/liangLouise/http_server/pkg/httpParser"
	p "github.com/liangLouise/http_server/pkg/httpProto"
)

type router struct {
}

func SimpleHandler(close chan interface{}, connection net.Conn, fs *fsService.FsService) {
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
			uri := req.GetUri()
			if uri == "/" {
				if fs.HasIndex {
					res.SetBody(fs.Cache["index.html"])
				}
			}
			_, file, isDir, err := fs.TryOpen(uri)
			if err != nil {
				fmt.Println(err)
				return
			}
			if isDir {
				length := fs.GetDirLen(file)
				entries := make(chan string, length)
				_, err := fs.WriteDirContent(file, entries)
				if err != nil {
					fmt.Println(err)
					return
				}
				body := "<pre>\r\n"
				body += "<h1>Directory listing for "
				body += uri + "</h1>\r\n<hr>\r\n"
				// files, err := file.ReadDir(-1)
				// if err != nil {
				// 	log.Printf("Error: %s", err)
				// 	return
				// }

				// for _, file := range files {
				// 	fileName := file.Name()
				// 	if file.IsDir() {
				// 		fileName += "/"
				// 	}
				for entry := range entries {
					body += "<a href=\"" + entry + "\">" + entry + "</a>\r\n"
				}

				// }
				body += "</hr>\r\n"
				body += "</pre>\r\n"
				log.Printf("\r\n%s", []byte(body))
				res.SetBody([]byte(body))
				res.InitHeader()
				res.SetProtocol(p.HTTP_1_1)
				res.SetStatus(200, "OK")
				res.AddHeader("Content-Type", "text.html")
				res.AddHeader("Content-Type", "charset=utf-8")
				res.AddHeader("Content-Length", strconv.Itoa(len(body)))
			} else {
				stat, err := file.Stat()
				if err != nil {
					log.Printf("Error: %s", err)
					return
				}
				size := stat.Size()
				fileoutput := make(chan []byte, size)
				_, err = fs.WriteFileContent(file, fileoutput)
				if err != nil {
					log.Printf("Error: %s", err)
					return
				}
				body := <-fileoutput
				res.SetBody(body)
				res.InitHeader()
				res.SetProtocol(p.HTTP_1_1)
				res.SetStatus(200, "OK")
				res.AddHeader("Content-Type", http.DetectContentType([]byte(body)))
				res.AddHeader("Content-Type", "charset=utf-8")
				res.AddHeader("Content-Length", strconv.Itoa(len(body)))

			}

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
