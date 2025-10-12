package service

import (
	"strings"
	"testing"
)

func TestFormatConvert_JSONtoYAML(t *testing.T) {
	input := `{"name":"Alice","age":30}`

	out, err := FormatConvert(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !strings.Contains(out, "name: Alice") || !strings.Contains(out, "age: 30") {
		t.Errorf("Expected YAML to contain both keys, got:\n%s", out)
	}
}

func TestFormatConvert_YAMLtoJSON(t *testing.T) {
	input := "name: Alice\nage: 30\n"
	want := "{\n \"age\": 30,\n \"name\": \"Alice\"\n}"

	out, err := FormatConvert(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Remove all whitespace for comparison
	clean := func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), "\n", ""), " ", ""), "\t", "")
	}
	if clean(out) != clean(want) {
		t.Errorf("Expected JSON:\n%s\nGot:\n%s", want, out)
	}
}

func TestFormatConvert_InvalidInput(t *testing.T) {
	input := "not: valid: yaml: or: json"

	_, err := FormatConvert(input)
	if err == nil {
		t.Errorf("Expected error for invalid input, got nil")
	}
	if err.Error() != "input is neither valid JSON nor YAML" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
