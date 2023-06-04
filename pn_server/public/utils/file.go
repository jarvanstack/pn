package utils

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"strings"
)

// FileToBase64 file to base64
func FileToBase64(file io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Base64ToReader
func Base64ToReader(base64Str string) (io.Reader, error) {
	// 如果前面有 data:image/jpeg;base64, 这种前缀，需要去掉
	if i := strings.Index(base64Str, ","); i > 0 {
		base64Str = base64Str[i+1:]
	}
	bs, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bs), nil
}
