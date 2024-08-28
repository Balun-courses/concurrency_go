package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"spider/internal/configuration"
	"spider/internal/initialization"
)

var (
	ConfigFileName = os.Getenv("CONFIG_FILE_NAME")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := &configuration.Config{}
	if ConfigFileName != "" {
		data, err := os.ReadFile(ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err = configuration.Load(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	initializer, err := initialization.NewInitializer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = initializer.StartDatabase(ctx); err != nil {
		log.Fatal(err)
	}
}
