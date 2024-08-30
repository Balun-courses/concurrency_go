package main

import (
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
	defer func() {
		if v := recover(); v != nil {
			log.Println("captured panic:", v)
		}
		c.Close()
	}()

	panic("internal error")
}
