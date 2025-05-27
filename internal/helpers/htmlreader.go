package helpers

import (
	"bytes"
	"golang.org/x/net/html/charset"
	"io"
)

func GetHTMLReader(body []byte, contentType string) (io.Reader, error) {
	return charset.NewReader(bytes.NewReader(body), contentType)
}
