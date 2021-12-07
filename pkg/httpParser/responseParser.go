// Construct a response object helping adding header, status code, body.
package httpParser

import (
	"net/textproto"
	"strconv"

	p "github.com/liangLouise/http_server/pkg/httpProto"
)

type Response struct {
	statusCode p.HTTP_STATUS_CODE
	StatusText p.HTTP_STATUS_TEXT
	header     textproto.MIMEHeader
	body       []byte
	protocol   p.HTTP_PROTOCOL_VERSION // "HTTP/1.1"
}

func (res *Response) InitHeader() {
	res.header = textproto.MIMEHeader{}
}

func (res *Response) SetStatus(statusCode p.HTTP_STATUS_CODE) {
	var sText p.HTTP_STATUS_TEXT
	switch statusCode {
	case p.OK_CODE:
		sText = p.OK_TEXT
	case p.NOT_MODIFIED_CODE:
		sText = p.NOT_MODIFIED_TEXT
	case p.FORBIDDEN_CODE:
		sText = p.FORBIDDEN_TEXT
	case p.NOT_FOUND_CODE:
		sText = p.NOT_FOUND_TEXT
	case p.METHOD_NOT_ALLOWED_CODE:
		sText = p.METHOD_NOT_ALLOWED_TEXT
	default:
		statusCode = p.INTERNAL_SERVER_ERROR_CODE
		sText = p.INTERNAL_SERVER_ERROR_TEXT
	}
	res.statusCode = statusCode
	res.StatusText = sText
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
	resText += string(res.protocol) + " " + strconv.Itoa(int(res.statusCode)) + " " + string(res.StatusText) + "\r\n"
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
