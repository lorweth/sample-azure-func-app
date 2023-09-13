package system

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/virsavik/sample-azure-func-app/internal/config"
	"github.com/virsavik/sample-azure-func-app/internal/logger"
	"github.com/virsavik/sample-azure-func-app/internal/waiter"
)

type System struct {
	cfg    config.AppConfig
	logger logger.Logger
	mux    *chi.Mux
	waiter waiter.Waiter
}

func New(cfg config.AppConfig) (*System, error) {
	s := &System{cfg: cfg}

	s.initMux()

	s.initLogger()

	s.initWaiter()

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
