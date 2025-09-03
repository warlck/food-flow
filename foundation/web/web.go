// Package web contains a small web framework extension.
package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, args ...any)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	log Logger
	*http.ServeMux
	mw []MidHandler
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(log Logger, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		mw:       mw,
		log:      log,
	}
}

// HandleFunc sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandleFunc(method string, group string, path string, handler Handler, mw ...MidHandler) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := a.handle(handler)

	a.ServeMux.HandleFunc(finalPath(method, group, path), h)
}

// HandleFuncNoMiddleware sets a handler function for a given HTTP method and
// path pair to the application server mux with no middleware.
func (a *App) HandleFuncNoMiddleware(method string, group string, path string, handler Handler, mw ...MidHandler) {
	h := a.handle(handler)
	a.ServeMux.HandleFunc(finalPath(method, group, path), h)
}

// handle is the function that wraps the handler function and adds the values to the context
// and validates the error.
func (a *App) handle(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			a.log(ctx, "web", "ERROR", err.Error())
			return
		}
	}
}

func finalPath(method string, group string, path string) string {
	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	return fmt.Sprintf("%s %s", method, finalPath)
}
