package ports

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, msg string) error
}
