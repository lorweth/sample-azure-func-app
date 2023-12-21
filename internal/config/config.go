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

type MongoDBConfig struct {
	URI    string
	DBName string
}

type AzBlobConfig struct {
	ConnectionString string
	ContainerName    string
}

type AzQueueConfig struct {
	ConnectionString string
	QueueName        string
}

// AppConfig representing an application configuration
type AppConfig struct {
	Web             WebConfig
	Otel            OtelConfig
	AzBlob          AzBlobConfig
	MongoDB         MongoDBConfig
	AzQueue         AzQueueConfig
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

	mongoDBURI := strings.TrimSpace(os.Getenv("MONGODB_URI"))
	if mongoDBURI == "" {
		return AppConfig{}, errors.New("mongotel uri is invalid")
	}

	mongoDBName := strings.TrimSpace(os.Getenv("MONGODB_DBNAME"))
	if mongoDBName == "" {
		return AppConfig{}, errors.New("mongotel database name is invalid")
	}

	blobContainerName := strings.TrimSpace(os.Getenv("AZURE_BLOB_STORAGE_CONTAINER_NAME"))
	if blobContainerName == "" {
		return AppConfig{}, errors.New("azure blob container name is invalid")
	}

	blobcs := strings.TrimSpace(os.Getenv("AZURE_BLOB_STORAGE_CONNECTION_STRING"))
	if blobcs == "" {
		return AppConfig{}, errors.New("azure blob connection string is invalid")
	}

	queuecs := strings.TrimSpace(os.Getenv("AZURE_STORAGE_QUEUE_CONNECTION_STRING"))
	if queuecs == "" {
		return AppConfig{}, errors.New("azure queue connection string is invalid")
	}

	queueName := strings.TrimSpace(os.Getenv("AZURE_STORAGE_QUEUE_NAME"))
	if queueName == "" {
		return AppConfig{}, errors.New("azure queue name is invalid")
	}

	return AppConfig{
		Web: WebConfig{
			Host: "0.0.0.0",
			Port: fmt.Sprintf(":%v", port),
		},
		MongoDB: MongoDBConfig{
			URI:    mongoDBURI,
			DBName: mongoDBName,
		},
		AzBlob: AzBlobConfig{
			ContainerName:    blobContainerName,
			ConnectionString: blobcs,
		},
		AzQueue: AzQueueConfig{
			QueueName:        queueName,
			ConnectionString: queuecs,
		},
		Otel: OtelConfig{
			ServiceName:      otelServiceName,
			ExporterEndpoint: otelExporterEndpoint,
		},
	}, nil
}
