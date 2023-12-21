package services

import (
	"context"
	"fmt"
	"io"

	pkgerrors "github.com/pkg/errors"

	"github.com/virsavik/sample-azure-func-app/internal/core/domain"
)

func (svc FileService) UploadFile(ctx context.Context, file io.Reader, info domain.FileInfo) error {
	// 1. Save file to Azure BLOB storage
	if err := svc.store.Save(ctx, info.Name, file); err != nil {
		return pkgerrors.WithStack(err)
	}

	// 2. Save info to DB
	if err := svc.repo.Save(ctx, info); err != nil {
		return err
	}

	// 3. Publish event save file success to SQS
	if err := svc.publisher.Publish(ctx, fmt.Sprintf("File upload successfully")); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}
