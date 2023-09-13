package httpio

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/virsavik/sample-azure-func-app/internal/logger"
)

// Middleware is a customizable HTTP middleware that provides logging, panic recovery,
// and response wrapping capabilities.
func Middleware(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// The main middleware function that will be executed for each HTTP request.
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Wrap response for getting response status and body size.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				// Recover from the panic and obtain the panic value.
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%+v", p)
					}

					// Capture and log the entire stack trace along with the error details.
					log.Errorf(err, "caught a panic, stacktrace: %s", debug.Stack())

					// Respond with a 500 Internal Server Error and log any encoding errors.
					WriteJSON(ww, r, Response[Message]{
						Status: http.StatusInternalServerError,
						Body:   MsgInternalServerError,
					})
				}
			}()

			defer func() {
				// Log when the function finishes processing the request.
				reqLogger := log.With(
					logger.String("host.name", r.Host),
					logger.String("url.path", r.URL.Path),
					logger.String("url.query", r.URL.RawQuery),
					logger.String("http.request.method_original", r.Method),
					logger.Int("http.request.body.size", int(r.ContentLength)),
					logger.String("http.request.proto", r.Proto),
					logger.String("http.request.remote_address", r.RemoteAddr),
					logger.String("user_agent.original", r.UserAgent()),
					logger.Int("http.response.status_code", ww.Status()),
					logger.Int("http.response.body.size", ww.BytesWritten()),
				)

				reqLogger.Infof("Served")
			}()

			// Handle the original request by calling the next handler in the chain,
			// passing the wrapped response writer and the modified request context.
			next.ServeHTTP(ww, r.WithContext(logger.SetInCtx(r.Context(), log)))
		}

		// Return an HTTP handler function wrapping the main middleware function.
		return http.HandlerFunc(fn)
	}
}
