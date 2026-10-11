package Ftp

import (
	"Helpers"
	"fmt"
	"github.com/jlaffaye/ftp"
	"os"
	"time"
)

type Client struct {
	c *ftp.ServerConn
}

func NewFtpClient(address string, port string) (*Client, error) {
	c, err := ftp.Dial(address+":"+port, ftp.DialWithTimeout(3*time.Second))
	if err != nil {
		Helpers.HandleError(err, 1015, fmt.Sprintf("Ftp connection failed: %v", err))
	}
	return &Client{c}, nil
}

func (c *Client) Login(username string, password string) (bool, error) {
	err := c.c.Login(username, password)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (c *Client) UploadFile(target string, source string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	if err = c.c.Stor(target, sourceFile); err != nil {
		return err
	}
	return nil
}

func (c *Client) Quit() error {
	if err := c.c.Quit(); err != nil {
		return err
	}
	return nil
}
