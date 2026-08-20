package authapi_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"
	"time"

	authapi "github.com/warlck/food-flow/app/domain/authapp"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/mux"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/password"
	"github.com/warlck/food-flow/business/types/role"
	"github.com/warlck/food-flow/foundation/web"
)

const (
	testEmail    = "login-admin@example.com"
	testPassword = "TestPassword123"
	testKID      = "test-kid"
	testIssuer   = "food-flow-auth"
	testTTL      = 8 * time.Hour
)

func Test_Login(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Login")

	ath := auth.New(auth.Config{
		Log:           db.Log,
		KeyLookup:     newTestKeyStore(t),
		UserBus:       db.BusDomain.User,
		Issuer:        testIssuer,
		ActiveKID:     testKID,
		TokenTTL:      testTTL,
		LoginMaxFails: 3,
		LoginLockout:  15 * time.Minute,
	})

	app := mux.WebAPI(mux.Config{
		Log:  db.Log,
		Auth: ath,
		DB:   db.DB,
		BusConfig: mux.BusConfig{
			UserBus: db.BusDomain.User,
		},
	}, routeAdderFunc(func(app *web.App, cfg mux.Config) {
		authapi.Routes(app, authapi.Config{
			Auth:    cfg.Auth,
			UserBus: cfg.UserBus,
		})
	}))

	ctx := context.Background()

	admin := seedUser(t, ctx, db.BusDomain.User, testEmail, role.Admin, true)
	seedUser(t, ctx, db.BusDomain.User, "login-user@example.com", role.User, true)
	seedUser(t, ctx, db.BusDomain.User, "login-disabled@example.com", role.Admin, false)
	seedUser(t, ctx, db.BusDomain.User, "login-lockout@example.com", role.Admin, true)

	t.Run("valid admin credentials", testLoginSuccess(app, ath, admin))
	t.Run("wrong password", testLoginFailure(app, `{"email":"`+testEmail+`","password":"**************"}`, http.StatusUnauthorized, "invalid credentials"))
	t.Run("unknown email", testLoginFailure(app, `{"email":"nobody@example.com","password":"`+testPassword+`"}`, http.StatusUnauthorized, "invalid credentials"))
	t.Run("disabled user", testLoginFailure(app, `{"email":"login-disabled@example.com","password":"`+testPassword+`"}`, http.StatusUnauthorized, "invalid credentials"))
	t.Run("malformed email", testLoginFailure(app, `{"email":"not-an-email","password":"`+testPassword+`"}`, http.StatusUnauthorized, "invalid credentials"))
	t.Run("non-admin user", testLoginFailure(app, `{"email":"login-user@example.com","password":"`+testPassword+`"}`, http.StatusForbidden, "admin role required"))
	t.Run("malformed json", testLoginFailure(app, `{`, http.StatusBadRequest, ""))
	t.Run("missing password", testLoginFailure(app, `{"email":"`+testEmail+`"}`, http.StatusBadRequest, ""))
	t.Run("lockout after repeated failures", testLoginLockout(app))
}

func testLoginLockout(app *web.App) func(t *testing.T) {
	return func(t *testing.T) {
		wrongPassword := `{"email":"login-lockout@example.com","password":"WrongPassw0rd"}`

		// The configured test limit is 3 consecutive failures.
		for i := range 3 {
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(wrongPassword)))

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("attempt %d: status = %d, want %d, body: %s", i+1, rec.Code, http.StatusUnauthorized, rec.Body.String())
			}
		}

		// The next attempts are locked out, even with the correct password.
		for _, body := range []string{
			wrongPassword,
			`{"email":"login-lockout@example.com","password":"` + testPassword + `"}`,
		} {
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body)))

			if rec.Code != http.StatusTooManyRequests {
				t.Fatalf("locked out: status = %d, want %d, body: %s", rec.Code, http.StatusTooManyRequests, rec.Body.String())
			}

			var errResp struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
				t.Fatalf("decoding error response: %s", err)
			}

			if errResp.Message != "too many attempts" {
				t.Errorf("message = %q, want %q", errResp.Message, "too many attempts")
			}
		}
	}
}

func testLoginSuccess(app *web.App, ath *auth.Auth, admin userbus.User) func(t *testing.T) {
	return func(t *testing.T) {
		body := `{"email":"` + testEmail + `","password":"` + testPassword + `"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
		}

		var resp struct {
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expiresAt"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("decoding response: %s", err)
		}

		if resp.Token == "" {
			t.Fatal("expected a token in the response")
		}

		if ttl := time.Until(resp.ExpiresAt); ttl < testTTL-time.Minute || ttl > testTTL {
			t.Errorf("token lifetime = %s, want approximately %s", ttl, testTTL)
		}

		// The issued token must authenticate end to end: signature, issuer,
		// expiry, and the user-still-enabled check.
		claims, err := ath.Authenticate(context.Background(), "Bearer "+resp.Token)
		if err != nil {
			t.Fatalf("issued token failed authentication: %s", err)
		}

		if claims.Subject != admin.ID.String() {
			t.Errorf("subject = %q, want %q", claims.Subject, admin.ID.String())
		}

		if claims.Issuer != testIssuer {
			t.Errorf("issuer = %q, want %q", claims.Issuer, testIssuer)
		}

		if len(claims.Roles) != 1 || claims.Roles[0] != "ADMIN" {
			t.Errorf("roles = %v, want [ADMIN]", claims.Roles)
		}
	}
}

func testLoginFailure(app *web.App, body string, expStatus int, expMessage string) func(t *testing.T) {
	return func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		if rec.Code != expStatus {
			t.Fatalf("status = %d, want %d, body: %s", rec.Code, expStatus, rec.Body.String())
		}

		if expMessage == "" {
			return
		}

		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
			t.Fatalf("decoding error response: %s", err)
		}

		if errResp.Message != expMessage {
			t.Errorf("message = %q, want %q", errResp.Message, expMessage)
		}
	}
}

// =============================================================================

func seedUser(t *testing.T, ctx context.Context, userBus *userbus.Business, email string, rle role.Role, enabled bool) userbus.User {
	t.Helper()

	usr, err := userBus.Create(ctx, userbus.NewUser{
		Name:     name.MustParse("Login Test"),
		Email:    mail.Address{Address: email},
		Roles:    []role.Role{rle},
		Password: password.MustParse(testPassword),
	})
	if err != nil {
		t.Fatalf("creating user %s: %s", email, err)
	}

	if !enabled {
		usr, err = userBus.Update(ctx, usr, userbus.UpdateUser{Enabled: dbtest.BoolPointer(false)})
		if err != nil {
			t.Fatalf("disabling user %s: %s", email, err)
		}
	}

	return usr
}

// =============================================================================

type routeAdderFunc func(*web.App, mux.Config)

func (fn routeAdderFunc) Add(app *web.App, cfg mux.Config) {
	fn(app, cfg)
}

// =============================================================================

// testKeyStore generates a throwaway RSA key pair per test run so no key
// material needs to be embedded in the test source.
type testKeyStore struct {
	privateKey string
	publicKey  string
}

func newTestKeyStore(t *testing.T) *testKeyStore {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating test key: %s", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshaling private key: %s", err)
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("marshaling public key: %s", err)
	}

	return &testKeyStore{
		privateKey: string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})),
		publicKey:  string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})),
	}
}

func (ks *testKeyStore) PrivateKey(kid string) (string, error) {
	return ks.privateKey, nil
}

func (ks *testKeyStore) PublicKey(kid string) (string, error) {
	return ks.publicKey, nil
}
