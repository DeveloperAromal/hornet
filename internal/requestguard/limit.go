package requestguard

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var (
	ErrURLTooLong   = errors.New("request URL exceeds maximum length")
	ErrBodyTooLarge = errors.New("request body exceeds maximum size")
)

type Limits struct {
	MaxURLLength int
	MaxBodySize  int64
}

func (lmt Limits) ValidateURL(req *http.Request) error {
	if lmt.MaxURLLength <= 0 {
		return nil
	}

	if len(req.URL.RequestURI()) > lmt.MaxURLLength {
		return ErrURLTooLong
	}

	return nil
}

func (lmt *Limits) ValidateBody(req *http.Request) error {

	if lmt.MaxBodySize <= 0 || req.Body == nil {
		return nil
	}

	originalBody := req.Body
	defer originalBody.Close()

	body, err := io.ReadAll(http.MaxBytesReader(nil, originalBody, lmt.MaxBodySize))
	if err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			return ErrBodyTooLarge
		}
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	return nil
}
