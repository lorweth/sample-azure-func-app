package waiter

import (
	"context"
)

// WaitFunc a function to wait for
type WaitFunc func(ctx context.Context) error

// CleanupFunc a function after waiting
type CleanupFunc func()

// Waiter is an interface that allows you to manage and control a set of functions,
// enabling you to add functions to wait for, specify cleanup functions to be executed
// after waiting, and initiate the waiting process. It also provides access to the
// underlying context and cancel function, which can be used to control the waiting
// process and handle cancellations.
type Waiter interface {
	// Add a function to wait for
	Add(fns ...WaitFunc)

	// Cleanup with a function after waiting
	Cleanup(fns ...CleanupFunc)

	// Wait for all added functions to complete
	Wait() error

	// Context return the context
	Context() context.Context

	// CancelFunc return the context cancel function
	CancelFunc() context.CancelFunc
}
