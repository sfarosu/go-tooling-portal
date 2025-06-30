package service

import (
	b64 "encoding/base64"
	"fmt"
)

// Base64Convert encodes or decodes a string using base64.
func Base64Convert(input, operation, format string) (string, error) {
	switch operation {
	case "decode":
		switch format {
		case "standard":
			r, err := b64.StdEncoding.DecodeString(input)
			if err != nil {
				return "", fmt.Errorf("error decoding in standard format: %v", err)
			}
			return string(r), nil
		case "url-compatible":
			r, err := b64.URLEncoding.DecodeString(input)
			if err != nil {
				return "", fmt.Errorf("error decoding in url-compatible format: %v", err)
			}
			return string(r), nil
		}
	case "encode":
		switch format {
		case "standard":
			return b64.StdEncoding.EncodeToString([]byte(input)), nil
		case "url-compatible":
			return b64.URLEncoding.EncodeToString([]byte(input)), nil
		}
	}
	return "", fmt.Errorf("unsupported operation or format: operation=[%v], format=[%v]", operation, format)
}
