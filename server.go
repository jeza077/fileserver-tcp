package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	channels map[string]*channel
	commands chan command
}

func newServer() *server {
	return &server{
		channels: make(map[string]*channel),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_USER:
			s.user(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_CHANNELS:
			s.listChannels(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_FILE:
			s.file(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("Nuevo cliente conectado: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		user:     "anonimo",
		commands: s.commands,
	}

	// Leer mensajes
	c.readInput()

}

func (s *server) user(c *client, args []string) {
	c.user = args[1]
	c.msg(fmt.Sprintf("Hola, %s", c.user))
}

func (s *server) join(c *client, args []string) {
	// Obtener nombre del canal que ingresa el cliente
	channelName := args[1]

	// Validamos si existe el canal con ese nombre sino creamos uno nuevo
	ch, ok := s.channels[channelName]
	if !ok {
		ch = &channel{
			name:    channelName,
			members: make(map[net.Addr]*client),
		}
		s.channels[channelName] = ch
	}

	// Agregamos el cliente actual al canal
	ch.members[c.conn.RemoteAddr()] = c

	// Salimos del canal actual
	s.quitCurrentChannel(c)
	c.channel = ch

	ch.broadcast(c, fmt.Sprintf("%s se ha unido al canal", c.user))
	c.msg(fmt.Sprintf("Bienvenido a %s", ch.name))
}

func (s *server) listChannels(c *client, args []string) {
	var channels []string
	for name := range s.channels {
		channels = append(channels, name)
	}

	c.msg(fmt.Sprintf("Canales disponibles: %s", strings.Join(channels, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if c.channel == nil {
		c.err(errors.New("debes unirte a un canal primero"))
		return
	}

	c.channel.broadcast(c, c.user+": "+strings.Join(args[1:len(args)], " "))
}

func (s *server) file(c *client, args []string) {
	if c.channel == nil {
		c.err(errors.New("debes unirte a un canal primero"))
		return
	}

	name, size, bytes, err := LoadFile(args[1])
	if err != nil {
		fmt.Println(err)
		// continue
	}

	f := &File{
		Name:    name,
		Size:    size,
		Content: bytes,
	}

	c.channel.broadcast(c, fmt.Sprintf("%s envio un archivo: %s, %d bytes: %d", c.user, f.Name, f.Size, f.Content))

}

func (s *server) quit(c *client, args []string) {
	log.Printf("Cliente se ha desconectado: %s", c.conn.RemoteAddr().String())

	s.quitCurrentChannel(c)

	c.msg(("¡Que vuelvas pronto!"))
	c.conn.Close()
}

func (s *server) quitCurrentChannel(c *client) {
	if c.channel != nil {
		delete(c.channel.members, c.conn.RemoteAddr())
		c.channel.broadcast(c, fmt.Sprintf("%s ha abandonado el canal", c.user))
	}
}
