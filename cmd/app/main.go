package main

import (
	"flag"
	"log"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/app"
)

func main() {
	configFilePath := flag.String("config", "", "Path to the config file (JSON or YAML)")
	flag.Parse()

	if *configFilePath == "" {
		log.Fatalf("No config file provided. Use -config <file_path> to specify the config file.")
	}

	config, err := config.Load(*configFilePath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app.Run(config)
}
