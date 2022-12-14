package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Message struct {
	conn    net.Conn
	message []byte
}

var connections []net.Conn
var messages = make(chan Message)
var addClient = make(chan net.Conn)
var removeClient = make(chan net.Conn)

func main() {
	server, err := net.Listen("tcp", ":8080")
	fmt.Printf("[+] Listening on localhost TCP port 8080. Press Ctrl+C to Terminate.\n")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	go startChannels()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		addClient <- conn
		go handleRequest(conn)
	}
}

func startChannels() {
	for {
		select {
		case message := <-messages:
			broadcastMessage(&message)
		case newClient := <-addClient:
			connections = append(connections, newClient)
			fmt.Println(len(connections))
		case deadClient := <-removeClient:
			removeConn(deadClient)
			fmt.Println(len(connections))
		}
	}
}
func handleRequest(conn net.Conn) {
	for {
		message := make([]byte, 1024)

		_, err := conn.Read(message)
		if err != nil {
			if err == io.EOF {
				removeClient <- conn
				conn.Close()
				return
			}
			log.Fatal(err)
		}
		m := Message{
			conn:    conn,
			message: message,
		}
		messages <- m
	}

}

func broadcastMessage(m *Message) {
	for _, conn := range connections {
		if m.conn != conn {
			_, err := conn.Write(m.message)

			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func removeConn(conn net.Conn) {
	var i int
	for i = range connections {
		if connections[i] == conn {
			break
		}
		connections = append(connections[:i], connections[i+1:]...)
	}
}
