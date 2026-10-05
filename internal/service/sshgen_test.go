package service

import (
	"strings"
	"testing"
)

func TestGenerateSSHKeyPair_Ed25519(t *testing.T) {
	priv, pub, err := GenerateSSHKeyPair("ed25519", "", "", "alice@example.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !strings.Contains(pub, "ssh-ed25519") {
		t.Errorf("Public key does not contain 'ssh-ed25519': %s", pub)
	}
	if !strings.Contains(pub, "alice@example.com") {
		t.Errorf("Public key does not contain email: %s", pub)
	}
	if !strings.HasPrefix(priv, "-----BEGIN OPENSSH PRIVATE KEY-----") {
		t.Errorf("Private key does not have expected header: %s", priv)
	}
}

func TestGenerateSSHKeyPair_ECDSA(t *testing.T) {
	priv, pub, err := GenerateSSHKeyPair("ecdsa", "256", "", "bob@example.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !strings.Contains(pub, "ecdsa-sha2-nistp256") {
		t.Errorf("Public key does not contain 'ecdsa-sha2-nistp256': %s", pub)
	}
	if !strings.Contains(pub, "bob@example.com") {
		t.Errorf("Public key does not contain email: %s", pub)
	}
	if !strings.HasPrefix(priv, "-----BEGIN EC PRIVATE KEY-----") {
		t.Errorf("Private key does not have expected header: %s", priv)
	}
}

func TestGenerateSSHKeyPair_RSA(t *testing.T) {
	priv, pub, err := GenerateSSHKeyPair("rsa", "", "2048", "carol@example.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !strings.Contains(pub, "ssh-rsa") {
		t.Errorf("Public key does not contain 'ssh-rsa': %s", pub)
	}
	if !strings.Contains(pub, "carol@example.com") {
		t.Errorf("Public key does not contain email: %s", pub)
	}
	if !strings.HasPrefix(priv, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Errorf("Private key does not have expected header: %s", priv)
	}
}

func TestGenerateSSHKeyPair_InvalidAlgorithm(t *testing.T) {
	_, _, err := GenerateSSHKeyPair("invalid", "", "", "nobody@example.com")
	if err == nil {
		t.Error("Expected error for invalid algorithm, got nil")
	}
}

func TestGenerateSSHKeyPair_MissingBits(t *testing.T) {
	_, _, err := GenerateSSHKeyPair("ecdsa", "", "", "nobody@example.com")
	if err == nil {
		t.Error("Expected error for missing ecdsaBits, got nil")
	}
	_, _, err = GenerateSSHKeyPair("rsa", "", "", "nobody@example.com")
	if err == nil {
		t.Error("Expected error for missing rsaBits, got nil")
	}
}
