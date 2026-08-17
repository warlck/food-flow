// Package storage provides an abstraction over object storage for direct
// client uploads via signed URLs. The GCS implementation is used in deployed
// environments while the local implementation keeps development self-contained.
package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned when an object does not exist in storage.
var ErrNotFound = errors.New("object not found")

// Supported backend names for Config.Backend.
const (
	BackendGCS   = "gcs"
	BackendLocal = "local"
)

// Config collects the settings required to construct a Signer.
type Config struct {
	// Backend selects the implementation: "gcs" or "local".
	Backend string

	// Bucket is the GCS bucket name (gcs backend only).
	Bucket string

	// ServiceAccount is the email used as GoogleAccessID for signing. When
	// empty on GCE/Cloud Run it is resolved from the metadata server.
	ServiceAccount string

	// PublicBaseURL prefixes object paths to form public URLs, e.g.
	// "https://storage.googleapis.com/my-bucket" (gcs backend only).
	PublicBaseURL string

	// URLTTL is how long a signed upload URL remains valid.
	URLTTL time.Duration

	// LocalDir is the filesystem directory for stored objects (local backend).
	LocalDir string

	// LocalBaseURL is the same-origin URL prefix the API serves local uploads
	// under, e.g. "/v1/images/local" (local backend only).
	LocalBaseURL string
}

// SignRequest describes an upload to authorize.
type SignRequest struct {
	ObjectPath  string
	ContentType string
	SizeBytes   int64
}

// SignedUpload contains everything a client needs to upload directly to
// storage without proxying bytes through the API.
type SignedUpload struct {
	UploadURL  string
	Method     string
	Headers    map[string]string
	PublicURL  string
	ObjectPath string
	ExpiresAt  time.Time
}

// ObjectInfo describes an object stored in the backend.
type ObjectInfo struct {
	ObjectPath  string
	SizeBytes   int64
	ContentType string
}

// Signer abstracts object storage so business code stays backend agnostic.
type Signer interface {
	// SignUpload returns a pre-authorized URL the client PUTs the file to.
	SignUpload(ctx context.Context, req SignRequest) (SignedUpload, error)

	// Stat returns metadata for an object or ErrNotFound.
	Stat(ctx context.Context, objectPath string) (ObjectInfo, error)

	// Delete removes an object. Deleting a missing object is not an error.
	Delete(ctx context.Context, objectPath string) error
}

// NewSigner constructs a Signer based on the configured backend.
func NewSigner(ctx context.Context, cfg Config) (Signer, error) {
	switch cfg.Backend {
	case BackendGCS:
		return newGCSSigner(ctx, cfg)
	case BackendLocal:
		return newLocalSigner(cfg)
	default:
		return nil, fmt.Errorf("unknown storage backend %q: expected %q or %q", cfg.Backend, BackendGCS, BackendLocal)
	}
}
