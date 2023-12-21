package v1

import (
	"net/http"

	"github.com/virsavik/sample-azure-func-app/internal/core/domain"
	"github.com/virsavik/sample-azure-func-app/internal/httpio"
)

func (h Handler) Upload() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		file, header, err := r.FormFile("file")
		if err != nil {
			return err
		}

		if err := h.fileService.UploadFile(ctx, file, domain.FileInfo{
			Name:     header.Filename,
			Metadata: "",
		}); err != nil {
			// TODO: convert service error
			return err
		}

		httpio.WriteJSON(w, r, httpio.Response[httpio.Message]{
			Status: http.StatusOK,
			Body: httpio.Message{
				Code: "upload_success",
				Desc: "Upload file successfully",
			},
		})

		return nil
	})
}
