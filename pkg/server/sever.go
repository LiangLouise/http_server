package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/liangLouise/http_server/pkg/router"
)

type HTTP_PROTOCOL_VERSION string

const (
	HTTP_1   string = "HTTP/1.0"
	HTTP_1_1 string = "HTTP/1.1"
	HTTP_2   string = "HTTP/2"
)

type Server interface {
	ListenRequest()
}

type server struct {
	Address  string
	Port     string
	Protocol string
	Listener net.Listener
	lock     sync.Mutex
}

func MakeServer(Adr, Port string, Protocol string) (s *server, err error) {
	port := Adr + ":" + Port
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	s = &server{
		Address:  Adr,
		Port:     Port,
		Protocol: Protocol,
		Listener: l,
	}

	return s, nil
}

func (s *server) ListenRequest() {
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go router.SimpleHandler(c)
	}
}
