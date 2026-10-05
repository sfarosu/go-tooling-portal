package helper

import (
	"fmt"
	"os"
)

// ReadFile reads the content of a file and returns it as a byte slice
func ReadFile(filePath string) ([]byte, error) {
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file in path [%v]: %v", filePath, err)
	}
	return byteData, nil
}
