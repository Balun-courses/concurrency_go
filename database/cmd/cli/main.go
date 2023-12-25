package main

import (
	"bufio"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"spider/internal/network"
	"spider/internal/tools"
	"time"
)

func main() {
	address := flag.String("address", "localhost:3223", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger, _ := zap.NewProduction()
	maxMessageSize, err := tools.ParseSize(*maxMessageSizeStr)
	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTCPClient(*address, maxMessageSize, *idleTimeout)
	if err != nil {
		logger.Fatal("failed to connect with server", zap.Error(err))
	}

	for {
		fmt.Print("[spider] > ")
		request, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("failed to read user query", zap.Error(err))
		}

		response, err := client.Send([]byte(request))
		if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(response))
	}
}
