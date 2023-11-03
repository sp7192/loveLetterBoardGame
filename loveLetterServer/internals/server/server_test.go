package server

import (
	"fmt"
	"log"
	"loveLetterBoardGame/internals/configs"
	"loveLetterBoardGame/models"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	conf := configs.Configs{PlayersInRoomCount: 4, ServerIP: "127.0.0.1", ServerPort: 8200} // create a mock config object
	s := NewServer(conf, log.Default())
	assert.Equal(t, conf.ServerIP, s.ip)
	assert.Equal(t, conf.ServerPort, uint(s.port))
}

func TestServer_listen(t *testing.T) {
	conf := configs.Configs{PlayersInRoomCount: 4, ServerIP: "127.0.0.1", ServerPort: 8000} // create a mock config object
	s := NewServer(conf, log.Default())

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
	conf := configs.Configs{PlayersInRoomCount: 4, ServerIP: "127.0.0.1", ServerPort: 9200} // create a mock config object
	s := NewServer(conf, log.Default())
	closer, err := s.listen()

	// Test listen error
	closer, err = s.listen()
	assert.Error(t, err)
	closer()
}

func TestGetClientsIds(t *testing.T) {

	// Create a mock server with some connections
	server := &Server{
		config:      configs.Configs{PlayersInRoomCount: 2},
		connections: NewSafeConnections(),
	}

	server.connections.Set(1, &net.TCPConn{})
	server.connections.Set(2, &net.TCPConn{})

	// Call the method being tested
	result := server.GetClientsIds()

	// Check that the result is as expected
	expected := []uint{1, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestStart(t *testing.T) {
	// Create a new Configs instance
	conf := configs.Configs{PlayersInRoomCount: 2, ServerIP: "127.0.0.1", ServerPort: 8080} // create a mock config object

	// Create a new Server instance
	srv := NewServer(conf, log.Default())

	// Use a channel to signal when acceptClients has finished accepting connections
	done := make(chan struct{})
	go func() {
		err := srv.Start()
		assert.NoError(t, err)
		close(done)
	}()

	// Connect two clients to the server
	conn1, err := net.Dial("tcp", "localhost:8080")
	require.NoError(t, err, "Failed to connect first client")
	defer conn1.Close()

	conn2, err := net.Dial("tcp", "localhost:8080")
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

func TestSendMessagesToClients(t *testing.T) {
	server := &Server{
		sendMessageChannel: make(chan models.ServerMessage),
		connections:        NewSafeConnections(),
	}

	clientConn, serverConn := net.Pipe()
	var clientID uint = 1
	server.connections.Set(clientID, serverConn)
	server.sendMessagesToClients()

	ch := make(chan func(), 2)

	// Add the read operation to the channel.
	ch <- func() {
		buf := make([]byte, 1024)
		n, err := clientConn.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, buf[:n], "Hello, client!")
	}

	// Add the write operation to the channel.
	ch <- func() {
		msg := models.ServerMessage{ToClientId: clientID, Message: "Hello, client!"}
		server.sendMessageChannel <- msg
	}

	// Run the functions on separate goroutines.
	for i := 0; i < 2; i++ {
		go func() {
			fn := <-ch
			fn()
		}()
	}
}

func TestAcceptClients(t *testing.T) {
	// Create a new Configs instance
	conf := configs.Configs{
		PlayersInRoomCount: 2,
		ServerIP:           "localhost",
		ServerPort:         1234,
	}

	// Create a new Server instance
	srv := NewServer(conf, log.Default())
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
	conf := configs.Configs{PlayersInRoomCount: 4, ServerIP: "127.0.0.1", ServerPort: 9000} // create a mock config object
	s := NewServer(conf, log.Default())

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
