package server

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/liangLouise/http_server/pkg/config"
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
	config   *config.ServerConfig
}

func MakeServer(config *config.ServerConfig, fs *fsService.FsService) (s *server, err error) {
	port := config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	s = &server{
		Address:  config.Server.Host,
		Port:     strconv.Itoa(config.Server.Port),
		Protocol: config.Server.Version,
		Listener: l,
		quit:     make(chan interface{}),
		fs:       fs,
		config:   config,
	}

	// s.wg.Add(1)
	return s, nil
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

		// New Connection, now increase wait group by 1
		s.wg.Add(1)
		go func() {
			router.SimpleHandler(s.quit, c, s.fs, s.config)
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
