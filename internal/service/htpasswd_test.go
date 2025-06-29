package service

import (
	"strings"
	"testing"
)

func TestGenerateHtpasswd(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		password  string
		algorithm string
		wantErr   bool
		wantAlgo  string // substring to check in hash
	}{
		{
			name:      "APR1 algorithm",
			username:  "alice",
			password:  "password123",
			algorithm: "apr1",
			wantErr:   false,
			wantAlgo:  "$apr1$",
		},
		{
			name:      "MD5 algorithm",
			username:  "bob",
			password:  "password123",
			algorithm: "1",
			wantErr:   false,
			wantAlgo:  "$1$",
		},
		{
			name:      "SHA256 algorithm",
			username:  "carol",
			password:  "password123",
			algorithm: "5",
			wantErr:   false,
			wantAlgo:  "$5$",
		},
		{
			name:      "SHA512 algorithm",
			username:  "dave",
			password:  "password123",
			algorithm: "6",
			wantErr:   false,
			wantAlgo:  "$6$",
		},
		{
			name:      "unsupported algorithm",
			username:  "eve",
			password:  "password123",
			algorithm: "sha256",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateHtpasswd(tt.username, tt.password, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateHtpasswd() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check username is present and hash contains expected algorithm marker
				parts := strings.SplitN(got, ":", 2)
				if len(parts) != 2 {
					t.Errorf("GenerateHtpasswd() output format invalid: [%v]", got)
				} else {
					if parts[0] != tt.username {
						t.Errorf("GenerateHtpasswd() username = [%v], want [%v]", parts[0], tt.username)
					}
					if !strings.Contains(parts[1], tt.wantAlgo) {
						t.Errorf("GenerateHtpasswd() hash = [%v], want to contain [%v]", parts[1], tt.wantAlgo)
					}
				}
			}
		})
	}
}

func Test_generateAPR1crypt(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid input",
			username: "alice",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			username: "bob",
			password: "",
			wantErr:  false, // APR1 allows empty passwords, but it can be set to true for enforcement
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateAPR1crypt(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateAPR1crypt() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parts := strings.SplitN(got, ":", 2)
				if len(parts) != 2 {
					t.Errorf("generateAPR1crypt() output format invalid: [%v]", got)
				} else {
					if parts[0] != tt.username {
						t.Errorf("generateAPR1crypt() username = [%v], want [%v]", parts[0], tt.username)
					}
					if !strings.HasPrefix(parts[1], "$apr1$") {
						t.Errorf("generateAPR1crypt() hash = [%v], want prefix [%v]", parts[1], "$apr1$")
					}
				}
			}
		})
	}
}

func Test_generateMD5crypt(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid input",
			username: "alice",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			username: "bob",
			password: "",
			wantErr:  false, // MD5 allows empty passwords
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateMD5crypt(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateMD5crypt() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parts := strings.SplitN(got, ":", 2)
				if len(parts) != 2 {
					t.Errorf("generateMD5crypt() output format invalid: [%v]", got)
				} else {
					if parts[0] != tt.username {
						t.Errorf("generateMD5crypt() username = [%v], want [%v]", parts[0], tt.username)
					}
					if !strings.HasPrefix(parts[1], "$1$") {
						t.Errorf("generateMD5crypt() hash = [%v], want prefix [%v]", parts[1], "$1$")
					}
				}
			}
		})
	}
}

func Test_generateSHA256crypt(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid input",
			username: "carol",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			username: "dave",
			password: "",
			wantErr:  false, // SHA256 allows empty passwords
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSHA256crypt(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSHA256crypt() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parts := strings.SplitN(got, ":", 2)
				if len(parts) != 2 {
					t.Errorf("generateSHA256crypt() output format invalid: [%v]", got)
				} else {
					if parts[0] != tt.username {
						t.Errorf("generateSHA256crypt() username = [%v], want [%v]", parts[0], tt.username)
					}
					if !strings.HasPrefix(parts[1], "$5$") {
						t.Errorf("generateSHA256crypt() hash = [%v], want prefix [%v]", parts[1], "$5$")
					}
				}
			}
		})
	}
}

func Test_generateSHA512crypt(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid input",
			username: "eve",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			username: "frank",
			password: "",
			wantErr:  false, // SHA512 allows empty passwords
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSHA512crypt(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSHA512crypt() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parts := strings.SplitN(got, ":", 2)
				if len(parts) != 2 {
					t.Errorf("generateSHA512crypt() output format invalid: [%v]", got)
				} else {
					if parts[0] != tt.username {
						t.Errorf("generateSHA512crypt() username = [%v], want [%v]", parts[0], tt.username)
					}
					if !strings.HasPrefix(parts[1], "$6$") {
						t.Errorf("generateSHA512crypt() hash = [%v], want prefix [%v]", parts[1], "$6$")
					}
				}
			}
		})
	}
}
