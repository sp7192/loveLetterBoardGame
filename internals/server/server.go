package server

import (
	"fmt"
	"net"
)

type Server struct {
	ip          string
	port        int
	listener    net.Listener
	connections SafeConnections
}

func NewServer(ip string, port int) Server {
	return Server{ip: ip, port: port, connections: NewSafeConnections()}
}

func (s *Server) listen() (func() error, error) {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		return nil, fmt.Errorf("error listening: %w", err)
	}
	return s.listener.Close, nil
}

func (s *Server) shutdown() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
