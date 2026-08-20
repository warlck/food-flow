package authclient_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/foundation/logger"
)

func testLogger() *logger.Logger {
	return logger.New(io.Discard, logger.LevelInfo, "TEST", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })
}

// Test_UnreachableFailsClosed proves an unreachable auth service yields an
// error (which middleware maps to a rejection), never a successful call.
func Test_UnreachableFailsClosed(t *testing.T) {
	t.Parallel()

	// Port 1 is reserved and refuses connections.
	cln := authclient.New(testLogger(), "http://127.0.0.1:1")

	if _, err := cln.Authenticate(context.Background(), "Bearer token"); err == nil {
		t.Fatal("authenticate against an unreachable auth service must fail")
	}

	if err := cln.Authorize(context.Background(), authclient.Authorize{}); err == nil {
		t.Fatal("authorize against an unreachable auth service must fail")
	}
}

// Test_ErrorMapping proves 401 and 403 responses from the auth service are
// decoded into typed errors with their codes preserved, so middleware can
// distinguish "not authenticated" from "not allowed".
func Test_ErrorMapping(t *testing.T) {
	t.Parallel()

	table := []struct {
		name   string
		status int
		code   errs.ErrCode
	}{
		{"401 maps to unauthenticated", http.StatusUnauthorized, errs.Unauthenticated},
		{"403 maps to permission denied", http.StatusForbidden, errs.PermissionDenied},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				json.NewEncoder(w).Encode(errs.Newf(tt.code, "test error"))
			}))
			defer srv.Close()

			cln := authclient.New(testLogger(), srv.URL)

			err := cln.Authorize(context.Background(), authclient.Authorize{})

			var appErr *errs.Error
			if !errors.As(err, &appErr) {
				t.Fatalf("error type = %T (%s), want *errs.Error", err, err)
			}

			if !appErr.Code.Equal(tt.code) {
				t.Errorf("code = %s, want %s", appErr.Code, tt.code)
			}
		})
	}
}

// Test_Timeout proves the client enforces its overall request timeout so a
// hung auth service cannot tie up callers indefinitely.
func Test_Timeout(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cln := authclient.New(testLogger(), srv.URL, authclient.WithTimeout(50*time.Millisecond))

	start := time.Now()
	if err := cln.Authorize(context.Background(), authclient.Authorize{}); err == nil {
		t.Fatal("request to a hung auth service must time out")
	}

	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("timeout took %s, want ~50ms", elapsed)
	}
}
