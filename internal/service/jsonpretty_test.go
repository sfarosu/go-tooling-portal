package service

import "testing"

func TestPrettyJSON_Valid(t *testing.T) {
    input := `{"name":"Alice","age":30}`
    out, err := PrettyJSON(input)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if out == "" {
        t.Fatalf("expected non-empty output")
    }

    // also test custom spaces
    out2, err := PrettyJSONWithSpaces(input, 4)
    if err != nil {
        t.Fatalf("expected no error for custom spaces, got %v", err)
    }
    if out2 == "" {
        t.Fatalf("expected non-empty output for custom spaces")
    }
}

func TestPrettyJSON_Invalid(t *testing.T) {
    input := `{name:Alice,age:30}`
    _, err := PrettyJSON(input)
    if err == nil {
        t.Fatalf("expected error for invalid JSON")
    }
}
