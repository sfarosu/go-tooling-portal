package service

import (
    "github.com/sfarosu/go-tooling-portal/internal/helper"
    "encoding/json"
    "fmt"
)

// PrettyJSON formats a JSON string with indentation and returns the formatted string
func PrettyJSON(input string) (string, error) {
    buf, err := helper.PrettyJSON(input)
    if err != nil {
        return "", err
    }
    return buf.String(), nil
}

// PrettyJSONWithSpaces formats JSON with a custom number of spaces for indentation.
func PrettyJSONWithSpaces(input string, spaces int) (string, error) {
    buf, err := helper.PrettyJSONWithSpaces(input, spaces)
    if err != nil {
        return "", err
    }
    return buf.String(), nil
}

// UnprettyJSON minifies JSON input (removes indentation/whitespace)
func UnprettyJSON(input string) (string, error) {
    var v any
    if err := json.Unmarshal([]byte(input), &v); err != nil {
        return "", fmt.Errorf("invalid JSON: %v", err)
    }
    out, err := json.Marshal(v)
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON: %v", err)
    }
    return string(out), nil
}
