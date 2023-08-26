package client

import (
	"fmt"
	"log"
	"loveLetterClient/internals/configs"
	"loveLetterClient/internals/logic"
	"net"
	"sync"
	"time"
)

type Client struct {
	config    *configs.Configs
	conn      net.Conn
	gameLogic *logic.GameLogic
	logger    *log.Logger
}

func NewClient(c *configs.Configs, l *log.Logger) Client {
	return Client{
		config:    c,
		gameLogic: logic.NewGameLogic(l),
		logger:    l,
	}
}

func (c *Client) connectToServer() bool {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.config.ServerIP, c.config.ServerPort))
	if err != nil {
		// TODO : reconnect
		c.logger.Printf("Error is : %s\n", err.Error())
		return false
	}
	return true
}

func (c *Client) tryConnection() bool {
	t := time.Now()
	isConnected := c.connectToServer()
	for !isConnected && time.Now().Before(t.Add(3*time.Second)) { // to be handled in config
		time.Sleep(100 * time.Millisecond)
		c.logger.Printf("Retrying to connect to server\n")
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
					c.logger.Println("Read timeout")
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

func (c *Client) sendToServerLoop(done <-chan struct{}) {
	go func() {
		for {
			select {
			case msg := <-c.gameLogic.SendMessageQueue:
				n, err := c.conn.Write([]byte(msg))
				c.logger.Printf("Wrote %s to server\n", msg[:n])
				if err != nil {
					// TODO: handle retry
					return
				}
			case <-done:
			}
		}
	}()
}

func (c *Client) Run() {
	fmt.Println("Client Started")
	if !c.tryConnection() {
		c.logger.Printf("COULD NOT CONNECT TO SERVER\n")
		return
	}
	defer c.conn.Close()

	done := make(chan struct{})
	defer close(done)
	c.sendToServerLoop(done)

	wg := sync.WaitGroup{}
	wg.Add(1)

	msgCh := c.receiveMessage(done, &wg)
	c.logger.Printf("")

	for {
		select {
		case <-done:
			break
		case msg := <-msgCh:
			c.gameLogic.ParseMessage(msg)
		}
	}
}
