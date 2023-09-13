package system

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/virsavik/sample-azure-func-app/internal/config"
	"github.com/virsavik/sample-azure-func-app/internal/logger"
)

// Service representing an application service
type Service interface {
	Config() config.AppConfig

	Logger() logger.Logger

	Mux() *chi.Mux
}

// Module representing an application module
type Module interface {
	Startup(context.Context, Service) error
}
