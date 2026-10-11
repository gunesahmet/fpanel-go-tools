package Cloudbackupmng

import (
	"CrypterWin"
	"encoding/json"
	"net"
	"strings"
)

type Connection struct {
	ServerId      string
	ServerAddress string
	ServerPort    string
	ServerKey     string
	ServerIv      string
	ServerType    ServerType
	Request       struct {
		Job        string   `json:"job"`
		Command    string   `json:"command"`
		Parameters []string `json:"parameters"`
	}
	Error   error
	Crypter *CrypterWin.Crypter
}

type Request struct {
	Job        string   `json:"job"`
	Command    string   `json:"command"`
	Parameters []string `json:"parameters"`
}

type ServerType int64

const (
	MAIN  ServerType = 0
	AGENT ServerType = 1
)

type GetName int

func (getName GetName) String() string {
	switch getName {
	case 0:
		return "MAIN"
	case 1:
		return "AGENT"
	default:
		return "UNDEFINED"
	}
}

func NewConnect() *Connection {
	return &Connection{}
}

func (c *Connection) SetServerId(serverId string, serverType ServerType) (*Connection, bool) {
	// Get server datas from mysql
	c.ServerId = serverId
	c.ServerType = serverType
	if serverType == MAIN {
		c.ServerPort = "8889"
	} else if serverType == AGENT {
		c.ServerPort = "8887"
	} else {
		Env.Logger.Error(1616, "Wrong server type. Error: %v", c.Error)
		return nil, false
	}

	c.Crypter, c.Error = CrypterWin.NewCrypter([]byte(c.ServerKey), []byte(c.ServerIv))
	if c.Error != nil {
		Env.Logger.Error(1617, "Server set crypter failed. Error: %v", c.Error)
		return nil, false
	}
	return c, true
}

func (c *Connection) SetRequest(request Request) *Connection {
	c.Request = request
	return c
}

func (c *Connection) Connect() bool {
	addr := strings.Join([]string{c.ServerAddress, c.ServerPort}, ":")
	var conn net.Conn
	conn, c.Error = net.Dial("tcp", addr)
	if c.Error != nil {
		Env.Logger.Error(1618, "Connection failed. Error: %v", c.Error)
		return false
	} else {
		var command []byte
		command, c.Error = json.Marshal(c.Request)
		encoded := c.Crypter.EncodeToString(string(command))
		if c.Error != nil {
			Env.Logger.Error(1619, "Connection drop. Error: %v", c.Error)
			return false
		}

		_, c.Error = conn.Write([]byte(encoded))
		if c.Error != nil {
			Env.Logger.Error(1620, "conn.Write. Error: %v", c.Error)
			return false
		}
	}
	return true
}

func (c *Connection) Response() []byte {

	return nil
}
