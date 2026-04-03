package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// PrettyJSON formats a JSON string with indentation for better readability
func PrettyJSON(insertedText string) (bytes.Buffer, error) {
	// default to 2 spaces to preserve previous external behavior
	return PrettyJSONWithSpaces(insertedText, 2)
}

// PrettyJSONWithSpaces formats a JSON string with indentation using the specified number of spaces.
func PrettyJSONWithSpaces(insertedText string, spaces int) (bytes.Buffer, error) {
	var pretty bytes.Buffer

	// must be non-negative
	if spaces < 0 {
		return bytes.Buffer{}, fmt.Errorf("spaces must be a valid positive integer")
	}

	// limit spaces to prevent excessive resource usage
	const maxSpaces = 16
	if spaces > maxSpaces {
		return bytes.Buffer{}, fmt.Errorf("spaces must be <= %d", maxSpaces)
	}

	indent := strings.Repeat(" ", spaces)
	err := json.Indent(&pretty, []byte(insertedText), "", indent)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error indenting JSON: %v", err)
	}
	return pretty, nil
}

// MarshalJSON converts a Go data structure to a JSON byte slice with indentation
// MarshalJSON marshals with a default indentation of 2 spaces.
func MarshalJSON(data any) ([]byte, error) {
	return MarshalJSONWithSpaces(data, 2)
}

// MarshalJSONWithSpaces marshals a Go data structure to JSON with the specified indentation spaces.
func MarshalJSONWithSpaces(data any, spaces int) ([]byte, error) {
	prefix := ""
	indent := strings.Repeat(" ", spaces)
	byteData, err := json.MarshalIndent(data, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	return byteData, nil
}

// UnmarshalJSON converts a JSON byte slice to a Go data structure
func UnmarshalJSON(byteData []byte) (map[string]any, error) {
	var jsonData map[string]any
	err := json.Unmarshal([]byte(byteData), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	return jsonData, err
}
