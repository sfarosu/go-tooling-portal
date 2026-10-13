package service

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/sfarosu/go-tooling-portal/internal/helper"
)

type URLDecodedData struct {
	Scheme    string `yaml:"Scheme"`
	Host      string `yaml:"Host"`
	Path      string `yaml:"Path"`
	QueryArgs []struct {
		Key   string `yaml:"Key"`
		Value string `yaml:"Value"`
	} `yaml:"QueryArgs"`
}

// URLDecode parses a URL and returns its components as YAML.
func URLDecode(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %w", err)
	}

	query, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("error parsing URL query: %w", err)
	}

	decoded := URLDecodedData{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   parsedURL.Path,
	}
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		values := query[key]
		for _, value := range values {
			decoded.QueryArgs = append(decoded.QueryArgs, struct {
				Key   string `yaml:"Key"`
				Value string `yaml:"Value"`
			}{Key: key, Value: value})
		}
	}

	result, err := helper.MarshalYAML(decoded)
	if err != nil {
		return "", fmt.Errorf("error formatting decoded URL: %w", err)
	}
	return string(result), nil
}
