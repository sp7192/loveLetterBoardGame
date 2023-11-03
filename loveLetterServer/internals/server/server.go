package server

import (
	"encoding/json"
	"fmt"
	"log"
	"loveLetterBoardGame/internals/configs"
	"loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/models"
	"net"
	"sort"
	"time"
)

type Server struct {
	ip               string
	port             int
	listener         net.Listener
	connections      SafeConnections
	config           configs.Configs
	receivedMessages chan models.ClientMessage
	logger           *log.Logger
}

func NewServer(conf configs.Configs, l *log.Logger) Server {
	return Server{ip: conf.ServerIP,
		port:             int(conf.ServerPort),
		connections:      NewSafeConnections(),
		config:           conf,
		receivedMessages: make(chan models.ClientMessage),
		logger:           l,
	}
}

func (s *Server) listen() (func() error, error) {
	var err error
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)
	s.logger.Printf("listening on %s\n", addr)
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

func (s *Server) handleClientMessage(id uint) {
	go func() {
		buffer := make([]byte, 4096)
		conn, err := s.connections.Get(id)
		if err != nil {
			s.logger.Println("Error connections ", err.Error())
			return
		}
		for {
			l, err := conn.Read(buffer)
			if err != nil {
				s.logger.Println("Error reading:", err.Error())
				return
			}

			if err != nil {
				s.logger.Println("Error reading:", err.Error())
				return
			}
			s.receivedMessages <- models.ClientMessage{ClientId: id, Message: string(buffer[:l])}
		}
	}()
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
		s.handleClientMessage(id)
		s.SendAndReceiveAck(id, models.InfoMessage, fmt.Sprintf("Your id set by server is : %d", id))
	}
	return nil
}

func (s *Server) SendToClient(msg models.ServerMessage) error {
	conn, err := s.connections.Get(msg.ToClientId)
	if err != nil {
		s.logger.Printf("Errror in reading connection :%s\n", err.Error())
		return err
	}
	if err := writeToConnection(conn, []byte(msg.Message)); err != nil {
		s.logger.Printf("Error in sending message: %s\n", err.Error())
		return err
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

func (s *Server) SendTo(id uint, msgType models.MessageType, msgPayload string) error {
	s.logger.Printf("In SendTo, id=%d: Type : %s, Payload: %s\n", id, msgType, msgPayload)
	data := models.Message{
		Type:    msgType,
		Payload: msgPayload,
	}
	msg, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}
	err = s.SendToClient(models.ServerMessage{ToClientId: id, Message: string(msg)})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) SendAndReceiveAck(id uint, msgType models.MessageType, msgPayload string) error {
	err := s.SendTo(id, msgType, msgPayload)
	if err != nil {
		return err
	}
	clientMessage, err := s.GetClientMessage()
	if err != nil {
		return err
	}
	var msg models.Message
	err = json.Unmarshal([]byte(clientMessage.Message), &msg)
	if err != nil {
		return err
	}

	if clientMessage.ClientId != id {
		return fmt.Errorf("ack client id is different, expected : %d, got : %d\n", id, clientMessage.ClientId)
	}

	if msg.Type != models.AckMessage || msg.Payload != string(msgType) {
		return fmt.Errorf("expected : %s, got : %s\n", models.AckMessage, msg.Type)
	}

	s.logger.Printf("Received ack of %s\n", msg.Payload)
	return nil
}

func (s *Server) SendToAllWithAck(state gamelogic.GameState) error {
	data, err := json.MarshalIndent(state, "", "	")
	if err != nil {
		return err
	}

	ids := s.GetClientsIds()
	for _, id := range ids {
		s.SendAndReceiveAck(id, models.UpdateMessage, string(data))
	}
	return nil
}

func (s *Server) GetClientMessage() (models.ClientMessage, error) {
	select {
	case ret := <-s.receivedMessages:
		return ret, nil
	case <-time.After(120 * time.Second): // TODO : Change magic number to read from config
		return models.ClientMessage{}, fmt.Errorf("time out")
	}
}
