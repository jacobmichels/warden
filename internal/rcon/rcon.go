package rcon

import (
	"fmt"
	"log"

	"github.com/gorcon/rcon"
)

type Client struct {
	conn *rcon.Conn
}

func NewClient(address, password string) (*Client, error) {
	conn, err := rcon.Dial(address, password)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rcon: %w", err)
	}

	log.Println("RCON connected")

	return &Client{
		conn,
	}, nil
}

// Sends a command to the rcon server and returns the response.
func (c *Client) SendCommand(command string) (string, error) {
	return c.conn.Execute(command)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
