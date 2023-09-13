package httpio

import (
	"errors"
	"net/http"
)

// HandlerFunc wraps an HTTP handler function that returns an error.
// It adds OpenTelemetry tracing and handles specific error types by responding with JSON.
// If the error is of type httpio.Error, a custom JSON response is generated.
// If the error is not of type httpio.Error, a generic internal server error response is generated.
func HandlerFunc(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute the provided handler function and handle errors
		if err := fn(w, r); err != nil {
			var apiErr Error
			if errors.As(err, &apiErr) {
				WriteJSON(w, r, Response[Message]{
					Status: apiErr.Status,
					Body: Message{
						Code: apiErr.Code,
						Desc: apiErr.Desc,
					},
				})

				return
			}

			// If the error is not of type "Error", respond with a generic internal server error
			WriteJSON(w, r, Response[Message]{
				Status: http.StatusInternalServerError,
				Body:   MsgInternalServerError,
			})
		}
	}
}
