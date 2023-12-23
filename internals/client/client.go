package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"loveLetterBoardGame/internals/configs"
	"loveLetterBoardGame/internals/logic"
	"loveLetterBoardGame/models"
	"net"
	"os"
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
		gameLogic: logic.NewGameLogic(l, bufio.NewScanner(os.Stdin)),
		logger:    l,
	}
}

func (c *Client) connectToServer() bool {
	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.config.ServerIP, c.config.ServerPort))
	if err != nil {
		c.logger.Fatalf("cannot connect to server : %s\n", err.Error())
		return false
	}
	c.gameLogic.SendMessage = c.conn.Write
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
			fmt.Printf("BUUUUUUUUFER : %s\n", string(buffer[:l]))
			ret <- string(buffer[:l])
		}
	}()
	return ret
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	msgCh := c.receiveMessage(done, &wg)

	c.logger.Printf("")

	for {
		select {
		case <-done:
			break
		case strMsg := <-msgCh:
			var msg models.Message
			err := json.Unmarshal([]byte(strMsg), &msg)
			if err != nil {
				// TODO: Send Not ok Message
				c.logger.Printf("ERROR PARSING MESSAGE %s: error is %s\n", strMsg, err.Error())
				continue
			}
			if c.gameLogic.DoesSendAck(msg) {
				c.gameLogic.SendReceivedAck(msg)
			}
			c.gameLogic.Update(msg)
		}
	}
}
