package storage

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	pkgerrors "github.com/pkg/errors"
)

type FileStorage struct {
	client        *azblob.Client
	containerName string
}

func New(client *azblob.Client, containerName string) FileStorage {
	return FileStorage{
		client:        client,
		containerName: containerName,
	}
}

func (b FileStorage) Save(ctx context.Context, fileName string, file io.Reader) error {
	_, err := b.client.UploadStream(ctx, b.containerName, fileName, file, nil)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}
