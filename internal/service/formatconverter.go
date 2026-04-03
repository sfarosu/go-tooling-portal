package service

import (
	"fmt"

	"github.com/sfarosu/go-tooling-portal/internal/helper"
)

func FormatConvert(input string) (string, error) {
	jsonData, errJSON := helper.UnmarshalJSON([]byte(input))
	yamlData, errYAML := helper.UnmarshalYAML([]byte(input))

	// determine if the input was json or yaml by unmarshaling both and see which throws an error
	// remember, yaml.unmarshal on a json does NOT throw error
	switch {
	case errJSON == nil && errYAML == nil:
		yamlOut, err := helper.MarshalYAML(jsonData)
		if err != nil {
			return "", fmt.Errorf("failed to marshal YAML: %v", err)
		}
		return string(yamlOut), nil
	case errYAML == nil:
		jsonOut, err := helper.MarshalJSON(yamlData)
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %v", err)
		}
		return string(jsonOut), nil
	default:
		return "", fmt.Errorf("input is neither valid JSON nor YAML")
	}
}
