// API to create a new Server instance, run it and API to shut it down.
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

// This struct contains the info and variables
// server needs to serve the content
type server struct {
	Address  string
	Port     string
	Protocol httpProto.HTTP_PROTOCOL_VERSION
	// TCP Sever to Listen new connections
	Listener net.Listener
	// Waitgroup to ensure all handler go routines to be shut
	// down before shuting down server itself
	Wg sync.WaitGroup
	// Dummy channel to send shutdown signal to handler
	quit   chan interface{}
	fs     *fsService.FsService
	config *config.ServerConfig
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

// Start the new server
//
// Create an infinite loop and try to listen New TCP connection as many as it can
// Until reach MAX_CONCURRENT_CONNECTIONS
// The listening loop will be broken by calling
// 	server.ShutDown()
// to signal a shutdown
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

// Signal to shut down the server.
//
// This will let sever stop listening new TCP connection
// and try to close all existing connections.
func (s *server) ShutDown() {
	close(s.quit)
	// Decrease the main server thread waiting
	s.Listener.Close()
	// Wait for the running handler to be done
	// As we have Timeout for each handler, so it
	// should not take long.
	s.Wg.Wait()
}
