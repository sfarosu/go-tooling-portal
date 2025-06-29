package helper

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestPrettyJSON(t *testing.T) {
	tests := []struct {
		name         string
		insertedText string
		want         string
		wantErr      bool
	}{
		{
			name:         "valid JSON",
			insertedText: `{"foo":"bar","baz":123}`,
			want: `{
  "baz": 123,
  "foo": "bar"
}
`,
			wantErr: false,
		},
		{
			name:         "invalid JSON",
			insertedText: `{"foo":}`,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "empty input",
			insertedText: "",
			want:         "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrettyJSON(tt.insertedText)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrettyJSON() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check that the output is valid JSON
				var v interface{}
				if err := json.Unmarshal(got.Bytes(), &v); err != nil {
					t.Errorf("PrettyJSON() output is not valid JSON: %v", err)
				}
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string // We'll check for substring presence due to map key order
		wantErr bool
	}{
		{
			name:    "simple map",
			input:   map[string]interface{}{"foo": "bar", "baz": 123},
			want:    `"foo": "bar"`, // Just check for presence
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			want:    "null",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(string(got), tt.want) {
				t.Errorf("MarshalJSON() = [%v], want substring [%v]", string(got), tt.want)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid JSON",
			input:   []byte(`{"foo":"bar","baz":123}`),
			want:    map[string]interface{}{"foo": "bar", "baz": float64(123)},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{"foo":}`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte(""),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}
