package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ConvertOperationsToString(operations json.RawMessage) (string, error) {
	var opsMap map[string]string
	err := json.Unmarshal(operations, &opsMap)
	if err != nil {
		return "", err
	}

	// Convert the map to the desired string format
	var ops []string
	for key, value := range opsMap {
		ops = append(ops, fmt.Sprintf("%s:%s", key, value))
	}
	return strings.Join(ops, ","), nil
}
