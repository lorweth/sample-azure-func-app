package system

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/virsavik/sample-azure-func-app/internal/config"
	"github.com/virsavik/sample-azure-func-app/internal/logger"
	"github.com/virsavik/sample-azure-func-app/internal/waiter"
)

type System struct {
	cfg     config.AppConfig
	logger  logger.Logger
	mux     *chi.Mux
	db      *mongo.Client
	storage *azblob.Client
	queue   *azqueue.QueueClient
	waiter  waiter.Waiter
	tp      *sdktrace.TracerProvider
}

func New(cfg config.AppConfig) (*System, error) {
	s := &System{cfg: cfg}

	s.initMux()

	s.initLogger()

	s.initWaiter()

	if err := s.initOpenTelemetry(); err != nil {
		return nil, err
	}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initBlobStorage(); err != nil {
		return nil, err
	}

	if err := s.initStorageQueue(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *System) Config() config.AppConfig {
	return s.cfg
}

func (s *System) initMux() {
	s.mux = chi.NewMux()
}

func (s *System) Mux() *chi.Mux {
	return s.mux
}

func (s *System) initLogger() {
	s.logger = logger.New()
}

func (s *System) Logger() logger.Logger {
	return s.logger
}

func (s *System) initDB() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(s.Config().MongoDB.URI))
	if err != nil {
		return err
	}

	s.db = client

	s.Waiter().Cleanup(func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("mongdo DB disconnect error: %v", err)
		}
	})

	return nil
}

func (s *System) DB() *mongo.Client {
	return s.db
}

func (s *System) initBlobStorage() error {
	client, err := azblob.NewClientFromConnectionString(s.Config().AzBlob.ConnectionString, nil)
	if err != nil {
		return fmt.Errorf("create azure blob client error: %w", err)
	}

	s.storage = client
	return nil
}

func (s *System) BlobStorage() *azblob.Client {
	return s.storage
}

func (s *System) initStorageQueue() error {
	client, err := azqueue.NewQueueClientFromConnectionString(s.Config().AzQueue.ConnectionString, s.Config().AzQueue.QueueName, nil)
	if err != nil {
		return fmt.Errorf("create azure queue client error: %w", err)
	}

	s.queue = client

	return nil
}

func (s *System) StorageQueue() *azqueue.QueueClient {
	return s.queue
}

// initOpenTelemetry Initializes an OTLP exporter, and configures the corresponding trace
// and metric providers.
// copy from sample https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go
func (s *System) initOpenTelemetry() error {
	// Skip when config is not provided
	if s.cfg.Otel.ServiceName == "" || s.cfg.Otel.ExporterEndpoint == "" {
		return nil
	}

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(s.cfg.Otel.ServiceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Set up a trace exporter
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(s.cfg.Otel.ExporterEndpoint),
	)
	if err != nil {
		return err
	}

	// Register the trace exporter with a TracerProvider, using
	// a batch span processor to aggregate spans before export.
	s.tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(s.tp)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	s.waiter.Cleanup(func() {
		if err := s.tp.Shutdown(context.Background()); err != nil {
			s.logger.Errorf(err, "ran into an issue shutting down the tracer provider")
		}
	})

	return nil
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Printf("web server started; listening at http://localhost%s\n", s.cfg.Web.Port)
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("web server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	return group.Wait()
}
