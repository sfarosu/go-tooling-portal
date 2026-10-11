package apis

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateHTMXErrorEscapesHTML(t *testing.T) {
	result := generateHTMXError(fmt.Errorf("<script>alert(1)</script>"))
	if strings.Contains(result, "<script>") {
		t.Fatalf("error fragment contains unescaped markup: %s", result)
	}
	if !strings.Contains(result, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Errorf("expected escaped error message in fragment, got: %s", result)
	}
}
