package services

import (
	"github.com/virsavik/sample-azure-func-app/internal/core/ports"
)

type FileService struct {
	repo      ports.FileRepository
	publisher ports.Publisher
	store     ports.BlobStorage
}

func New(repo ports.FileRepository, publisher ports.Publisher, store ports.BlobStorage) FileService {
	return FileService{
		repo:      repo,
		publisher: publisher,
		store:     store,
	}
}
