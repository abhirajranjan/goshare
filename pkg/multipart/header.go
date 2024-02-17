package multipart

import (
	"errors"
	"mime"
	"net/http"
)

type headers struct {
	Boundary string
}

func getheaders(r *http.Request) (*headers, error) {
	mediatype, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return &headers{}, errors.New("incorrect content type")

	}

	if mediatype != "multipart/form-data" {
		return &headers{}, errors.New("incorrect media type")
	}

	boundary, ok := params["boundary"]
	if !ok {
		return &headers{}, errors.New("incorrect boundary")
	}

	return &headers{Boundary: boundary}, nil
}
