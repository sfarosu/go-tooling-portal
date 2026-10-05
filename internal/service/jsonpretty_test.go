package service

import (
    "encoding/json"
    "reflect"
    "testing"
)

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

func TestUnprettyJSON_Valid(t *testing.T) {
    input := `{
  "name": "Alice",
  "age": 30
}`
    out, err := UnprettyJSON(input)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    // The order of object keys is not guaranteed; compare via unmarshalling
    var got map[string]any
    if err := json.Unmarshal([]byte(out), &got); err != nil {
        t.Fatalf("output is not valid JSON: %v", err)
    }
    want := map[string]any{"name": "Alice", "age": float64(30)}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("unexpected minified output content: got=%v want=%v", got, want)
    }
}

func TestUnprettyJSON_Invalid(t *testing.T) {
    input := `{name:Alice,age:30}`
    _, err := UnprettyJSON(input)
    if err == nil {
        t.Fatalf("expected error for invalid JSON in UnprettyJSON")
    }
}
