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

func NewServer(conf configs.Configs) Server {
	return Server{ip: conf.ServerIP, port: int(conf.ServerPort), connections: NewSafeConnections(), config: conf}
}

func (s *Server) listen() (func() error, error) {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		return nil, fmt.Errorf("error listening: %w", err)
	}
	return s.listener.Close, nil
}

func (s *Server) acceptClients() error {
	for {
		if s.connections.Count() >= s.config.PlayersInRoomCount {
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting client: %w", err)
		}
		s.connections.Set(s.connections.Count()+1, conn)
	}
	return nil
}

func (s *Server) Start() error {
	closer, err := s.listen()
	if err != nil {
		return err
	}
	defer closer()

	err = s.acceptClients()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) shutdown() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
