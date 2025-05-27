package helpers

import (
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
)

func GetHTMLReader(resp *http.Response) (io.Reader, error) {
	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return reader, nil
}
