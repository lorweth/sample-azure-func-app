package httpio

import (
	"encoding/json"
	"net/http"

	"github.com/virsavik/sample-azure-func-app/internal/logger"
)

// Response represents an HTTP response structure with a generic body 'T'.
// It includes the HTTP status code, headers, and the response body.
// The type 'T' can be any type, making it flexible for different response structures.
type Response[T any] struct {
	Status  int
	Headers map[string]string
	Body    T
}

// WriteJSON writes a JSON response to the provided http.ResponseWriter.
// It takes an HTTP request, a Response structure containing the response data,
// and encodes the data into JSON format, updating the HTTP status code and headers accordingly.
func WriteJSON[T any](w http.ResponseWriter, r *http.Request, data Response[T]) {
	// Prepare HTTP response headers
	w.Header().Set("Content-Type", "application/json")
	for key, val := range data.Headers {
		w.Header().Set(key, val)
	}

	// Update the HTTP status code
	w.WriteHeader(data.Status)

	// Encode the response body to JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(data.Body); err != nil {
		logger.FromCtx(r.Context()).Errorf(err, "json encode error")
	}
}
