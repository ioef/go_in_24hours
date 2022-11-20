package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	server, err := net.Listen("tcp", ":8080")
	fmt.Printf("[+] Listening on localhost TCP port 8080. Press Ctrl+C to Terminate.\n")
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	io.Copy(conn, conn)
}
