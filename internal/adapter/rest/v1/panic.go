package v1

import (
	"errors"
	"net/http"

	"github.com/virsavik/sample-azure-func-app/internal/httpio"
)

func (h Handler) Panic() http.HandlerFunc {
	return httpio.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		panic(errors.New("sample panic exception"))

		return nil
	})
}
