package storage

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const testServiceAccount = "staging-workload@test-project.iam.gserviceaccount.com"

// newTestGCSSigner builds a signer backed by an unauthenticated storage
// client, optionally pointed at a fake JSON API endpoint. SignedURL performs
// no network calls when SignBytes is provided, and Stat/Delete hit the fake
// server, so no GCP credentials are required.
func newTestGCSSigner(t *testing.T, endpoint string, signBlob func(context.Context, []byte) ([]byte, error)) *gcsSigner {
	t.Helper()

	opts := []option.ClientOption{option.WithoutAuthentication()}
	if endpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoint))
	}

	client, err := gcs.NewClient(context.Background(), opts...)
	if err != nil {
		t.Fatalf("creating storage client: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	return &gcsSigner{
		objects:        client,
		signBlob:       signBlob,
		serviceAccount: testServiceAccount,
		bucket:         "test-bucket",
		publicBaseURL:  "https://storage.googleapis.com/test-bucket",
		ttl:            time.Minute,
	}
}

func TestNormalizePublicBaseURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"no trailing slash", "https://storage.googleapis.com/test-bucket", "https://storage.googleapis.com/test-bucket"},
		{"trailing slash", "https://storage.googleapis.com/test-bucket/", "https://storage.googleapis.com/test-bucket"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizePublicBaseURL(tc.in); got != tc.want {
				t.Errorf("normalizePublicBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNewGCSSignerValidation(t *testing.T) {
	ctx := context.Background()
	base := Config{
		Backend:        BackendGCS,
		Bucket:         "test-bucket",
		ServiceAccount: testServiceAccount,
		PublicBaseURL:  "https://storage.googleapis.com/test-bucket",
		URLTTL:         time.Minute,
	}

	cases := []struct {
		name   string
		mutate func(*Config)
	}{
		{"missing bucket", func(c *Config) { c.Bucket = "" }},
		{"missing public base url", func(c *Config) { c.PublicBaseURL = "" }},
		{"zero url ttl", func(c *Config) { c.URLTTL = 0 }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := base
			tc.mutate(&cfg)
			if _, err := newGCSSigner(ctx, cfg); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestGCSSignUpload(t *testing.T) {
	ctx := context.Background()
	fakeSig := []byte("fake-signature")

	var gotPayload []byte
	signer := newTestGCSSigner(t, "", func(_ context.Context, payload []byte) ([]byte, error) {
		gotPayload = append([]byte(nil), payload...)
		return fakeSig, nil
	})

	got, err := signer.SignUpload(ctx, SignRequest{
		ObjectPath:  "restaurants/r1/cover.png",
		ContentType: "image/png",
		SizeBytes:   1024,
	})
	if err != nil {
		t.Fatalf("signing upload: %v", err)
	}

	if got.Method != http.MethodPut {
		t.Errorf("method = %q, want PUT", got.Method)
	}
	if got.Headers["Content-Type"] != "image/png" {
		t.Errorf("content-type header = %q, want image/png", got.Headers["Content-Type"])
	}
	if want := "https://storage.googleapis.com/test-bucket/restaurants/r1/cover.png"; got.PublicURL != want {
		t.Errorf("public url = %q, want %q", got.PublicURL, want)
	}
	if got.ObjectPath != "restaurants/r1/cover.png" {
		t.Errorf("object path = %q, want restaurants/r1/cover.png", got.ObjectPath)
	}
	if got.ExpiresAt.Before(time.Now().Add(30 * time.Second)) {
		t.Errorf("expires at %v should be roughly a minute out", got.ExpiresAt)
	}

	u, err := url.Parse(got.UploadURL)
	if err != nil {
		t.Fatalf("parsing upload url: %v", err)
	}
	if u.Scheme != "https" || u.Host != "storage.googleapis.com" {
		t.Errorf("url = %q, want https://storage.googleapis.com/...", got.UploadURL)
	}
	if u.Path != "/test-bucket/restaurants/r1/cover.png" {
		t.Errorf("url path = %q, want /test-bucket/restaurants/r1/cover.png", u.Path)
	}

	q := u.Query()
	if v := q.Get("X-Goog-Algorithm"); v != "GOOG4-RSA-SHA256" {
		t.Errorf("algorithm = %q, want GOOG4-RSA-SHA256", v)
	}
	if v := q.Get("X-Goog-Credential"); !strings.HasPrefix(v, testServiceAccount+"/") {
		t.Errorf("credential = %q, want prefix %q", v, testServiceAccount+"/")
	}
	if v := q.Get("X-Goog-Expires"); v != "59" && v != "60" {
		t.Errorf("expires = %q, want the 60s ttl (truncated to 59 or 60)", v)
	}
	if v := q.Get("X-Goog-SignedHeaders"); !strings.Contains(v, "content-type") {
		t.Errorf("signed headers = %q, want content-type bound into the signature", v)
	}
	if v, want := q.Get("X-Goog-Signature"), hex.EncodeToString(fakeSig); v != want {
		t.Errorf("signature = %q, want hex of fake signBlob output %q", v, want)
	}

	// The payload handed to signBlob must be the V4 string-to-sign.
	if !strings.HasPrefix(string(gotPayload), "GOOG4-RSA-SHA256") {
		t.Errorf("string-to-sign should start with GOOG4-RSA-SHA256, got %q", gotPayload)
	}
}

func TestGCSSignUploadTrimsLeadingSlash(t *testing.T) {
	ctx := context.Background()
	signer := newTestGCSSigner(t, "", func(_ context.Context, payload []byte) ([]byte, error) {
		return []byte("sig"), nil
	})

	got, err := signer.SignUpload(ctx, SignRequest{
		ObjectPath:  "/restaurants/r1/cover.png",
		ContentType: "image/png",
		SizeBytes:   1,
	})
	if err != nil {
		t.Fatalf("signing upload: %v", err)
	}

	if got.ObjectPath != "restaurants/r1/cover.png" {
		t.Errorf("object path = %q, want leading slash trimmed", got.ObjectPath)
	}
	if want := "https://storage.googleapis.com/test-bucket/restaurants/r1/cover.png"; got.PublicURL != want {
		t.Errorf("public url = %q, want %q", got.PublicURL, want)
	}

	u, err := url.Parse(got.UploadURL)
	if err != nil {
		t.Fatalf("parsing upload url: %v", err)
	}
	if u.Path != "/test-bucket/restaurants/r1/cover.png" {
		t.Errorf("signed url path = %q, want the trimmed object path", u.Path)
	}
}

func TestGCSSignUploadEscapesPublicURL(t *testing.T) {
	ctx := context.Background()
	signer := newTestGCSSigner(t, "", func(_ context.Context, payload []byte) ([]byte, error) {
		return []byte("sig"), nil
	})

	got, err := signer.SignUpload(ctx, SignRequest{
		ObjectPath:  "restaurants/r1/my cover#1.png",
		ContentType: "image/png",
		SizeBytes:   1,
	})
	if err != nil {
		t.Fatalf("signing upload: %v", err)
	}

	// The object name itself stays raw; only the URL representation escapes.
	if got.ObjectPath != "restaurants/r1/my cover#1.png" {
		t.Errorf("object path = %q, want the raw object name", got.ObjectPath)
	}
	if want := "https://storage.googleapis.com/test-bucket/restaurants/r1/my%20cover%231.png"; got.PublicURL != want {
		t.Errorf("public url = %q, want %q", got.PublicURL, want)
	}
	if !strings.Contains(got.UploadURL, "my%20cover%231.png") {
		t.Errorf("signed url = %q, want the escaped object path", got.UploadURL)
	}
}

func TestGCSSignUploadSignError(t *testing.T) {
	signer := newTestGCSSigner(t, "", func(context.Context, []byte) ([]byte, error) {
		return nil, errors.New("iam unavailable")
	})

	_, err := signer.SignUpload(context.Background(), SignRequest{
		ObjectPath:  "restaurants/r1/cover.png",
		ContentType: "image/png",
		SizeBytes:   1,
	})
	if err == nil || !strings.Contains(err.Error(), "signing upload url") {
		t.Fatalf("got %v, want wrapped signing error", err)
	}
}

// fakeGCSAPI emulates the storage JSON API endpoints Stat and Delete use.
func fakeGCSAPI(t *testing.T) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.EscapedPath(), "missing"):
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":{"code":404,"message":"not found","errors":[{"message":"not found","domain":"global","reason":"notFound"}]}}`)
		case strings.Contains(r.URL.EscapedPath(), "broken"):
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":{"code":500,"message":"boom"}}`)
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			fmt.Fprint(w, `{"kind":"storage#object","bucket":"test-bucket","name":"restaurants/r1/cover.png","size":"123","contentType":"image/png","updated":"2026-08-18T00:00:00.000Z"}`)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestGCSStat(t *testing.T) {
	ctx := context.Background()
	signer := newTestGCSSigner(t, fakeGCSAPI(t).URL+"/storage/v1/", nil)

	info, err := signer.Stat(ctx, "restaurants/r1/cover.png")
	if err != nil {
		t.Fatalf("stat existing object: %v", err)
	}
	if info.ObjectPath != "restaurants/r1/cover.png" {
		t.Errorf("object path = %q, want restaurants/r1/cover.png", info.ObjectPath)
	}
	if info.SizeBytes != 123 {
		t.Errorf("size = %d, want 123", info.SizeBytes)
	}
	if info.ContentType != "image/png" {
		t.Errorf("content type = %q, want image/png", info.ContentType)
	}

	if _, err := signer.Stat(ctx, "missing.png"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("stat missing object: got %v, want ErrNotFound", err)
	}

	// The SDK retries 5xx with backoff until the context ends, so bound it.
	errCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if _, err := signer.Stat(errCtx, "broken.png"); err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("stat server error: got %v, want non-notfound error", err)
	}
}

func TestGCSDelete(t *testing.T) {
	ctx := context.Background()
	signer := newTestGCSSigner(t, fakeGCSAPI(t).URL+"/storage/v1/", nil)

	if err := signer.Delete(ctx, "restaurants/r1/cover.png"); err != nil {
		t.Fatalf("delete existing object: %v", err)
	}

	// Deleting a missing object must not fail.
	if err := signer.Delete(ctx, "missing.png"); err != nil {
		t.Fatalf("delete missing object: got %v, want nil", err)
	}

	// The SDK retries 5xx with backoff until the context ends, so bound it.
	errCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := signer.Delete(errCtx, "broken.png"); err == nil {
		t.Fatal("delete server error: got nil, want error")
	}
}
