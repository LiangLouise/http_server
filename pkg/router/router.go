// The handler logics for handling incomming HTTP request
//
// Mainly three types of possible responses
//
// - client sends a directory path, and there is no index.html file,
// so server needs to send back the list of the dir content, i.e. `/a/`
//
// - client sends a directory path, and there is an index.html file,
// so server needs to send back the content of index.html
//
// - client sends a file path, server will try to read the file and send
// back file content with correct MIME type
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
	"github.com/liangLouise/http_server/pkg/config"
	"github.com/liangLouise/http_server/pkg/fsService"
	"github.com/liangLouise/http_server/pkg/httpParser"
	p "github.com/liangLouise/http_server/pkg/httpProto"
)

// this is the main logic of the connection handler.
// It will call different helper handler to further handle the request
// based on situation
func SimpleHandler(close chan interface{}, connection net.Conn, fs *fsService.FsService, config *config.ServerConfig) {

	defer connection.Close()

	for {
		select {
		case <-close:
			return
		default:
		}
		// renew timeout timer
		if config.RunTime.HasPersistant {
			timeout := time.Duration(config.RunTime.TimeoutDuration) * (time.Second)
			err := connection.SetDeadline(time.Now().Add(timeout))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		reqs, err := httpParser.ParseRequest(connection)
		// check if anyone error has occur while trying to read from the connection
		if err != nil {
			if os.IsTimeout(err) {
				log.Println("connection timeouted")
			} else if err == io.EOF {
				log.Println("encountered EOF error")
			} else {
				log.Printf("error while parsing request: %s", err)
			}
			return
		}

		var maxReq float64
		if config.Server.Version == p.HTTP_1 || !config.RunTime.HasPipelining {
			maxReq = 1
		} else {
			// limit max pipelined requests as specified in the config file
			maxReq = math.Min(float64(len(reqs)), float64(config.RunTime.MaxPipelining))
		}
		reqs = reqs[:int(maxReq)]
		for _, req := range reqs {
			res, keepOpen := singleReqHandler(connection, req, fs, config.Server.Version)

			// render the response
			fmt.Fprintf(connection, "%s", res.ParseResponse())

			// check if connection should be kept alive
			if !keepOpen {
				return
			}
		}
	}

}

func singleReqHandler(connection net.Conn, req httpParser.Request, fs *fsService.FsService, protocol p.HTTP_PROTOCOL_VERSION) (response httpParser.Response, persistence bool) {
	keepOpen := true
	log.Printf("Address: %s", connection.RemoteAddr())
	var res httpParser.Response
	res.InitHeader()
	res.SetProtocol(protocol)
	// only GET method is allowed
	if req.GetMethod() != "GET" {
		res = onErrorHandler(res, p.METHOD_NOT_ALLOWED_CODE)
	} else {
		// HTTP/1.1 keep connection alive unless specified or timeouted
		regex := regexp.MustCompile("(?i)keep-alive")
		match := regex.Match([]byte(req.GetConnection()))
		if !match || protocol == p.HTTP_1 {
			log.Printf("keep-alive not supported, closing the connection %s", connection.RemoteAddr())
			keepOpen = false
		} else {
			res.AddHeader("Keep-Alive", "timeout=5")
			res.AddHeader("Keep-Alive", "max=5")
			keepOpen = true
		}
		uri := req.GetUri()
		t := req.GetHeader().Get("If-Modified-Since")
		// try to open the file
		basepath, file, isDir, err := fs.TryOpen(uri)
		// file does not exist, call coresponding handler
		if err != nil {
			log.Println(err)
			res = onErrorHandler(res, p.NOT_FOUND_CODE)
		} else {
			// request is asking if file has been modified, call coresponding handler
			if t != "" {
				res = IfModSinceHandler(t, res, file, isDir, basepath, uri)
			} else {
				// call diretory handler
				if isDir {
					res = dirHandler(res, basepath, file, uri)
					// call file handler
				} else {
					res = fileHandler(res, file)
				}
			}

		}
	}

	return res, keepOpen
}

// response handler for directory
// it will try to render index.html first if it exists
// otherwise, serve the directory content as response body
func dirHandler(res httpParser.Response, basePath string, dir *os.File, uri string) (response httpParser.Response) {
	// close directory after use
	defer dir.Close()

	indexFile, err := fsService.TryOpenIndex(basePath)

	// exist but failed to open, special status code is required
	if !os.IsNotExist(err) && err != nil {
		// handle permission deny
		if os.IsPermission(err) {
			res = onErrorHandler(res, p.FORBIDDEN_CODE)
			return res
		}
		log.Printf("DirHandler: %s", err)
		return
	}

	if indexFile == nil {

		body := "<html>\r\n"
		body += "<head>\r\n"
		body += "<title>Directory listing for " + uri + "</title>\r\n"
		body += "</head>\r\n"
		body += "<body>\r\n"
		body += "<h1>Directory listing for " + uri + "</h1>\r\n"
		body += "<hr>\r\n"
		body += "<ul>\r\n"

		files, err := dir.ReadDir(-1)
		if err != nil {
			log.Printf("DirHandler: %s", err)
			return
		}

		// Write the dir entries to output channel
		for _, file := range files {
			fileName := file.Name()
			if file.IsDir() {
				fileName += "/"
			}
			body += "<li><a href=\"" + fileName + "\">" + fileName + "</a></li>\r\n"

		}
		body += "</ul>\r\n"
		body += "<hr>\r\n"
		body += "<body>\r\n"
		body += "</html>\r\n"
		log.Printf("\r\n%s", []byte(body))
		res.SetBody([]byte(body))
		res.SetProtocol(p.HTTP_1_1)
		res.SetStatus(p.OK_CODE)
		res.AddHeader("Content-Type", "text.html")
		res.AddHeader("Content-Type", "charset=utf-8")
		res.AddHeader("Content-Length", strconv.Itoa(len(body)))
		fileinfo, err := os.Stat(dir.Name())
		if err != nil {
			log.Printf("Cannot get file info: %s", err)
		}
		LastModTime := fileinfo.ModTime()
		res.SetHeader("Last-Modified", LastModTime.Format(time.RFC1123))
		return res
	} else {
		return fileHandler(res, indexFile)
	}

}

// response handler for file
// load the file content as response body
func fileHandler(res httpParser.Response, file *os.File) (response httpParser.Response) {
	// close file after use
	defer file.Close()

	// Detect file Mime Type
	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		log.Printf("FileHandler: %s", err)
		res = onErrorHandler(res, p.INTERNAL_SERVER_ERROR_CODE)
		return res
	}
	file.Seek(0, io.SeekStart)

	buf := bytes.NewBuffer(nil)
	size, err := io.Copy(buf, file)
	if err != nil {
		log.Printf("FileHandler: %s", err)
		res = onErrorHandler(res, p.INTERNAL_SERVER_ERROR_CODE)
		return res
	}
	body := buf.Bytes()

	res.SetBody(body)
	res.SetStatus(p.OK_CODE)

	res.AddHeader("Content-Type", mtype.String())
	res.AddHeader("Content-Length", strconv.FormatInt(size, 10))
	fileinfo, err := os.Stat(file.Name())
	if err != nil {
		log.Printf("Cannot get file info: %s", err)
	}
	LastModTime := fileinfo.ModTime()
	res.SetHeader("Last-Modified", LastModTime.Format(time.RFC1123))
	return res
}

