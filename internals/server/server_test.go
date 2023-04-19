package server

import (
	"fmt"
	"loveLetterBoardGame/internals/configs"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	conf := configs.Configs{PlayersInRoomCount: 4} // create a mock config object
	s := NewServer(ip, port, conf)
	assert.Equal(t, ip, s.ip)
	assert.Equal(t, port, s.port)
}

func TestServer_listen(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	conf := configs.Configs{PlayersInRoomCount: 4} // create a mock config object
	s := NewServer(ip, port, conf)

	// Test successful listen
	closer, err := s.listen()
	assert.NoError(t, err)
	assert.NotNil(t, s.listener)

	// Test listener closing
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	closer()

	// Wait for the listener to be closed before checking if a connection can be established
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			break
		}
		conn.Close()
		time.Sleep(time.Millisecond * 100)
	}
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestServer_listen_error(t *testing.T) {
	ip := "invalid_ip"
	port := 8080

	conf := configs.Configs{PlayersInRoomCount: 4} // create a mock config object
	s := NewServer(ip, port, conf)

	// Test listen error
	closer, err := s.listen()
	assert.Error(t, err)
	assert.Nil(t, closer)
	assert.Contains(t, err.Error(), "lookup")
}

func TestServer_GetAllConnections(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	conf := configs.Configs{PlayersInRoomCount: 4} // create a mock config object
	s := NewServer(ip, port, conf)

	expected := map[uint]net.Conn{
		1: &net.TCPConn{},
		2: &net.TCPConn{},
		3: &net.TCPConn{},
	}
	for id, conn := range expected {
		s.connections.Write(id, conn)
	}

	// Test getting all connections
	got := s.connections.GetAllConnections()
	assert.Equal(t, expected, got)
}
