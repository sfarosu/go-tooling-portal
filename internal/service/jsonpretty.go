package service

import (
    "github.com/sfarosu/go-tooling-portal/internal/helper"
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
