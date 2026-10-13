package service

import (
	"strings"
	"testing"
)

func TestURLDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "decode URL components and query values",
			input: "https://example.com/foo+bar%21?service=ab%2Bc&product=de%24f&service=second",
			want: []string{
				"Scheme: https",
				"Host: example.com",
				"Path: /foo+bar!",
				"Key: product",
				"Value: de$f",
				"Key: service",
				"Value: ab+c",
				"Key: service",
				"Value: second",
			},
		},
		{
			name:  "trim surrounding whitespace and preserve repeated query values",
			input: " https://example.com/path?key=first&key=second ",
			want:  []string{"Value: first", "Value: second"},
		},
		{
			name:    "reject empty URL",
			input:   " \t\n ",
			wantErr: true,
		},
		{
			name:    "invalid path escape",
			input:   "https://example.com/%zz",
			wantErr: true,
		},
		{
			name:    "invalid query escape",
			input:   "https://example.com/?key=%zz",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := URLDecode(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("URLDecode() error = %v, wantErr %v", err, test.wantErr)
			}
			for _, part := range test.want {
				if !strings.Contains(got, part) {
					t.Errorf("URLDecode() output missing %q in:\n%s", part, got)
				}
			}
		})
	}
}
