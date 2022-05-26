package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	user     string
	channel  *channel
	commands chan<- command
}

func (c *client) readInput() {
	for {
		// msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0]) //Eliminar espacios

		// msg = strings.TrimSpace(string(msg))
		// args := strings.Split(msg, " ")
		// cmd := strings.TrimSpace(args[0]) //Eliminar espacios

		switch cmd {
		case "/user":
			c.commands <- command{
				id:     CMD_USER,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/channels":
			c.commands <- command{
				id:     CMD_CHANNELS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("Comando desconocido: %s", cmd))
		}

	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
