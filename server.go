package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (host *server) run() {
	for cmd := range host.commands {
		switch cmd.id {
		case CMD_NICK:
			host.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			host.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			host.listRooms(cmd.client, cmd.args)
		case CMD_MESSAGE:
			host.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			host.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymus",
		commands: s.commands,
	}

	c.readInput()

}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("Your nickname is now %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	roonName := args[1]
	r, ok := s.rooms[roonName]

	if !ok {
		r = &room{
			name:    roonName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roonName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("You are now in room %s", r.name))
}

func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	if len(rooms) == 0 {
		c.msg("There are no rooms")
		return
	}

	c.msg(fmt.Sprintf("Avaible rooms are: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("You are in a room, you must leave it before quitting"))
		return
	}
	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:len(args)], " "))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("Client %s disconnected", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.msg("Sad to see you go :( ")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {

	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
