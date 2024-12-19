package main

import (
	"errors"
	"log"
	"net"
)

// nc 127.0.0.1 12345

func main() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go ClientHandler(conn)
	}
}

func ClientHandler(c net.Conn) {
	panic(errors.New("internal error"))
}
