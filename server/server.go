package server

import (
	"fmt"
	"net"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/Amr-Nashaatx/http-from-tcp/request"
)

type Server struct {
	socket      net.Listener
	connections []net.Conn
	requests    chan *request.Request
	// atomic.Bool because Close() and listen() run in different goroutines — a plain bool has no visibility guarantee across cores.
	closed atomic.Bool
	// guards connections: listen() appends and each handle() goroutine deletes concurrently.
	mu sync.Mutex
}

func Serve(port int) (*Server, error) {
	socket, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("cannot open listen on connections %w", err)
	}
	server := &Server{socket: socket, requests: make(chan *request.Request, 10)}
	go server.listen()

	return server, nil
}
func (s *Server) Close() error {
	closeErr := s.socket.Close()
	if closeErr != nil {
		return closeErr
	}
	s.closed.Store(true)
	return nil
}
func (s *Server) listen() {
	for {
		conn, acceptErr := s.socket.Accept()
		// Accept returns an error when the socket is closed; distinguish intentional shutdown from real errors.
		if acceptErr != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println(acceptErr)
		}
		// Lock because a slice header (ptr/len/cap) is not atomically writable — concurrent append and delete would corrupt it.
		s.mu.Lock()
		s.connections = append(s.connections, conn)
		s.mu.Unlock()
		go func() {
			req, err := request.RequestFromReader(conn)
			if err != nil {
				fmt.Println(err)
				return
			}
			s.handle(conn)
			s.requests <- req
		}()
	}
}
func (s *Server) handle(conn net.Conn) {
	response := "HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"

	conn.Write([]byte(response))
	conn.Close()

	// Same lock as in listen() — multiple handle() goroutines can finish at the same time and race on this slice.
	s.mu.Lock()
	s.connections = slices.DeleteFunc(s.connections, func(openConn net.Conn) bool {
		return openConn == conn
	})
	s.mu.Unlock()
}
