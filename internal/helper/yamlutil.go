package helper

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// MarshalYAML converts a Go data structure to a YAML byte slice with indentation
func MarshalYAML(data interface{}) ([]byte, error) {
	var byteData bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&byteData)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling YAML: %v", err)
	}
	return byteData.Bytes(), nil
}

// UnmarshalYAML converts a YAML byte slice to a Go data structure
func UnmarshalYAML(byteData []byte) (map[string]interface{}, error) {
	var yamlData map[string]interface{}
	err := yaml.Unmarshal([]byte(byteData), &yamlData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}
	return yamlData, err
}
