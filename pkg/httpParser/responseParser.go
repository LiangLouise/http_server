// Construct a response object helping adding header, status code, body.
package httpParser

import (
	"net/textproto"
	"strconv"

	p "github.com/liangLouise/http_server/pkg/httpProto"
)

type Response struct {
	statusCode int
	statusText string
	header     textproto.MIMEHeader
	body       []byte
	protocol   p.HTTP_PROTOCOL_VERSION // "HTTP/1.1"
}

func (res *Response) InitHeader() {
	res.header = textproto.MIMEHeader{}
}

func (res *Response) SetStatus(statusCode int, statusText string) {
	res.statusCode = statusCode
	res.statusText = statusText
}

func (res *Response) AddHeader(key, value string) {
	res.header.Add(key, value)
}

func (res *Response) SetHeader(key, value string) {
	res.header.Set(key, value)
}

func (res *Response) SetBody(body []byte) {
	res.body = body
}

func (res *Response) SetProtocol(protocol p.HTTP_PROTOCOL_VERSION) {
	res.protocol = protocol
}

func (res *Response) ParseResponse() string {
	var resText string
	var values string
	resText += string(res.protocol) + " " + strconv.Itoa(res.statusCode) + " " + res.statusText + "\r\n"
	for k, v := range res.header {
		values = ""
		if len(v) > 1 {
			for i := 0; i < len(v); i++ {
				values += v[i]
				if i != len(v)-1 {
					values += "; "
				}
			}
		} else {
			values = v[0]
		}
		resText += k + ": " + values + "\r\n"
	}
	resText += "\r\n" + string(res.body)
	return resText
}

func (res *Response) ConstructRes() {
	res.InitHeader()
	res.SetProtocol(p.HTTP_1_1)
	res.SetStatus(200, "OK")
	res.AddHeader("Content-Type", "text/html")
	res.AddHeader("Content-Type", "charset=utf-8")
	res.AddHeader("Content-Length", "20")
	res.SetBody([]byte("<h1>hello world</h1>"))
}
