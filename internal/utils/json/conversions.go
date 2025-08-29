package json

import (
	"encoding/json"
	"fmt"
)

// StructToJSON converts any struct to JSON string
func StructToJSON(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// StructToJSONIndented converts any struct to indented JSON string for better readability
func StructToJSONIndented(data interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal struct to indented JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// StructToJSONBytes converts any struct to JSON bytes
func StructToJSONBytes(data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON bytes: %w", err)
	}
	return jsonBytes, nil
}
