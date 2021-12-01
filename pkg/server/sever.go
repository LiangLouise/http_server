package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liangLouise/http_server/pkg/httpProto"
	"github.com/liangLouise/http_server/pkg/router"
)

type Server interface {
	ListenRequest()
}

type server struct {
	Address  string
	Port     string
	Protocol httpProto.HTTP_PROTOCOL_VERSION
	Listener net.Listener
	lock     sync.Mutex
}

func MakeServer(Adr, Port string, Protocol httpProto.HTTP_PROTOCOL_VERSION) (s *server, err error) {
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

func SetTimeout(conn net.Conn, timer int) error {
	timeout := time.Duration(timer) * (time.Second)
	err := conn.SetDeadline(time.Now().Add(timeout))
	return err
}

func (s *server) ListenRequest() {
	for {

		c, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// err = SetTimeout(c, 5)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }

		go router.SimpleHandler(c)
	}
}
