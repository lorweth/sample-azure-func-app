package v1

import (
	"net/http"

	"github.com/virsavik/sample-azure-func-app/internal/httpio"
)

func (h Handler) Ping() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		httpio.WriteJSON(w, r, httpio.Response[httpio.Message]{
			Status: http.StatusOK,
			Body: httpio.Message{
				Code: "connected",
				Desc: "Connected",
			},
		})
		return nil
	})
}
