package storage_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/warlck/food-flow/foundation/storage"
)

func newTestSigner(t *testing.T) storage.Signer {
	t.Helper()

	signer, err := storage.NewSigner(context.Background(), storage.Config{
		Backend:      storage.BackendLocal,
		LocalDir:     t.TempDir(),
		LocalBaseURL: "/v1/images/local",
		URLTTL:       15 * time.Minute,
	})
	if err != nil {
		t.Fatalf("creating signer: %v", err)
	}
	return signer
}

func TestLocalSignUpload(t *testing.T) {
	signer := newTestSigner(t)

	got, err := signer.SignUpload(context.Background(), storage.SignRequest{
		ObjectPath:  "restaurants/abc/cover.jpg",
		ContentType: "image/jpeg",
		SizeBytes:   1024,
	})
	if err != nil {
		t.Fatalf("signing upload: %v", err)
	}

	if got.Method != http.MethodPut {
		t.Errorf("method = %q, want PUT", got.Method)
	}
	if got.Headers["Content-Type"] != "image/jpeg" {
		t.Errorf("content-type header = %q, want image/jpeg", got.Headers["Content-Type"])
	}
	wantURL := "/v1/images/local/" + "restaurants%2Fabc%2Fcover.jpg"
	if got.UploadURL != wantURL {
		t.Errorf("upload url = %q, want %q", got.UploadURL, wantURL)
	}
	if got.PublicURL != got.UploadURL {
		t.Errorf("public url = %q, want same-origin %q", got.PublicURL, got.UploadURL)
	}
	if got.ExpiresAt.Before(time.Now()) {
		t.Error("expires at should be in the future")
	}
}

func TestLocalSignUploadRejectsTraversal(t *testing.T) {
	signer := newTestSigner(t)

	bad := []string{"", "..", "../escape.jpg", "a/../../escape.jpg", "a/../b/.."}
	for _, path := range bad {
		if _, err := signer.SignUpload(context.Background(), storage.SignRequest{
			ObjectPath:  path,
			ContentType: "image/png",
			SizeBytes:   1,
		}); err == nil {
			t.Errorf("path %q: expected error, got nil", path)
		}
	}
}

func TestLocalStatAndDelete(t *testing.T) {
	signer := newTestSigner(t)
	ctx := context.Background()

	_, err := signer.Stat(ctx, "menu-items/x/y.jpg")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("stat missing object: got %v, want ErrNotFound", err)
	}

	// Deleting a missing object must not fail.
	if err := signer.Delete(ctx, "menu-items/x/y.jpg"); err != nil {
		t.Fatalf("delete missing object: %v", err)
	}

	// Simulate the dev PUT handler writing the object.
	dir := t.TempDir()
	signer, err = storage.NewSigner(ctx, storage.Config{
		Backend:      storage.BackendLocal,
		LocalDir:     dir,
		LocalBaseURL: "/v1/images/local",
		URLTTL:       time.Minute,
	})
	if err != nil {
		t.Fatalf("creating signer: %v", err)
	}

	objectPath := "restaurants/r1/cover.png"
	fullPath := filepath.Join(dir, objectPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte("png-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	info, err := signer.Stat(ctx, objectPath)
	if err != nil {
		t.Fatalf("stat existing object: %v", err)
	}
	if info.SizeBytes != int64(len("png-bytes")) {
		t.Errorf("size = %d, want %d", info.SizeBytes, len("png-bytes"))
	}

	if err := signer.Delete(ctx, objectPath); err != nil {
		t.Fatalf("delete existing object: %v", err)
	}
	if _, err := signer.Stat(ctx, objectPath); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("stat after delete: got %v, want ErrNotFound", err)
	}
}

func TestNewSignerUnknownBackend(t *testing.T) {
	_, err := storage.NewSigner(context.Background(), storage.Config{Backend: "s3"})
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}
