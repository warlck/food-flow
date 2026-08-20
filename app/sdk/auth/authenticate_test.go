package auth_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/warlck/food-flow/app/sdk/auth"
)

// strictKeyStore errors on unknown kids, unlike the permissive keyStore in
// auth_test.go, so unknown-kid rejection is actually exercised.
type strictKeyStore struct{}

func (ks *strictKeyStore) PrivateKey(k string) (string, error) {
	if k != kid {
		return "", errors.New("unknown kid")
	}

	return privateKeyPEM, nil
}

func (ks *strictKeyStore) PublicKey(k string) (string, error) {
	if k != kid {
		return "", errors.New("unknown kid")
	}

	return publicKeyPEM, nil
}

func tokenClaims(issuer string, exp time.Time) auth.Claims {
	return auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   "5cf37266-3473-4006-984f-9325122678b7",
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
		},
		Roles: []string{"ADMIN"},
	}
}

func signToken(t *testing.T, method jwt.SigningMethod, kidHeader string, claims auth.Claims, key any) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = kidHeader

	str, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("signing token: %s", err)
	}

	return str
}

// Test_AuthenticateRejects proves the token verification chain (parser
// pinned to RS256 plus the OPA io.jwt.decode_verify policy) rejects the
// classic JWT attack classes.
func Test_AuthenticateRejects(t *testing.T) {
	log := newUnit(t)

	ath := auth.New(auth.Config{
		Log:       log,
		KeyLookup: &strictKeyStore{},
		Issuer:    "foodflow.test",
		ActiveKID: kid,
		TokenTTL:  8 * time.Hour,
	})

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		t.Fatalf("parsing private key: %s", err)
	}

	validClaims := func() auth.Claims {
		return tokenClaims("foodflow.test", time.Now().UTC().Add(time.Hour))
	}

	t.Run("valid token accepted", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodRS256, kid, validClaims(), rsaKey)
		if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err != nil {
			t.Fatalf("valid token should authenticate: %s", err)
		}
	})

	t.Run("alg none rejected", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodNone, kid, validClaims(), jwt.UnsafeAllowNoneSignatureType)
		if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err == nil {
			t.Fatal("alg=none token must be rejected")
		}
	})

	t.Run("wrong alg HS256 rejected", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodHS256, kid, validClaims(), []byte("hmac-secret"))
		if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err == nil {
			t.Fatal("HS256-signed token must be rejected")
		}
	})

	t.Run("unknown kid rejected", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodRS256, "no-such-kid", validClaims(), rsaKey)
		if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err == nil {
			t.Fatal("token with unknown kid must be rejected")
		}
	})

	t.Run("wrong issuer rejected", func(t *testing.T) {
		claims := tokenClaims("evil-issuer", time.Now().UTC().Add(time.Hour))
		token := signToken(t, jwt.SigningMethodRS256, kid, claims, rsaKey)
		if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err == nil {
			t.Fatal("token with wrong issuer must be rejected")
		}
	})

	t.Run("expired token rejected", func(t *testing.T) {
		claims := tokenClaims("foodflow.test", time.Now().UTC().Add(-time.Hour))
		token := signToken(t, jwt.SigningMethodRS256, kid, claims, rsaKey)
		if _, err := ath.Authenticate(context.Background(), "Bearer "+token); err == nil {
			t.Fatal("expired token must be rejected")
		}
	})

	t.Run("tampered signature rejected", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodRS256, kid, validClaims(), rsaKey)

		parts := strings.Split(token, ".")
		sig := []byte(parts[2])
		if sig[0] == 'A' {
			sig[0] = 'B'
		} else {
			sig[0] = 'A'
		}
		parts[2] = string(sig)

		if _, err := ath.Authenticate(context.Background(), "Bearer "+strings.Join(parts, ".")); err == nil {
			t.Fatal("token with tampered signature must be rejected")
		}
	})
}
