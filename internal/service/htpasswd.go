package service

import (
	"fmt"

	"github.com/GehirnInc/crypt/apr1_crypt"
	"github.com/GehirnInc/crypt/md5_crypt"
	"github.com/GehirnInc/crypt/sha256_crypt"
	"github.com/GehirnInc/crypt/sha512_crypt"
)

// GenerateHtpasswd generates an htpasswd entry for the given username and password
// using the specified algorithm ("apr1", "1", "5", or "6")
// Returns the formatted htpasswd line or an error if the algorithm is unsupported or hashing fails
func GenerateHtpasswd(username string, password string, algorithm string) (string, error) {
	switch algorithm {
	case "apr1":
		hash, err := generateAPR1crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	case "1":
		hash, err := generateMD5crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	case "5":
		hash, err := generateSHA256crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	case "6":
		hash, err := generateSHA512crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	default:
		return "", fmt.Errorf("unsupported algorithm [%v]; use [apr1], [1], [5] or [6] openssl cryptographic options", algorithm)
	}
}

// generateAPR1crypt generates an APR1 (Apache MD5) htpasswd entry for the given username and password
func generateAPR1crypt(username string, password string) (string, error) {
	hash, err := apr1_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", fmt.Errorf("error generating SHA512 htpasswd for user [%v]: %v", username, err)
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}

// generateMD5crypt generates an MD5 htpasswd entry for the given username and password
func generateMD5crypt(username string, password string) (string, error) {
	hash, err := md5_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", fmt.Errorf("error generating MD5 htpasswd for user [%v]: %v", username, err)
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}

// generateSHA256crypt generates a SHA256 htpasswd entry for the given username and password
func generateSHA256crypt(username string, password string) (string, error) {
	hash, err := sha256_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", fmt.Errorf("error generating SHA256 htpasswd for user [%v]: %v", username, err)
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}

// generateSHA512crypt generates a SHA512 htpasswd entry for the given username and password
func generateSHA512crypt(username string, password string) (string, error) {
	hash, err := sha512_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", fmt.Errorf("error generating SHA512 htpasswd for [%v]: %v", username, err)
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}
