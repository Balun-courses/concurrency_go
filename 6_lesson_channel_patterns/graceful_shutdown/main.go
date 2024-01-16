package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", "127.0.0.1:3322")
	if err != nil {
		log.Fatal(err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			connection, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break
				}

				log.Print(err.Error())
				continue
			}

			request := make([]byte, 2048)
			_, err = connection.Read(request)
			if err != nil {
				if err != io.EOF {
					log.Print(err.Error())
				}

				break
			}

			log.Println(fmt.Sprintf("request: %s", string(request)))
			if _, err := connection.Write([]byte("Hello from server")); err != nil {
				log.Print(err.Error())
				break
			}
		}

		log.Print("Server was stopped")
	}()

	<-interruptCh
	if err := listener.Close(); err != nil {
		log.Print(err.Error())
	}

	wg.Wait()
	log.Print("Application was stopped")
}
