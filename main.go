package main

import (
	"fmt"
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

	fmt.Println("1. '/nickName' for to define an alias")
	fmt.Println("2. '/rooms' for to list rooms chat")
	fmt.Println("3. '/join' for to join a room chat")
	fmt.Println("4. '/msg' for to send a message to a room chat")
	fmt.Println("5. '/quit' for to quit the chat")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Error accepting connection: %s", err.Error())
			continue
		}

		log.Printf("New client connected: %s", conn.RemoteAddr().String())

		go server.newClient(conn)
	}

}
