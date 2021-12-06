package router

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"regexp"
	"time"

	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/liangLouise/http_server/pkg/fsService"
	"github.com/liangLouise/http_server/pkg/httpParser"
	p "github.com/liangLouise/http_server/pkg/httpProto"
)

type router struct {
}

func SimpleHandler(close chan interface{}, connection net.Conn, fs *fsService.FsService) {
	keepOpen := true

	defer connection.Close()

	for {
		select {
		case <-close:
			return
		default:
		}

		timeout := time.Duration(5) * (time.Second)
		err := connection.SetDeadline(time.Now().Add(timeout))
		if err != nil {
			fmt.Println(err)
			return
		}
		reqs, err := httpParser.ParseRequest(connection)
		if err != nil {
			log.Printf("error while parsing request: %s", err)
			return
		}
		maxReq := math.Min(float64(len(reqs)), 5)
		reqs = reqs[:int(maxReq)]
		for _, req := range reqs {
			log.Printf("Address: %s", connection.RemoteAddr())
			var res httpParser.Response
			res.InitHeader()
			// HTTP/1.1 keep connection alive unless specified or timeouted
			regex := regexp.MustCompile("(?i)keep-alive")
			match := regex.Match([]byte(req.GetConnection()))
			if !match {
				log.Printf("closing the connection %s", connection.RemoteAddr())
				keepOpen = false
			} else {
				res.AddHeader("Keep-Alive", "timeout=5")
				res.AddHeader("Keep-Alive", "max=5")
				// res.SetHeader("Last-Modified", time.Now().Format("01-02-2006 15:04:05"))
				keepOpen = true
			}

			uri := req.GetUri()
			time := req.GetHeader().Get("If-Modified-Since")
			if uri == "/" {
				if fs.HasIndex {
					res.SetBody(fs.Cache["index.html"])
				}
			}
			_, file, isDir, err := fs.TryOpen(uri)
			if err != nil {
				log.Println(err)
				res = NotFoundHadnler(res)
			} else {
				if time != "" {
					res = IfModSinceHandler(res, file, isDir)
				} else {
					if isDir {
						res = DirHandler(res, fs, file, uri)
					} else {
						res = FileHandler(res, file)
					}
				}

			}

			fmt.Fprintf(connection, "%s", res.ParseResponse())

			if !keepOpen {
				return
			}
		}
	}

}

func DirHandler(res httpParser.Response, fs *fsService.FsService, dir *os.File, uri string) (response httpParser.Response) {
	entries := make(chan string)
	_, err := fs.WriteDirContent(dir, entries)
	if err != nil {
		log.Println(err)
		return
	}
	body := "<html>\r\n"
	body += "<head>\r\n"
	body += "<title>Directory listing for " + uri + "</title>\r\n"
	body += "</head>\r\n"
	body += "<body>\r\n"
	body += "<h1>Directory listing for " + uri + "</h1>\r\n"
	body += "<hr>\r\n"
	body += "<ul>\r\n"
	for entry := range entries {
		body += "<li><a href=\"" + entry + "\">" + entry + "</a></li>\r\n"
	}
	body += "</ul>\r\n"
	body += "</hr>\r\n"
	body += "</body>\r\n"
	body += "</html>\r\n"
	log.Printf("\r\n%s", []byte(body))
	res.SetBody([]byte(body))
	res.SetProtocol(p.HTTP_1_1)
	res.SetStatus(200, "OK")
	res.AddHeader("Content-Type", "text.html")
	res.AddHeader("Content-Type", "charset=utf-8")
	res.AddHeader("Content-Length", strconv.Itoa(len(body)))
	return res
}

func FileHandler(res httpParser.Response, file *os.File) (response httpParser.Response) {

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	size, err := io.Copy(buf, file)
	if err != nil {
		log.Printf("FileHandler: %s", err)
	}
	body := buf.Bytes()

	res.SetBody(body)
	res.SetProtocol(p.HTTP_1_1)
	res.SetStatus(200, "OK")
	mtype := mimetype.Detect(body[:512])
	log.Printf("%s %s", "Content-Type", mtype.String())

	res.AddHeader("Content-Type", mtype.String())
	// res.AddHeader("Content-Type", "charset=utf-8")
	res.AddHeader("Content-Length", strconv.FormatInt(size, 10))
	return res
}

func NotFoundHadnler(res httpParser.Response) (response httpParser.Response) {
	body := "<html>\r\n"
	body += "<head>\r\n"
	body += "<title>Error response</title>\r\n"
	body += "</head>\r\n"

	body += "<h1>Error response</h1>\r\n"
	body += "<p>Error code: 404</p>\r\n"
	body += "<p>Message: File not found.</p>\r\n"
	body += "<body>\r\n"
	body += "</body>\r\n"
	body += "</html>\r\n"
	log.Printf("\r\n%s", []byte(body))
	res.SetBody([]byte(body))
	res.SetProtocol(p.HTTP_1_1)
	res.SetStatus(404, "File not found")
	res.AddHeader("Content-Type", "text.html")
	res.AddHeader("Content-Type", "charset=utf-8")
	res.AddHeader("Content-Length", strconv.Itoa(len(body)))
	return res
}

func IfModSinceHandler(res httpParser.Response, file *os.File, isDir bool) (response httpParser.Response) {
	return res
}
