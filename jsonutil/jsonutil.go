package jsonutil

import (
	"encoding/json"
)

func Encode(v interface{}) (string, error) {
    bytes, err := json.Marshal(v)
    if err != nil {
        return "", err
    }
	return string(bytes), nil
}

func Decode(s string) (v interface{}, err error) {
	err = json.Unmarshal([]byte(s), &v)
	return
}
