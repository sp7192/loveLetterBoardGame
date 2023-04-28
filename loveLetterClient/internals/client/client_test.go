package client

import (
	"fmt"
	"loveLetterClient/internals/configs"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Create a new Configs object for testing purposes
	testConfig := &configs.Configs{
		ServerIP:   "localhost",
		ServerPort: 1234,
	}

	client := NewClient(testConfig)

	assert.Equal(t, client.config, testConfig, "NewClient() failed: config field does not match")
	assert.Nil(t, client.conn, "NewClient() failed: conn field should be nil")
}

func TestConnectToServer(t *testing.T) {
	testConfig := &configs.Configs{
		ServerIP:   "localhost",
		ServerPort: 1234,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", testConfig.ServerIP, testConfig.ServerPort))
	assert.NoError(t, err, fmt.Sprintf("Error creating test listener: %s", err))
	defer listener.Close()

	listenerReady := make(chan struct{})

	go func() {
		listenerReady <- struct{}{}
		conn, err := listener.Accept()
		assert.NoError(t, err, fmt.Sprintf("Error accepting test connection: %s", err))
		conn.Close()
	}()

	<-listenerReady

	client := &Client{
		config: testConfig,
	}

	success := client.connectToServer()
	assert.True(t, success, "connectToServer() failed: expected success but got false")
	assert.NotNil(t, client.conn, "connectToServer() failed: expected conn field to be non-nil")
}

func TestClient_tryConnection(t *testing.T) {
	c := &Client{config: &configs.Configs{
		ServerIP:   "127.0.0.1",
		ServerPort: 1235,
	}}

	// start listener
	listener, err := net.Listen("tcp", "127.0.0.1:1235")
	assert.NoError(t, err)
	defer listener.Close()

	// use a waitgroup to ensure that tryConnection() is called only after the listener starts accepting connections
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		// signal that the listener is ready
		wg.Done()

		conn, err := listener.Accept()
		assert.NoError(t, err)
		defer conn.Close()
	}()

	// wait for the listener to be ready
	wg.Wait()

	// test that tryConnection() succeeds
	isConnected := c.tryConnection()
	assert.True(t, isConnected)
}

func TestClient_receiveMessage(t *testing.T) {
	c := NewClient(&configs.Configs{
		ServerIP:   "127.0.0.1",
		ServerPort: 2334,
	})

	fmt.Println("Client Started")
	if !c.tryConnection() {
		fmt.Printf("COULD NOT CONNECT TO SERVER\n")
		return
	}
	defer c.conn.Close()

	wg := sync.WaitGroup{}
	done := make(chan struct{})
	defer close(done)

	messageChan := c.receiveMessage(done, &wg)

	// Test that receiveMessage returns a channel
	assert.NotNil(t, messageChan)

	// Test that receiveMessage receives messages
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", "localhost:2334")
		assert.NoError(t, err)
		defer conn.Close()
		conn.Write([]byte("Hello"))
	}()

	select {
	case message := <-messageChan:
		assert.Equal(t, "Hello", message)
	case <-time.After(1 * time.Second):
		t.Error("Timed out while waiting for message")
	}

	// Test that receiveMessage handles timeouts
	c.conn = nil // Disconnect the client to trigger timeout error
	select {
	case message := <-messageChan:
		t.Errorf("Expected no message but got %s", message)
	case <-time.After(2 * time.Second):
	}
}
