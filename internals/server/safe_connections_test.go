package server

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSafeConnections(t *testing.T) {
	sc := NewSafeConnections()
	assert.NotNil(t, sc.connections)
}

func TestSafeConnections_Read(t *testing.T) {
	sc := NewSafeConnections()
	id := uint(1)
	expected := &net.TCPConn{}
	sc.Set(id, expected)

	// Test reading a valid connection
	got, err := sc.Get(id)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, got, expected)

	// Test reading an invalid connection
	invalidConn, err := sc.Get(id + 1)
	assert.Nil(t, invalidConn)
}

func TestSafeConnections_Write(t *testing.T) {
	sc := NewSafeConnections()
	id := uint(1)
	conn := &net.TCPConn{}

	// Test writing a new connection
	sc.Set(id, conn)
	assert.Equal(t, conn, sc.connections[id])

	// Test overwriting an existing connection
	newConn := &net.UDPConn{}
	sc.Set(id, newConn)
	assert.Equal(t, newConn, sc.connections[id])
}

func TestSafeConnections_GetAllConnections(t *testing.T) {
	sc := NewSafeConnections()
	conn1 := &net.TCPConn{}
	conn2 := &net.UDPConn{}
	sc.Set(1, conn1)
	sc.Set(2, conn2)

	// Test getting all connections
	allConns := sc.GetAllConnections()
	assert.Len(t, allConns, 2)
	assert.Equal(t, allConns[1], conn1)
	assert.Equal(t, allConns[2], conn2)

	// Test getting all connections after modifying the original map
	delete(sc.connections, 1)
	modifiedConns := sc.GetAllConnections()
	assert.Len(t, modifiedConns, 1)
	assert.Equal(t, modifiedConns[2], conn2)
}
