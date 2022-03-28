package main

import (
	"log"
	"net"
)

func main() {
	server := newServer()
	go server.run()

	listener, err := net.Listen("tcp", ":8888")

	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Listening on port 8888")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Error accepting connection: %s", err.Error())
			continue
		}

		go server.newClient(conn)
	}

}
