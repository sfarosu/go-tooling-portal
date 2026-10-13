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

func TestFormatConvert_JSONArrayToYAML(t *testing.T) {
	out, err := FormatConvert(`[{"name":"Alice"}, 2, true]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "- name: Alice") || !strings.Contains(out, "- 2") || !strings.Contains(out, "- true") {
		t.Errorf("expected YAML sequence with object and scalar items, got:\n%s", out)
	}
}

func TestFormatConvert_YAMLSequenceToJSON(t *testing.T) {
	out, err := FormatConvert("- Alice\n- 2\n- true\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[\n  \"Alice\",\n  2,\n  true\n]"
	if out != want {
		t.Errorf("expected JSON:\n%s\ngot:\n%s", want, out)
	}
}

func TestFormatConvert_JSONScalarToYAML(t *testing.T) {
	out, err := FormatConvert(`"hello"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello\n" {
		t.Errorf("expected YAML scalar %q, got %q", "hello\n", out)
	}
}

func TestFormatConvert_YAMLScalarToJSON(t *testing.T) {
	out, err := FormatConvert("hello\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `"hello"` {
		t.Errorf("expected JSON scalar %q, got %q", `"hello"`, out)
	}
}

func TestFormatConvert_InvalidInput(t *testing.T) {
	input := "not: [valid"

	_, err := FormatConvert(input)
	if err == nil {
		t.Errorf("Expected error for invalid input, got nil")
	}
	if err.Error() != "input is neither valid JSON nor YAML" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
