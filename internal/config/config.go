package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// WebConfig representing a web configuration
type WebConfig struct {
	Host string
	Port string
}

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s%v", c.Host, c.Port)
}

// AppConfig representing an application configuration
type AppConfig struct {
	Web             WebConfig
	ShutdownTimeout time.Duration
}

// ReadConfigFromEnv reads all environment variables, validates it and parses it into AppConfig struct
func ReadConfigFromEnv() (AppConfig, error) {
	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")))
	if err != nil {
		return AppConfig{}, errors.New("port is invalid")
	}

	return AppConfig{
		Web: WebConfig{
			Host: "0.0.0.0",
			Port: fmt.Sprintf(":%v", port),
		},
	}, nil
}
