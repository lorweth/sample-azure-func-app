package ports

import (
	"context"
	"io"

	"github.com/virsavik/sample-azure-func-app/internal/core/domain"
)

type FileService interface {
	UploadFile(ctx context.Context, file io.Reader, info domain.FileInfo) error
}

type FileRepository interface {
	Save(ctx context.Context, info domain.FileInfo) error
}
