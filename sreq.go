package sreq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// UnmarshalResp unmarshalls the response body into a new instance of T
func UnmarshalResp[T any](r *http.Response) (*T, error) {
	defer func() { _ = r.Body.Close() }()
	payload, err := Unmarshal[T](r.Body)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response for [%s %s]: %w", r.Request.Method, r.Request.URL.String(), err)
	}
	return payload, nil
}

// UnmarshalReq unmarshalls the request body into a new instance of T
func UnmarshalReq[T any](r *http.Request) (*T, error) {
	defer func() { _ = r.Body.Close() }()
	payload, err := Unmarshal[T](r.Body)
	if err != nil {
		return nil, fmt.Errorf("unmarshal request for [%s %s]: %w", r.Method, r.URL.String(), err)
	}
	return payload, nil
}

// Unmarshal unmarshalls io.Reader into a new instance of T
func Unmarshal[T any](r io.Reader) (*T, error) {
	var payload T
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// U creates a new url.URL with the given path appended to the base url
// Example: U("https://example.com/api", "/user/create") -> "https://example.com/api/user/create"
func U(base, path string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", base, strings.TrimPrefix(path, "/")))
	if err != nil {
		return nil, fmt.Errorf("parse url [%s]: %w", path, err)
	}
	return u, nil
}
