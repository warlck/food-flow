// Package authclient provides support to access the auth service.
package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/foundation/logger"
)

// This provides a default client configuration, but it's recommended
// this is replaced by the user with application specific settings using
// the WithClient function at the time a AuthAPI is constructed.
// DualStack Deprecated: Fast Fallback is enabled by default. To disable, set FallbackDelay to a negative value.
var defaultClient = http.Client{
	// Overall request timeout: the sales service calls the auth service on
	// every authenticated request, so a hung auth service must fail fast
	// (fail-closed) instead of tying up sales handlers indefinitely.
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

// Client represents a client that can talk to the auth service.
type Client struct {
	log  *logger.Logger
	url  string
	http *http.Client
}

// New constructs an Auth that can be used to talk with the auth service.
func New(log *logger.Logger, url string, options ...func(cln *Client)) *Client {
	// Copy the default client so per-client options never mutate the shared
	// package-level value.
	hc := defaultClient

	cln := Client{
		log:  log,
		url:  url,
		http: &hc,
	}

	for _, option := range options {
		option(&cln)
	}

	return &cln
}

// WithClient adds a custom client for processing requests. It's recommend
// to not use the default client and provide your own.
func WithClient(hc *http.Client) func(cln *Client) {
	return func(cln *Client) {
		cln.http = hc
	}
}

// WithTimeout overrides the overall request timeout of the default client.
func WithTimeout(d time.Duration) func(cln *Client) {
	return func(cln *Client) {
		cln.http.Timeout = d
	}
}

// Authenticate calls the auth service to authenticate the user.
func (cln *Client) Authenticate(ctx context.Context, authorization string) (AuthenticateResp, error) {
	endpoint := fmt.Sprintf("%s/v1/auth/authenticate", cln.url)

	headers := map[string]string{
		"authorization": authorization,
	}

	var resp AuthenticateResp
	if err := cln.do(ctx, http.MethodGet, endpoint, headers, nil, &resp); err != nil {
		return AuthenticateResp{}, err
	}

	return resp, nil
}

// Authorize calls the auth service to authorize the user.
func (cln *Client) Authorize(ctx context.Context, auth Authorize) error {
	endpoint := fmt.Sprintf("%s/v1/auth/authorize", cln.url)

	if err := cln.do(ctx, http.MethodPost, endpoint, nil, auth, nil); err != nil {
		return err
	}

	return nil
}

func (cln *Client) do(ctx context.Context, method string, endpoint string, headers map[string]string, body any, v any) error {
	var statusCode int

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parsing endpoint: %w", err)
	}
	base := path.Base(u.Path)

	cln.log.Info(ctx, "authclient: rawRequest: started", "method", method, "call", base, "endpoint", endpoint)
	defer func() {
		cln.log.Info(ctx, "authclient: rawRequest: completed", "status", statusCode)
	}()

	// ctx, span := otel.AddSpan(ctx, fmt.Sprintf("app.sdk.authclient.%s", base), attribute.String("endpoint", endpoint))
	// defer func() {
	// 	span.SetAttributes(attribute.Int("status", statusCode))
	// 	span.End()
	// }()

	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return fmt.Errorf("encoding error: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, &b)
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for key, value := range headers {
		// Never log credential material: the authorization header carries a
		// bearer token, and tokens in logs are credential leakage.
		logValue := value
		if strings.EqualFold(key, "authorization") {
			logValue = "[REDACTED]"
		}
		cln.log.Info(ctx, "authclient: rawRequest", "key", key, "value", logValue)
		req.Header.Set(key, value)
	}

	// otel.AddTraceToRequest(ctx, req)

	resp, err := cln.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: error: %w", err)
	}
	defer resp.Body.Close()

	// Assign so it can be logged in the defer above.
	statusCode = resp.StatusCode

	if statusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	switch statusCode {
	case http.StatusOK:
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), err)
		}
		return nil

	case http.StatusUnauthorized, http.StatusForbidden:
		var err *errs.Error
		if err := json.Unmarshal(data, &err); err != nil {
			return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), err)
		}
		return err

	default:
		return fmt.Errorf("failed: response: %s", string(data))
	}
}
