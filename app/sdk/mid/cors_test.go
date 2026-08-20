package mid_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/foundation/web"
)

func corsTestHandler() web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func Test_CORS(t *testing.T) {
	t.Parallel()

	allowed := []string{"https://admin.example.com"}

	table := []struct {
		name           string
		origins        []string
		method         string
		origin         string
		expStatus      int
		expAllowOrigin string
	}{
		{"no origin header passes through without CORS headers", allowed, http.MethodGet, "", http.StatusOK, ""},
		{"allowed origin reflected", allowed, http.MethodGet, "https://admin.example.com", http.StatusOK, "https://admin.example.com"},
		{"disallowed origin gets no CORS headers", allowed, http.MethodGet, "https://evil.example.com", http.StatusOK, ""},
		{"allowed preflight answers 204", allowed, http.MethodOptions, "https://admin.example.com", http.StatusNoContent, "https://admin.example.com"},
		{"disallowed preflight rejected", allowed, http.MethodOptions, "https://evil.example.com", http.StatusForbidden, ""},
		{"wildcard allows any origin", []string{"*"}, http.MethodGet, "https://anything.example.com", http.StatusOK, "https://anything.example.com"},
		{"empty config blocks all origins", nil, http.MethodGet, "https://admin.example.com", http.StatusOK, ""},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			handler := mid.CORS(tt.origins)(corsTestHandler())

			req := httptest.NewRequest(tt.method, "/v1/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			rec := httptest.NewRecorder()

			if err := handler(context.Background(), rec, req); err != nil {
				t.Fatalf("handler error: %s", err)
			}

			if rec.Code != tt.expStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.expStatus)
			}

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tt.expAllowOrigin {
				t.Errorf("Allow-Origin = %q, want %q", got, tt.expAllowOrigin)
			}

			if tt.method == http.MethodOptions && tt.expStatus == http.StatusNoContent {
				if got := rec.Header().Get("Access-Control-Allow-Headers"); got == "" {
					t.Error("preflight should set Allow-Headers")
				}
			}
		})
	}
}
