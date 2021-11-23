package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/liangLouise/http_server/pkg/router"
)

type HTTP_PROTOCOL_VERSION int

const (
	HTTP_1 HTTP_PROTOCOL_VERSION = iota
	HTTP_1_1
	HTTP_2
)


type Server interface {
	ListenRequest()
}

type server struct {
	Address  string
	Port     string
	Protocol HTTP_PROTOCOL_VERSION
	Listener net.Listener
	lock     sync.Mutex
}

func MakeServer(Adr, Port string, Protocol HTTP_PROTOCOL_VERSION) (s *server, err error) {
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

