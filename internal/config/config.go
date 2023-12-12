package config

import (
	"errors"
	"fmt"
	"log"
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

type OtelConfig struct {
	ServiceName      string
	ExporterEndpoint string
}

// AppConfig representing an application configuration
type AppConfig struct {
	Web             WebConfig
	Otel            OtelConfig
	ShutdownTimeout time.Duration
}

// ReadConfigFromEnv reads all environment variables, validates it and parses it into AppConfig struct
func ReadConfigFromEnv() (AppConfig, error) {
	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")))
	if err != nil {
		return AppConfig{}, errors.New("port is invalid")
	}

	otelServiceName := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME"))
	if otelServiceName == "" {
		log.Print("open telemetry service name have not been set")
	}

	otelExporterEndpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if otelExporterEndpoint == "" {
		log.Print("open telemetry exporter endpoint have not been set")
	}

	return AppConfig{
		Web: WebConfig{
			Host: "0.0.0.0",
			Port: fmt.Sprintf(":%v", port),
		},
		Otel: OtelConfig{
			ServiceName:      otelServiceName,
			ExporterEndpoint: otelExporterEndpoint,
		},
	}, nil
}
