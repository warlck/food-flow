package storage

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// localSigner implements Signer against the local filesystem for development.
// Upload URLs point back at the API itself, which serves a same-origin PUT
// handler, so no external storage dependency is required.
type localSigner struct {
	dir     string
	baseURL string
	ttl     time.Duration
}

func newLocalSigner(cfg Config) (*localSigner, error) {
	if cfg.LocalDir == "" {
		return nil, errors.New("local backend: local dir is required")
	}
	if cfg.LocalBaseURL == "" {
		return nil, errors.New("local backend: local base url is required")
	}
	if cfg.URLTTL <= 0 {
		return nil, errors.New("local backend: url ttl must be greater than zero")
	}

	if err := os.MkdirAll(cfg.LocalDir, 0o755); err != nil {
		return nil, fmt.Errorf("local backend: creating local dir: %w", err)
	}

	return &localSigner{
		dir:     cfg.LocalDir,
		baseURL: strings.TrimSuffix(cfg.LocalBaseURL, "/"),
		ttl:     cfg.URLTTL,
	}, nil
}

// SignUpload returns a same-origin URL the API's local upload handler serves.
func (s *localSigner) SignUpload(_ context.Context, req SignRequest) (SignedUpload, error) {
	if _, err := s.LocalPath(req.ObjectPath); err != nil {
		return SignedUpload{}, err
	}

	escaped := "/" + url.PathEscape(req.ObjectPath)

	return SignedUpload{
		UploadURL: s.baseURL + escaped,
		Method:    http.MethodPut,
		Headers: map[string]string{
			"Content-Type": req.ContentType,
		},
		PublicURL:  s.baseURL + escaped,
		ObjectPath: req.ObjectPath,
		ExpiresAt:  time.Now().Add(s.ttl),
	}, nil
}

// Stat returns file metadata, mapping a missing file to ErrNotFound.
func (s *localSigner) Stat(_ context.Context, objectPath string) (ObjectInfo, error) {
	fullPath, err := s.LocalPath(objectPath)
	if err != nil {
		return ObjectInfo{}, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ObjectInfo{}, ErrNotFound
		}
		return ObjectInfo{}, fmt.Errorf("stat object: %w", err)
	}

	return ObjectInfo{
		ObjectPath:  objectPath,
		SizeBytes:   info.Size(),
		ContentType: "",
	}, nil
}

// Delete removes the file. A missing file is treated as already deleted.
func (s *localSigner) Delete(_ context.Context, objectPath string) error {
	fullPath, err := s.LocalPath(objectPath)
	if err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

// LocalPath maps an object path to a filesystem path inside the configured
// directory. Paths that would escape the directory are rejected.
func (s *localSigner) LocalPath(objectPath string) (string, error) {
	if objectPath == "" {
		return "", errors.New("object path is empty")
	}

	// Reject traversal attempts outright rather than silently confining them.
	for _, seg := range strings.Split(objectPath, "/") {
		if seg == ".." {
			return "", fmt.Errorf("object path %q contains invalid segment", objectPath)
		}
	}

	// Rooting the path before cleaning keeps the result inside the directory.
	clean := filepath.Clean("/" + objectPath)
	fullPath := filepath.Join(s.dir, clean)

	if !strings.HasPrefix(fullPath, filepath.Clean(s.dir)+string(os.PathSeparator)) {
		return "", fmt.Errorf("object path %q escapes the storage directory", objectPath)
	}

	return fullPath, nil
}
