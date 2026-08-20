package mid

import (
	"context"
	"net/http"
	"slices"

	"github.com/warlck/food-flow/foundation/web"
)

// CORS handles cross-origin browser requests according to the configured
// allowed origins. In the deployed topology the frontends proxy API calls
// same-origin through nginx, so browsers rarely need CORS at all; this is
// defense in depth for direct browser access to a service. A single "*"
// entry allows any origin (local development). An empty list emits no CORS
// headers, which makes browsers block all cross-origin reads.
func CORS(allowedOrigins []string) web.MidHandler {
	allowAny := slices.Contains(allowedOrigins, "*")

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return handler(ctx, w, r)
			}

			if !allowAny && !slices.Contains(allowedOrigins, origin) {
				// Disallowed preflights are rejected outright; disallowed
				// simple requests run but get no CORS headers, so the
				// browser blocks the response from the calling script.
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusForbidden)
					return nil
				}

				return handler(ctx, w, r)
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return nil
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
