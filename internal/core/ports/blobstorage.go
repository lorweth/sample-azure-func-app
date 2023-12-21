package ports

import (
	"context"
	"io"
)

type BlobStorage interface {
	Save(ctx context.Context, fileName string, file io.Reader) error
}
