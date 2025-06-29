package helper

import (
	"os"
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	var err error

	type args struct {
		filePath string
	}
	// Create a temporary file with known content
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write known content to the temporary file
	content := []byte("test-content")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "read existing file",
			args:    args{filePath: tmpFile.Name()},
			want:    content,
			wantErr: false,
		},
		{
			name:    "read non-existent file",
			args:    args{filePath: "does_not_exist.txt"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
