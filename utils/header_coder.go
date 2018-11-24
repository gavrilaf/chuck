package utils

import (
	"encoding/json"
	"net/http"
)

func EncodeHeaders(header http.Header) ([]byte, error) {
	buf, err := json.MarshalIndent(header, "", "\t")
	return buf, err
}

func DecodeHeaders(buf []byte) (http.Header, error) {
	var header http.Header
	err := json.Unmarshal(buf, &header)
	if err != nil {
		return nil, err
	}

	return header, nil
}
