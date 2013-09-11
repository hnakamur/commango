package jsonutil

import (
	"bytes"
	"encoding/json"
	"strings"
)

func Encode(v interface{}) (string, error) {
	var buf bytes.Buffer
    err := json.NewEncoder(&buf).Encode(v)
    if err != nil {
        return "", err
    }
	return strings.TrimRight(buf.String(), "\n"), nil
}

func Decode(s string) (v interface{}, err error) {
	err = json.NewDecoder(bytes.NewReader([]byte(s))).Decode(&v)
	return
}
