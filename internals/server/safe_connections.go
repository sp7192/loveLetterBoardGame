package server

import (
	"fmt"
	"net"
	"sync"
)

type SafeConnections struct {
	connections map[uint]net.Conn
	lock        sync.RWMutex
}

func NewSafeConnections() SafeConnections {
	return SafeConnections{connections: make(map[uint]net.Conn), lock: sync.RWMutex{}}
}

func (s *SafeConnections) Read(id uint) (net.Conn, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	conn, ok := s.connections[id]
	if !ok {
		return nil, fmt.Errorf("connection not found for id %d", id)
	}
	return conn, nil
}

func (s *SafeConnections) Write(id uint, conn net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.connections[id] = conn
}

func (s *SafeConnections) Count() uint {
	return uint(len(s.connections))
}

func (s *SafeConnections) GetAllConnections() map[uint]net.Conn {
	s.lock.RLock()
	defer s.lock.RUnlock()
	copy := make(map[uint]net.Conn, len(s.connections))
	for k, v := range s.connections {
		copy[k] = v
	}
	return copy
}
