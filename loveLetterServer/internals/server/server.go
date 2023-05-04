package server

import (
	"encoding/json"
	"fmt"
	"loveLetterBoardGame/internals/configs"
	"loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/models"
	"net"
	"sort"
)

type Server struct {
	ip             string
	port           int
	listener       net.Listener
	connections    SafeConnections
	config         configs.Configs
	messageChannel chan models.ServerMessage
}

func NewServer(conf configs.Configs) Server {
	return Server{ip: conf.ServerIP,
		port:           int(conf.ServerPort),
		connections:    NewSafeConnections(),
		config:         conf,
		messageChannel: make(chan models.ServerMessage),
	}
}

func (s *Server) listen() (func() error, error) {
	var err error
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)
	fmt.Printf("listening on %s\n", addr)
	s.listener, err = net.Listen("tcp", addr)
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
		id := s.connections.Count() + 1
		s.connections.Set(id, conn)
		s.messageChannel <- models.ServerMessage{ToClientId: id, Message: fmt.Sprintf("Your id set by server is : %d", id)}
	}
	return nil
}

func (s *Server) sendMessagesToClients() {
	go func() {
		for {
			msg := <-s.messageChannel
			conn, err := s.connections.Get(msg.ToClientId)
			if err != nil {
				fmt.Printf("Errror in reading connection :%s\n", err.Error())
				continue
			}
			if err := writeToConnection(conn, []byte(msg.Message)); err != nil {
				fmt.Printf("Error in sending message: %s\n", err.Error())
				continue
			}
		}
	}()
}

func (s *Server) Start() error {
	closer, err := s.listen()
	if err != nil {
		return err
	}
	defer closer()

	s.sendMessagesToClients()
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

func (s *Server) SendTo(id uint, msg string) (int, error) {
	conn, err := s.connections.Get(id)
	if err != nil {
		return 0, err
	}
	return conn.Write([]byte(msg))
}

func (s *Server) SendToAll(g *gamelogic.GameLogic) error {
	data, err := json.Marshal(g)
	if err != nil {
		return err
	}
	ids := s.GetClientsIds()
	for _, id := range ids {
		s.SendTo(id, string(data))
	}
	return nil
}
