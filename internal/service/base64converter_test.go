package service

import "testing"

func TestBase64Convert(t *testing.T) {
	type args struct {
		input     string
		operation string
		format    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "encode standard",
			args:    args{input: "hello", operation: "encode", format: "standard"},
			want:    "aGVsbG8=",
			wantErr: false,
		},
		{
			name:    "decode standard",
			args:    args{input: "aGVsbG8=", operation: "decode", format: "standard"},
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "encode url-compatible",
			args:    args{input: "hello", operation: "encode", format: "url-compatible"},
			want:    "aGVsbG8=",
			wantErr: false,
		},
		{
			name:    "decode url-compatible",
			args:    args{input: "aGVsbG8=", operation: "decode", format: "url-compatible"},
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "invalid base64 decode",
			args:    args{input: "!!!", operation: "decode", format: "standard"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "unsupported operation",
			args:    args{input: "hello", operation: "compress", format: "standard"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "unsupported format",
			args:    args{input: "hello", operation: "encode", format: "foobar"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Convert(tt.args.input, tt.args.operation, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base64Convert() error = [%v], wantErr [%v]", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64Convert() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}
