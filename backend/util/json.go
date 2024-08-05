package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Encode encodes the given data structure into JSON format and returns the JSON byte slice.
func Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// Decode decodes the JSON body from the given request into the provided data structure.
func Decode(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read the request body")
	}

	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	return json.Unmarshal(body, v)
}
