package main

import (
	"EventSpace/internal/server"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func main() {

	doc, err := os.ReadFile("config/api-gateway/config.toml")

	if err != nil {
		log.Fatalf("Error to load api-gateway config: %v", err)
	}

	var config server.Config
	err = toml.Unmarshal(doc, &config)

	if err != nil {
		log.Fatalf("Error to parse api-gateway config: %v", err)
	}

	server := server.NewServer(&config)

	if err = server.Start(); err != nil {
		log.Fatalf("Error to start api-gateway: %v", err)
	}
}
