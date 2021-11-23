// Read Raw http request message, convert it to object

import (
	"github.com/liangLouise/http_server/pkg/server"
)

type request struct {
	method string // GET, POST, etc.
	header textproto.MIMEHeader
	body   []byte
	uri    string                // The raw URI from the request
	proto  HTTP_PROTOCOL_VERSION // "HTTP/1.1"
}
