package TCPServer

import (
	Crypt "CrypterWin"
	"Helpers"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	conn    net.Conn
	host    string
	port    int
	Channel chan interface{}
	Key     string
	Iv      string
}

type Request struct {
	Command string `json:"command"`
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Servers(host string, port int) *Client {
	c.host = host
	c.port = port
	c.Channel = make(chan interface{}, 2)
	return c
}

func (c *Client) Secrets(key string, iv string) *Client {
	c.Key = key
	c.Iv = iv
	return c
}

func (c *Client) Connect() (*Client, error) {
	addr := strings.Join([]string{c.host, strconv.Itoa(c.port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		Helpers.HandleError(err, 1030, fmt.Sprintf("Connection drop. Error: %v", err))
		return c, err
	}
	c.conn = conn
	return c, nil
}

func (c *Client) Send(request Request) string {
	crypt, _ := Crypt.NewCrypter([]byte(c.Key), []byte(c.Iv))
	command := crypt.EncodeToString(request.Command)

	_, errF := c.conn.Write([]byte(command))
	if errF != nil {
		Helpers.HandleError(errF, 1031, fmt.Sprintf("conn.Write. Error: %v", errF))
		return ""
	}
	reader := bufio.NewReader(c.conn)

	msg, _ := reader.ReadString('\n') // \n this is our delimiter
	if len(Helpers.Trim(msg)) > 0 {
		fmt.Println(Helpers.Trim(msg))
		return crypt.DecodeToString(msg)
	}
	return msg
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
