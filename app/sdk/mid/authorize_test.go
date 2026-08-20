package mid

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

func authorizeTestLogger() *logger.Logger {
	return logger.New(io.Discard, logger.LevelInfo, "TEST", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })
}

func authorizeOKHandler() web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func authorizeTestClaims() auth.Claims {
	return auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: uuid.NewString(),
		},
		Roles: []string{"USER"},
	}
}

// Test_AuthorizePropagatesForbidden proves a 403 from the auth service
// (authenticated but not allowed) reaches the caller as PermissionDenied
// instead of being collapsed to 401, which would log the user out.
func Test_AuthorizePropagatesForbidden(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(errs.Newf(errs.PermissionDenied, "not allowed"))
	}))
	defer srv.Close()

	cln := authclient.New(authorizeTestLogger(), srv.URL)
	handler := Authorize(cln, auth.RuleAdminOnly)(authorizeOKHandler())

	ctx := setClaims(context.Background(), authorizeTestClaims())
	req := httptest.NewRequest(http.MethodGet, "/v1/orders", nil)
	rec := httptest.NewRecorder()

	err := handler(ctx, rec, req)

	var appErr *errs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type = %T (%s), want *errs.Error", err, err)
	}

	if !appErr.Code.Equal(errs.PermissionDenied) {
		t.Errorf("code = %s, want %s", appErr.Code, errs.PermissionDenied)
	}
}

// Test_AuthorizeTransportFailureFailsClosed proves a transport failure to
// the auth service denies the request (401) rather than letting it through.
func Test_AuthorizeTransportFailureFailsClosed(t *testing.T) {
	t.Parallel()

	cln := authclient.New(authorizeTestLogger(), "http://127.0.0.1:1")
	handler := Authorize(cln, auth.RuleAdminOnly)(authorizeOKHandler())

	ctx := setClaims(context.Background(), authorizeTestClaims())
	req := httptest.NewRequest(http.MethodGet, "/v1/orders", nil)
	rec := httptest.NewRecorder()

	err := handler(ctx, rec, req)

	var appErr *errs.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error type = %T (%s), want *errs.Error", err, err)
	}

	if !appErr.Code.Equal(errs.Unauthenticated) {
		t.Errorf("code = %s, want %s", appErr.Code, errs.Unauthenticated)
	}
}
