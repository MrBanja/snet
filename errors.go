package snet

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// //////////////////////////////////////////////////////////////
// WrongStatusError

var _ error = WrongStatusError{}

type WrongStatusError struct {
	ResponseBody       []byte `json:"response_body"`
	ResponseStatusCode int    `json:"response_status_code"`
	RequestURL         string `json:"request_url"`
	RequestMethod      string `json:"request_method"`
}

func (err WrongStatusError) Error() string {
	return fmt.Sprintf(
		"wrong status code for [%s %s]; Code: [%d] with reponse [%s] ",
		err.RequestMethod, err.RequestURL,
		err.ResponseStatusCode, err.ResponseBody,
	)
}

func NewWrongStatusError(r *http.Response) error {
	defer func() { _ = r.Body.Close() }()
	b, _ := io.ReadAll(r.Body)
	return WrongStatusError{
		ResponseBody:       b,
		ResponseStatusCode: r.StatusCode,
		RequestURL:         r.Request.URL.String(),
		RequestMethod:      r.Request.Method,
	}
}

func IsWrongStatusError(err error) *WrongStatusError {
	var e WrongStatusError
	if errors.As(err, &e) {
		return &e
	} else {
		return nil
	}
}
