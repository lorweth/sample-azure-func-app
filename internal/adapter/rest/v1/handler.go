package v1

import (
	"github.com/virsavik/sample-azure-func-app/internal/core/services"
)

type Handler struct {
	fileService services.FileService
}

func New(fileService services.FileService) Handler {
	return Handler{
		fileService: fileService,
	}
}
