package utils

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"strconv"
)

// convert string to int
func StringToInt(str string) int {
	var i int
	i, _ = strconv.Atoi(str)
	return i
}

// zipString zip string
func ZipString(str string) string {
	return Base64Encode(zipStr(str))
}

func zipStr(origin string) (content string) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(origin))
	w.Close()
	return b.String()
}

// string to base64 string
func Base64Encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}
