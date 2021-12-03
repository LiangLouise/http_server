package server

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/liangLouise/http_server/pkg/fsService"
	"github.com/liangLouise/http_server/pkg/httpProto"
	"github.com/liangLouise/http_server/pkg/router"
)

type Server interface {
	ListenRequest()
	ShutDown()
}

type server struct {
	Address  string
	Port     string
	Protocol httpProto.HTTP_PROTOCOL_VERSION
	Listener net.Listener
	wg       sync.WaitGroup
	quit     chan interface{}
	fs       *fsService.FsService
}

func MakeServer(Adr, Port string, Protocol httpProto.HTTP_PROTOCOL_VERSION, fs *fsService.FsService) (s *server, err error) {
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
		quit:     make(chan interface{}),
		fs:       fs,
	}

	// s.wg.Add(1)
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
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
				continue
			}
		}
		// err = SetTimeout(c, 5)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }

		// New Connection, now increase wait group by 1
		s.wg.Add(1)
		go func() {
			router.SimpleHandler(s.quit, c)
			s.wg.Done()
		}()
	}
}

func (s *server) ShutDown() {
	close(s.quit)
	// Decrease the main server thread waiting
	s.Listener.Close()
	// Wait for the running handler to be done
	// As we have Timeout for each handler, so it
	// should not take long.
	s.wg.Wait()
}
