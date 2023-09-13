package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// waiter implementation of the Waiter interface
type waiter struct {
	ctx          context.Context
	waitFuncs    []WaitFunc
	cleanupFuncs []CleanupFunc
	cancel       context.CancelFunc
}

type waiterCfg struct {
	parentCtx    context.Context
	catchSignals bool
}

func New(options ...Option) Waiter {
	cfg := &waiterCfg{
		parentCtx:    context.Background(),
		catchSignals: false,
	}

	for _, option := range options {
		option(cfg)
	}

	w := &waiter{
		waitFuncs:    []WaitFunc{},
		cleanupFuncs: []CleanupFunc{},
	}

	w.ctx, w.cancel = context.WithCancel(cfg.parentCtx)
	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return w
}

// Add a function to wait for
func (w *waiter) Add(fns ...WaitFunc) {
	w.waitFuncs = append(w.waitFuncs, fns...)
}

// Cleanup with a function after waiting
func (w *waiter) Cleanup(fns ...CleanupFunc) {
	w.cleanupFuncs = append(w.cleanupFuncs, fns...)
}

// Wait for all added functions to complete
func (w *waiter) Wait() (err error) {
	g, ctx := errgroup.WithContext(w.ctx)

	g.Go(func() error {
		<-ctx.Done()
		w.cancel()
		return nil
	})

	for _, fn := range w.waitFuncs {
		waitFunc := fn // Avoid closure capture
		g.Go(func() error { return waitFunc(ctx) })
	}

	for _, fn := range w.cleanupFuncs {
		cleanupFunc := fn // Avoid closure capture
		defer cleanupFunc()
	}

	return g.Wait()
}

func (w *waiter) Context() context.Context {
	return w.ctx
}

func (w *waiter) CancelFunc() context.CancelFunc {
	return w.cancel
}
