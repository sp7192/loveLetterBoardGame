package server

import (
	"fmt"
	"loveLetterBoardGame/internals/configs"
	"net"
)

type Server struct {
	ip          string
	port        int
	listener    net.Listener
	connections SafeConnections
	config      configs.Configs
}

func NewServer(ip string, port int, conf configs.Configs) Server {
	return Server{ip: ip, port: port, connections: NewSafeConnections(), config: conf}
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
