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

	// test connection retry timeout
	start := time.Now()
	_, err = listener.Accept()
	assert.NoError(t, err)
	end := time.Now()
	elapsed := end.Sub(start)
	assert.True(t, elapsed > 3*time.Second)
}
