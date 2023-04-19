package server

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	s := NewServer(ip, port)
	assert.Equal(t, ip, s.ip)
	assert.Equal(t, port, s.port)
	assert.NotNil(t, s.connections)
}

func TestServer_listen(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	s := NewServer(ip, port)

	// Test successful listen
	done := make(chan bool)
	go func() {
		close(done)
	}()
	defer func() {
		<-done
	}()

	closeFn := s.listen()
	assert.NotNil(t, s.listener)

	// Test listener closing
	closeFn()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestServer_listen_error(t *testing.T) {
	ip := "invalid_ip"
	port := 8080
	s := NewServer(ip, port)

	// Test listen error
	assert.Panics(t, func() {
		s.listen()
	})
}

func TestServer_GetAllConnections(t *testing.T) {
	ip := "127.0.0.1"
	port := 8080
	s := NewServer(ip, port)

	expected := map[uint]net.Conn{
		1: &net.TCPConn{},
		2: &net.TCPConn{},
		3: &net.TCPConn{},
	}
	for id, conn := range expected {
		s.connections.Write(id, conn)
	}

	// Test getting all connections
	got := s.GetAllConnections()
	assert.Equal(t, expected, got)
}
