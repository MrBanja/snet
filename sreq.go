package sreq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// NewRequest creates a new http.Request with the given method, url, and body.
// Body is an optional json-serializable struct
func NewRequest(ctx context.Context, method string, url string, body any) (*http.Request, error) {
	var buf = new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encode body for [%s %s]: %w", method, url, err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("create request for [%s %s]: %w", method, url, err)
	}
	return req, nil
}

// Unmarshal unmarshals the response body into a new instance of T
func Unmarshal[T any](r *http.Response) (*T, error) {
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("unmarshal response for [%s %s]: %w", r.Request.Method, r.Request.URL.String(), err)
	}
	return &payload, r.Body.Close()
}

// U creates a new url.URL with the given path appended to the base url
// Example: U("https://test.localhost/api", "/user/create") -> "https://test.localhost/api/user/create"
func U(base, path string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", base, strings.TrimPrefix(path, "/")))
	if err != nil {
		return nil, fmt.Errorf("parse url [%s]: %w", path, err)
	}
	return u, nil
}
