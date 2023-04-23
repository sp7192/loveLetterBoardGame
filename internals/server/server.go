package server

import (
	"fmt"
	"loveLetterBoardGame/internals/configs"
	"net"
	"sort"
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
		return func() error { return nil }, fmt.Errorf("error listening: %w", err)
	}
	return func() error {
		if s == nil || s.listener == nil {
			return fmt.Errorf("listener is nil")
		}
		return s.listener.Close()
	}, nil
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

func (s *Server) GetClientsIds() []uint {
	ret := make([]uint, 0, s.config.PlayersInRoomCount)
	conns := s.connections.GetAllConnections()
	for k := range conns {
		ret = append(ret, k)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

func (s *Server) shutdown() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
