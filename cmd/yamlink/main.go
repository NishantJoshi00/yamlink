package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/NishantJoshi00/yamlink"
	"gopkg.in/yaml.v3"
)

func logConfig() *slog.HandlerOptions {
	logLevelEnv := os.Getenv("LOG_LEVEL")

	// Parse log level
	var level slog.Level
	switch strings.ToUpper(logLevelEnv) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		fmt.Printf("Invalid log level: %s. Defaulting to INFO.\n", logLevelEnv)
		level = slog.LevelInfo
	}

	return &slog.HandlerOptions{
		Level: level,
	}
}

func main() {
	lconfig := logConfig()

	yamlink.Logger = slog.New(slog.NewJSONHandler(os.Stderr, lconfig))

	config_file, err1 := os.LookupEnv("CONFIG_FILE")

	if err1 != true {
		os.Exit(1)
	}

	if _, err := os.Stat(config_file); err != nil {
		yamlink.Logger.Error("Failed while loading config file", "error", err)
	}

	file, err := os.Open(config_file)

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var config yamlink.Config

	err = decoder.Decode(&config)

	if err != nil {
		yamlink.Logger.Error("Failed while decoding config file", "error", err)
		os.Exit(1)
	}

	yamlink.Logger.Info("Config file loaded successfully")

	server, err := yamlink.Init(&config)

	if err != nil {
		yamlink.Logger.Error("Failed while initializing server", "error", err)
	}

	yamlink.Logger.Info(fmt.Sprintf("Starting server on %s:%d", config.Host, config.Port))

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), server)

	if err != nil {
		yamlink.Logger.Error("Failed while starting server", "error", err)
		os.Exit(1)
	}
}