// response handler when user is asking if file has been modified
// return 304 if it is not, otherwise call FileHandler/DirHandler
// accordingly to serve the updated content
func IfModSinceHandler(t string, res httpParser.Response, file *os.File, isDir bool, basePath string, uri string) (response httpParser.Response) {
	IfModSince, err := time.Parse(time.RFC1123, t)
	if err != nil {
		log.Printf("Cannot parse date: %s", err)
		res = onErrorHandler(res, p.FORBIDDEN_CODE)
		return res
	}
	fileinfo, err := os.Stat(file.Name())
	if err != nil {
		log.Printf("Cannot get file info: %s", err)
		res = onErrorHandler(res, p.FORBIDDEN_CODE)
		return res
	}
	LastModTime := fileinfo.ModTime()
	updated := LastModTime.After(IfModSince)
	if updated {
		if isDir {
			res = dirHandler(res, basePath, file, uri)
		} else {
			res = fileHandler(res, file)
		}
	} else {
		body := ""
		log.Printf("\r\n%s", []byte(body))
		res.SetBody([]byte(body))
		res.SetStatus(p.NOT_MODIFIED_CODE)
		res.AddHeader("Content-Type", "text.html")
		res.AddHeader("Content-Type", "charset=utf-8")
		res.AddHeader("Content-Length", strconv.Itoa(len(body)))
	}
	return res
}

// response handler when there is an error occuring
//
// 1. when use sent any request other than GET
//
// 2. when user has not access to the file
//
// 3. when interal server error happens
//
// 4. when requested file not found
func onErrorHandler(res httpParser.Response, sCode p.HTTP_STATUS_CODE) (response httpParser.Response) {

	body := "<html>\r\n"
	body += "<head>\r\n"
	body += "<title>Error response</title>\r\n"
	body += "</head>\r\n"
	body += "<h1>Error response</h1>\r\n"
	body += "<p>Error code: " + strconv.Itoa(int(sCode)) + "</p>\r\n"
	res.SetStatus(sCode)
	body += "<p>Message: " + string(res.StatusText) + ".</p>\r\n"
	body += "<body>\r\n"
	body += "</body>\r\n"
	body += "</html>\r\n"
	log.Printf("\r\n%s", []byte(body))
	res.SetBody([]byte(body))
	res.AddHeader("Content-Type", "text.html")
	res.AddHeader("Content-Type", "charset=utf-8")
	res.AddHeader("Content-Length", strconv.Itoa(len(body)))
	return res
}
