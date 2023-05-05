package client

import (
	"fmt"
	"loveLetterClient/internals/configs"
	"net"
	"sync"
	"time"
)

type Client struct {
	config *configs.Configs
	conn   net.Conn
}

func NewClient(c *configs.Configs) Client {
	return Client{config: c}
}

func (c *Client) connectToServer() bool {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.config.ServerIP, c.config.ServerPort))
	if err != nil {
		// TODO : reconnect
		fmt.Printf("Error is : %s\n", err.Error())
		return false
	}
	return true
}

func (c *Client) tryConnection() bool {
	t := time.Now()
	isConnected := c.connectToServer()
	for !isConnected && time.Now().Before(t.Add(3*time.Second)) { // to be handled in config
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Retrying to connect to server\n")
		isConnected = c.connectToServer()
	}
	return isConnected
}

func (c *Client) receiveMessage(done <-chan struct{}, wg *sync.WaitGroup) <-chan string {
	ret := make(chan string)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096)

		for {
			select {
			case <-done:
				return
			default:
			}

			c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			l, err := c.conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					fmt.Println("Read timeout")
					continue
				}
				fmt.Println("Error reading:", err.Error())
				return
			}

			ret <- string(buffer[:l])
		}
	}()
	return ret
}

func (c *Client) Run() {
	fmt.Println("Client Started")
	if !c.tryConnection() {
		fmt.Printf("COULD NOT CONNECT TO SERVER\n")
		return
	}
	defer c.conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	done := make(chan struct{})
	defer close(done)
	msgCh := c.receiveMessage(done, &wg)
	fmt.Printf("")

	for {
		select {
		case <-done:
			break
		case msg := <-msgCh:
			fmt.Printf("Server >> %s\n", msg)
		}
	}
	wg.Wait()
}
