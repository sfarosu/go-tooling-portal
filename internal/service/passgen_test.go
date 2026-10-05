package service

import (
	"strings"
	"testing"
)

func TestGeneratePassword_Valid_AllFlags(t *testing.T) {
	pwd, err := GeneratePassword(16, true, true, true, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(pwd) != 16 {
		t.Fatalf("Expected length 16, got %d", len(pwd))
	}
	// Basic sanity checks
	if !strings.ContainsAny(pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Errorf("Expected at least one uppercase letter in %s", pwd)
	}
	if !strings.ContainsAny(pwd, "abcdefghijklmnopqrstuvwxyz") {
		t.Errorf("Expected at least one lowercase letter in %s", pwd)
	}
	if !strings.ContainsAny(pwd, "0123456789") {
		t.Errorf("Expected at least one digit in %s", pwd)
	}
	if !strings.ContainsAny(pwd, ";#$%&'()*+,-.:;<=>?@[]^_`{|}~") {
		t.Errorf("Expected at least one special char in %s", pwd)
	}
}

func TestGeneratePassword_SingleClass_Numbers(t *testing.T) {
	pwd, err := GeneratePassword(8, false, false, true, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(pwd) != 8 {
		t.Fatalf("Expected length 8, got %d", len(pwd))
	}
	for _, r := range pwd {
		if !strings.ContainsAny(string(r), "0123456789") {
			t.Fatalf("Expected only digits in %s", pwd)
		}
	}
}

func TestGeneratePassword_InvalidLength(t *testing.T) {
	_, err := GeneratePassword(0, true, true, true, true)
	if err == nil {
		t.Fatal("Expected error for length 0, got nil")
	}
}

func TestGeneratePassword_NoClassSelected(t *testing.T) {
	_, err := GeneratePassword(10, false, false, false, false)
	if err == nil {
		t.Fatal("Expected error when no char classes selected, got nil")
	}
}
