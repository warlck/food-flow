package keystore_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"testing"
	"testing/fstest"

	"github.com/warlck/food-flow/foundation/keystore"
)

func rsaPKCS1PEM(t *testing.T) (string, *rsa.PublicKey) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %s", err)
	}

	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}

	return string(pem.EncodeToMemory(block)), &key.PublicKey
}

func rsaPKCS8PEM(t *testing.T) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %s", err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshaling PKCS8 key: %s", err)
	}

	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

func ecPEM(t *testing.T) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating EC key: %s", err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshaling EC key: %s", err)
	}

	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

// secretPayload builds the {"key","pem"} document exactly the way
// bootstrap-auth-secret.sh and make dev-auth-keys build it with jq:
//
//	jq -n --arg key "$KID" --rawfile pem key.pem '{key: $key, pem: $pem}'
//
// This guards the Secret Manager / k8s secret contract.
func secretPayload(t *testing.T, kid string, privatePEM string) string {
	t.Helper()

	doc, err := json.Marshal(map[string]string{"key": kid, "pem": privatePEM})
	if err != nil {
		t.Fatalf("marshaling payload: %s", err)
	}

	return string(doc)
}

func publicKeyMatches(t *testing.T, publicPEM string, want *rsa.PublicKey) {
	t.Helper()

	block, _ := pem.Decode([]byte(publicPEM))
	if block == nil {
		t.Fatal("public PEM does not decode")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Fatalf("parsing public key: %s", err)
	}

	got, ok := key.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("public key type = %T, want *rsa.PublicKey", key)
	}

	if got.N.Cmp(want.N) != 0 || got.E != want.E {
		t.Error("derived public key does not match the private key's public part")
	}
}

func Test_LoadByJSON(t *testing.T) {
	t.Parallel()

	rsaPEM, pub := rsaPKCS1PEM(t)

	t.Run("secret manager payload", func(t *testing.T) {
		ks := keystore.New()

		n, err := ks.LoadByJSON(secretPayload(t, "prod-kid", rsaPEM))
		if err != nil {
			t.Fatalf("load: %s", err)
		}
		if n != 1 {
			t.Errorf("key count = %d, want 1", n)
		}

		gotPrivate, err := ks.PrivateKey("prod-kid")
		if err != nil {
			t.Fatalf("PrivateKey: %s", err)
		}
		if gotPrivate != rsaPEM {
			t.Error("PrivateKey does not return the loaded PEM")
		}

		gotPublic, err := ks.PublicKey("prod-kid")
		if err != nil {
			t.Fatalf("PublicKey: %s", err)
		}
		publicKeyMatches(t, gotPublic, pub)
	})

	t.Run("pkcs8 private key accepted", func(t *testing.T) {
		ks := keystore.New()

		if _, err := ks.LoadByJSON(secretPayload(t, "pkcs8-kid", rsaPKCS8PEM(t))); err != nil {
			t.Fatalf("PKCS8 key should load: %s", err)
		}
	})

	t.Run("empty document is a no-op", func(t *testing.T) {
		ks := keystore.New()

		n, err := ks.LoadByJSON("")
		if err != nil {
			t.Fatalf("empty document should not error: %s", err)
		}
		if n != 0 {
			t.Errorf("key count = %d, want 0", n)
		}
	})

	t.Run("malformed JSON rejected", func(t *testing.T) {
		ks := keystore.New()

		if _, err := ks.LoadByJSON("{not json"); err == nil {
			t.Fatal("malformed JSON must be rejected")
		}
	})

	t.Run("non-RSA key rejected", func(t *testing.T) {
		ks := keystore.New()

		if _, err := ks.LoadByJSON(secretPayload(t, "ec-kid", ecPEM(t))); err == nil {
			t.Fatal("EC key must be rejected, only RSA is supported")
		}
	})

	t.Run("non-PEM payload rejected", func(t *testing.T) {
		ks := keystore.New()

		if _, err := ks.LoadByJSON(secretPayload(t, "garbage-kid", "this is not a pem")); err == nil {
			t.Fatal("non-PEM payload must be rejected")
		}
	})

	t.Run("multiple kids accumulate and reload overwrites", func(t *testing.T) {
		ks := keystore.New()

		if _, err := ks.LoadByJSON(secretPayload(t, "kid-1", rsaPEM)); err != nil {
			t.Fatalf("load kid-1: %s", err)
		}

		otherPEM, _ := rsaPKCS1PEM(t)
		n, err := ks.LoadByJSON(secretPayload(t, "kid-2", otherPEM))
		if err != nil {
			t.Fatalf("load kid-2: %s", err)
		}
		if n != 2 {
			t.Errorf("key count = %d, want 2", n)
		}

		n, err = ks.LoadByJSON(secretPayload(t, "kid-1", otherPEM))
		if err != nil {
			t.Fatalf("reload kid-1: %s", err)
		}
		if n != 2 {
			t.Errorf("key count after reload = %d, want 2", n)
		}

		got, err := ks.PrivateKey("kid-1")
		if err != nil {
			t.Fatalf("PrivateKey: %s", err)
		}
		if got != otherPEM {
			t.Error("reloading a kid should replace the stored key")
		}
	})
}

func Test_LoadByFileSystem(t *testing.T) {
	t.Parallel()

	rsaPEM, pub := rsaPKCS1PEM(t)

	t.Run("loads pem files by name and ignores others", func(t *testing.T) {
		otherPEM, _ := rsaPKCS1PEM(t)

		fsys := fstest.MapFS{
			"local-dev.pem": &fstest.MapFile{Data: []byte(rsaPEM)},
			"second.pem":    &fstest.MapFile{Data: []byte(otherPEM)},
			"notes.txt":     &fstest.MapFile{Data: []byte("not a key")},
			".gitkeep":      &fstest.MapFile{Data: []byte("")},
		}

		ks := keystore.New()

		n, err := ks.LoadByFileSystem(fsys)
		if err != nil {
			t.Fatalf("load: %s", err)
		}
		if n != 2 {
			t.Errorf("key count = %d, want 2 (non-pem files must be ignored)", n)
		}

		if _, err := ks.PrivateKey("local-dev"); err != nil {
			t.Errorf("kid local-dev should exist: %s", err)
		}
		if _, err := ks.PrivateKey("second"); err != nil {
			t.Errorf("kid second should exist: %s", err)
		}

		gotPublic, err := ks.PublicKey("local-dev")
		if err != nil {
			t.Fatalf("PublicKey: %s", err)
		}
		publicKeyMatches(t, gotPublic, pub)
	})

	t.Run("missing directory errors", func(t *testing.T) {
		ks := keystore.New()

		if _, err := ks.LoadByFileSystem(os.DirFS("does-not-exist")); err == nil {
			t.Fatal("a missing directory must error so startup fails fast")
		}
	})

	t.Run("invalid pem file errors", func(t *testing.T) {
		fsys := fstest.MapFS{
			"broken.pem": &fstest.MapFile{Data: []byte("not a pem")},
		}

		ks := keystore.New()

		if _, err := ks.LoadByFileSystem(fsys); err == nil {
			t.Fatal("an unreadable key file must error so startup fails fast")
		}
	})
}

func Test_UnknownKID(t *testing.T) {
	t.Parallel()

	ks := keystore.New()

	if _, err := ks.PrivateKey("no-such-kid"); err == nil {
		t.Error("PrivateKey for an unknown kid must error")
	}

	if _, err := ks.PublicKey("no-such-kid"); err == nil {
		t.Error("PublicKey for an unknown kid must error")
	}
}
