// Read Raw http request message, convert it to object
package httpParser

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"

	"github.com/liangLouise/http_server/pkg/httpProto"
)

type Request struct {
	method   string // GET, POST, etc.
	header   textproto.MIMEHeader
	body     []byte
	uri      string // The raw URI from the request
	protocol httpProto.HTTP_PROTOCOL_VERSION
}

func (req *Request) GetMethod() string {
	return req.method
}

func (req *Request) GetHeader() textproto.MIMEHeader {
	return req.header
}

func (req *Request) GetBody() []byte {
	return req.body
}

func (req *Request) GetUri() string {
	return req.uri
}

func (req *Request) GetProtocol() httpProto.HTTP_PROTOCOL_VERSION {
	return req.protocol
}

func (req *Request) GetConnection() string {
	return req.GetHeader().Get("Connection")
}

func ParseRequest(connection net.Conn) ([]Request, error) {
	var req Request
	var err error
	reqList := make([]Request, 0)
	bufioReader := bufio.NewReader(connection)
	textprotoReader := textproto.NewReader(bufioReader)
	for {
		reqLine, err := textprotoReader.ReadLine()
		if err != nil {
			log.Println(err)
			return reqList, err
		}
		reqLineSplitted := strings.Split(reqLine, " ")

		req = Request{
			method:   reqLineSplitted[0],
			uri:      reqLineSplitted[1],
			protocol: httpProto.HTTP_PROTOCOL_VERSION(reqLineSplitted[2]),
		}

		// req.method, req.uri, req.protocol = reqLineSplitted[0], reqLineSplitted[1], httpProto.HTTP_PROTOCOL_VERSION(reqLineSplitted[2])

		header, err := textprotoReader.ReadMIMEHeader()
		if err != nil {
			log.Println(err)
			return reqList, err
		}
		req.header = header
		log.Printf("Request Method: %s", req.GetMethod())
		log.Printf("Request URI: %s", req.GetUri())
		log.Printf("Request Protocol: %s", req.GetProtocol())
		log.Printf("Request Headers: %s", req.GetHeader())
		if req.GetHeader().Get("Content-Length") != "" {
			length := req.GetHeader().Get("Content-Length")[0]
			body := make([]byte, length)
			_, err = io.ReadFull(bufioReader, body) // ReadAll?
			if err != nil {
				log.Println(err)
				return reqList, err
			}
			req.body = body
			log.Printf("Request Body: %s", req.GetBody())
		}
		log.Printf("This request does not have a body")
		log.Printf("Connection: %s", req.GetConnection())
		reqList = append(reqList, req)
		remain := bufioReader.Buffered()
		if remain == 0 {
			log.Printf("Remaining bytes: %d", remain)
			break
		}

	}

	return reqList, err
}
