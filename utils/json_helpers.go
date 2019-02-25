package utils

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"mime"
	"net/http"
	"strings"
)

func IsRespHasJsonContent(r *http.Response) bool {
	contentType := r.Header.Get("Content-type")
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err == nil && strings.Contains(t, "json") {
			return true
		}
	}
	return false
}

func FormatJson(b []byte) []byte {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return b
	}

	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetIndent("", "\t")
	if err := enc.Encode(data); err != nil {
		return b
	}

	return out.Bytes()
}
