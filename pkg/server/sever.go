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
	"golang.org/x/net/netutil"
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
	Wg       sync.WaitGroup
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

	// Set the upper bound of the TCP connections the
	// server can accpet
	if config.RunTime.MaxConnections <= 0 {
		log.Fatal("The max TCP connection number is expected to be greater than 0")
	}
	l = netutil.LimitListener(l, config.RunTime.MaxConnections)

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
			// Listner close and got shutdown signal
			case <-s.quit:
				// break the listener loop
				return
			default:
				log.Println("accept error", err)
				continue
			}
		}

		// New Connection, now increase wait group by 1
		s.Wg.Add(1)

		go func() {
			router.SimpleHandler(s.quit, c, s.fs, s.config)
			s.Wg.Done()
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
	s.Wg.Wait()
}
