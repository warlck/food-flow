package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
		publicBaseURL:  cfg.PublicBaseURL,
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

// SignUpload mints a V4 signed URL for a PUT of the exact content type given.
func (s *gcsSigner) SignUpload(ctx context.Context, req SignRequest) (SignedUpload, error) {
	expires := time.Now().Add(s.ttl)

	url, err := s.objects.Bucket(s.bucket).SignedURL(req.ObjectPath, &gcs.SignedURLOptions{
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
		UploadURL: url,
		Method:    http.MethodPut,
		Headers: map[string]string{
			"Content-Type": req.ContentType,
		},
		PublicURL:  s.publicBaseURL + "/" + req.ObjectPath,
		ObjectPath: req.ObjectPath,
		ExpiresAt:  expires,
	}, nil
}

// Stat returns object metadata, mapping a missing object to ErrNotFound.
func (s *gcsSigner) Stat(ctx context.Context, objectPath string) (ObjectInfo, error) {
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
