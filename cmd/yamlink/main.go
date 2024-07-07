package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/NishantJoshi00/yamlink"
	"gopkg.in/yaml.v3"
)

func main() {
	yamlink.Logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	config_file, err1 := os.LookupEnv("CONFIG_FILE")

	if err1 != true {
		os.Exit(1)
	}

	if _, err := os.Stat(config_file); err != nil {
		yamlink.Logger.Error("Failed while loading config file", err)
	}

	file, err := os.Open(config_file)

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var config yamlink.Config

	err = decoder.Decode(&config)

	if err != nil {
		yamlink.Logger.Error("Failed while decoding config file", err)
		os.Exit(1)
	}

	yamlink.Logger.Info("Config file loaded successfully")

	server, err := yamlink.Init(&config)

	if err != nil {
		yamlink.Logger.Error("Failed while initializing server", err)
	}

	yamlink.Logger.Info(fmt.Sprintf("Starting server on %s:%d", config.Host, config.Port))

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), server)

	if err != nil {
		yamlink.Logger.Error("Failed while starting server", err)
		os.Exit(1)
	}
}
