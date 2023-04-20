package server

import (
	"fmt"
	"loveLetterBoardGame/internals/configs"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestStart(t *testing.T) {
	// Create a new Configs instance
	conf := configs.Configs{}
	conf.PlayersInRoomCount = 2

	// Create a new Server instance
	srv := NewServer("localhost", 1234, conf)

	// Use a channel to signal when acceptClients has finished accepting connections
	done := make(chan struct{})
	go func() {
		err := srv.Start()
		assert.NoError(t, err)
		close(done)
	}()

	// Connect two clients to the server
	conn1, err := net.Dial("tcp", "localhost:1234")
	require.NoError(t, err, "Failed to connect first client")
	defer conn1.Close()

	conn2, err := net.Dial("tcp", "localhost:1234")
	require.NoError(t, err, "Failed to connect second client")
	defer conn2.Close()

	// Wait for acceptClients to finish accepting connections
	select {
	case <-done:
		// acceptClients has finished accepting connections
	case <-time.After(time.Second):
		// Wait for at most 1 second
		t.Fatalf("acceptClients did not finish accepting connections")
	}

	// Check that the connections were added to the SafeConnections instance
	assert.Equal(t, 2, int(srv.connections.Count()), "Expected two connections to be added")

	// Check that the listener was closed
	assert.ErrorIs(t, srv.shutdown(), net.ErrClosed, "Expected the listener to be closed")
}

func TestAcceptClients(t *testing.T) {
	// Create a new Configs instance
	conf := configs.Configs{}
	conf.PlayersInRoomCount = 2

	// Create a new Server instance
	srv := NewServer("localhost", 1234, conf)
	closer, err := srv.listen()
	assert.NoError(t, err)
	defer closer()

	// Use a channel to signal when acceptClients has finished accepting connections
	done := make(chan struct{})
	go func() {
		err := srv.acceptClients()
		assert.NoError(t, err)
		close(done)
	}()

	// Connect two clients to the server
	conn1, err := net.Dial("tcp", "localhost:1234")
	require.NoError(t, err, "Failed to connect first client")
	defer conn1.Close()

	conn2, err := net.Dial("tcp", "localhost:1234")
	require.NoError(t, err, "Failed to connect second client")
	defer conn2.Close()

	// Wait for acceptClients to finish accepting connections
	select {
	case <-done:
		// acceptClients has finished accepting connections
	case <-time.After(time.Second):
		// Wait for at most 1 second
		t.Fatalf("acceptClients did not finish accepting connections")
	}

	// Check that the connections were added to the SafeConnections instance
	assert.Equal(t, 2, int(srv.connections.Count()), "Expected two connections to be added")
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
		s.connections.Set(id, conn)
	}

	// Test getting all connections
	got := s.connections.GetAllConnections()
	assert.Equal(t, expected, got)
}
