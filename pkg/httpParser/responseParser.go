// Construct a response object helping adding header, status code, body.
package httpParser

import (
	"net/textproto"

	p "github.com/liangLouise/http_server/pkg/protocols"
)

type Response struct {
	statusCode int
	statusText string
	header     textproto.MIMEHeader
	body       []byte
	protocol   p.HTTP_PROTOCOL_VERSION // "HTTP/1.1"
}

func (res *Response) SetStatus(statusCode int, statusText string) {
	res.statusCode = statusCode
	res.statusText = statusText
}

func (res *Response) SetHeader(key, value string) {
	res.header.Add(key, value)
}

func (res *Response) SetBody(body []byte) {
	res.body = body
}

func (res *Response) SetProtocol(protocol p.HTTP_PROTOCOL_VERSION) {
	res.protocol = protocol
}

func (res *Response) ParseResponse() []byte {
	var resText string
	var values string
	resText += string(res.protocol) + " " + string(res.statusCode) + " " + res.statusText + "\r\n"
	for k, v := range res.header {
		if len(v) > 1 {
			for i := 0; i < len(v); i++ {
				values += v[i]
				if i != len(v)-1 {
					values += ";"
				}
			}
		}
		values := v[0]
		resText += k + ": " + values + "\r\n"
	}
	resText += "\r\n" + string(res.body)
	return []byte(resText)
}
