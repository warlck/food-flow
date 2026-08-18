package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// gcsSigner signs V4 upload URLs for Google Cloud Storage. Cloud Run and GCE
// workloads have no service-account private key available, so signing is
// delegated to the IAM Credentials signBlob API using the ambient identity.
type gcsSigner struct {
	objects        *gcs.Client
	signing        *credentials.IamCredentialsClient
	signBlob       func(ctx context.Context, payload []byte) ([]byte, error)
	serviceAccount string
	bucket         string
	publicBaseURL  string
	ttl            time.Duration
}

func newGCSSigner(ctx context.Context, cfg Config) (*gcsSigner, error) {
	if cfg.Bucket == "" {
		return nil, errors.New("gcs backend: bucket is required")
	}
	if cfg.PublicBaseURL == "" {
		return nil, errors.New("gcs backend: public base url is required")
	}
	if cfg.URLTTL <= 0 {
		return nil, errors.New("gcs backend: url ttl must be greater than zero")
	}

	serviceAccount := cfg.ServiceAccount
	if serviceAccount == "" {
		email, err := metadata.EmailWithContext(ctx, "default")
		if err != nil {
			return nil, fmt.Errorf("gcs backend: resolving service account from metadata server: %w", err)
		}
		serviceAccount = email
	}

	objects, err := gcs.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs backend: creating storage client: %w", err)
	}

	signing, err := credentials.NewIamCredentialsClient(ctx, option.WithEndpoint("iamcredentials.googleapis.com"))
	if err != nil {
		objects.Close()
		return nil, fmt.Errorf("gcs backend: creating iam credentials client: %w", err)
	}

	signer := &gcsSigner{
		objects:        objects,
		signing:        signing,
		serviceAccount: serviceAccount,
		bucket:         cfg.Bucket,
		publicBaseURL:  normalizePublicBaseURL(cfg.PublicBaseURL),
		ttl:            cfg.URLTTL,
	}
	signer.signBlob = func(ctx context.Context, payload []byte) ([]byte, error) {
		resp, err := signing.SignBlob(ctx, &credentialspb.SignBlobRequest{
			Name:    "projects/-/serviceAccounts/" + serviceAccount,
			Payload: payload,
		})
		if err != nil {
			return nil, err
		}
		return resp.SignedBlob, nil
	}

	return signer, nil
}

// normalizePublicBaseURL strips a trailing slash so joining the base URL with
// an object path never produces a double slash.
func normalizePublicBaseURL(baseURL string) string {
	return strings.TrimSuffix(baseURL, "/")
}

// SignUpload mints a V4 signed URL for a PUT of the exact content type given.
// Leading slashes are trimmed from the object path before signing so the
// signed object and the returned public URL always refer to the same object.
func (s *gcsSigner) SignUpload(ctx context.Context, req SignRequest) (SignedUpload, error) {
	objectPath := strings.TrimLeft(req.ObjectPath, "/")
	if err := validateObjectPath(objectPath); err != nil {
		return SignedUpload{}, err
	}
	expires := time.Now().Add(s.ttl)

	signedURL, err := s.objects.Bucket(s.bucket).SignedURL(objectPath, &gcs.SignedURLOptions{
		Scheme:         gcs.SigningSchemeV4,
		Method:         http.MethodPut,
		Expires:        expires,
		ContentType:    req.ContentType,
		GoogleAccessID: s.serviceAccount,
		SignBytes: func(b []byte) ([]byte, error) {
			return s.signBlob(ctx, b)
		},
	})
	if err != nil {
		return SignedUpload{}, fmt.Errorf("signing upload url: %w", err)
	}

	return SignedUpload{
		UploadURL: signedURL,
		Method:    http.MethodPut,
		Headers: map[string]string{
			"Content-Type": req.ContentType,
		},
		PublicURL:  s.publicBaseURL + "/" + escapeObjectPath(objectPath),
		ObjectPath: objectPath,
		ExpiresAt:  expires,
	}, nil
}

// validateObjectPath rejects empty object paths before any remote call, so
// misconfiguration surfaces as a clean local error instead of a cryptic GCS
// API failure.
func validateObjectPath(objectPath string) error {
	if objectPath == "" {
		return errors.New("object path is empty")
	}
	return nil
}

// escapeObjectPath percent-escapes each path segment while preserving the
// slash separators, so object names with spaces or special characters still
// produce a valid public URL.
func escapeObjectPath(objectPath string) string {
	segments := strings.Split(objectPath, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}

// Stat returns object metadata, mapping a missing object to ErrNotFound.
func (s *gcsSigner) Stat(ctx context.Context, objectPath string) (ObjectInfo, error) {
	if err := validateObjectPath(objectPath); err != nil {
		return ObjectInfo{}, err
	}

	attrs, err := s.objects.Bucket(s.bucket).Object(objectPath).Attrs(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return ObjectInfo{}, ErrNotFound
		}
		return ObjectInfo{}, fmt.Errorf("stat object: %w", err)
	}

	return ObjectInfo{
		ObjectPath:  objectPath,
		SizeBytes:   attrs.Size,
		ContentType: attrs.ContentType,
	}, nil
}

// Delete removes the object. A missing object is treated as already deleted.
func (s *gcsSigner) Delete(ctx context.Context, objectPath string) error {
	if err := validateObjectPath(objectPath); err != nil {
		return err
	}

	err := s.objects.Bucket(s.bucket).Object(objectPath).Delete(ctx)
	if err != nil && !errors.Is(err, gcs.ErrObjectNotExist) {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

// Close releases the underlying clients.
func (s *gcsSigner) Close() error {
	signErr := s.signing.Close()
	objErr := s.objects.Close()
	return errors.Join(signErr, objErr)
}
