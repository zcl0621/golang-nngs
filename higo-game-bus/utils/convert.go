package utils

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"io"
)

// ZipString zip string
func ZipString(str string) string {
	return Base64Encode(zipStr(str))
}

func zipStr(origin string) (content string) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_,_ = w.Write([]byte(origin))
	_ = w.Close()
	return b.String()
}

// Base64Encode string to base64 string
func Base64Encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

func UnzipString(str string) (string, error) {
	b, err := Base64Decode(str)
	if err != nil {
		return "", err
	}

	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, r); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Base64Decode base64 string to string
func Base64Decode(str string) ([]byte, error) {
	return b64.StdEncoding.DecodeString(str)
}
