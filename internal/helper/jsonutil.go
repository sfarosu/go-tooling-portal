package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PrettyJSON formats a JSON string with indentation for better readability
func PrettyJSON(insertedText string) (bytes.Buffer, error) {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(insertedText), "", "    ")
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error indenting JSON: %v", err)
	}
	return pretty, nil
}

// MarshalJSON converts a Go data structure to a JSON byte slice with indentation
func MarshalJSON(data interface{}) ([]byte, error) {
	byteData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	return byteData, nil
}

// UnmarshalJSON converts a JSON byte slice to a Go data structure
func UnmarshalJSON(byteData []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(byteData), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	return jsonData, err
}
