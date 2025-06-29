package helper

import (
	"reflect"
	"strings"
	"testing"
)

func TestMarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []string // We'll check for presence of these substrings due to map key order
		wantErr bool
	}{
		{
			name:    "simple map",
			input:   map[string]interface{}{"foo": "bar", "baz": 123},
			want:    []string{"foo: bar", "baz: 123"},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			want:    []string{"null"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalYAML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotStr := string(got)
				for _, substr := range tt.want {
					if !strings.Contains(gotStr, substr) {
						t.Errorf("MarshalYAML() output missing substring %q in %q", substr, gotStr)
					}
				}
			}
		})
	}
}

func TestUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid YAML",
			input:   []byte("foo: bar\nbaz: 123\n"),
			want:    map[string]interface{}{"foo": "bar", "baz": 123},
			wantErr: false,
		},
		{
			name:    "invalid YAML",
			input:   []byte("foo: ["),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte(""),
			want:    map[string]interface{}{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalYAML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Normalize nil and empty map for comparison
				if got == nil && tt.want == nil {
					return
				}
				if got == nil && len(tt.want) == 0 {
					return
				}
				if tt.want == nil && len(got) == 0 {
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("UnmarshalYAML() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
