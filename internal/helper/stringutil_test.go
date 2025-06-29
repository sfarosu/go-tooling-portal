package helper

import (
	"fmt"
	"testing"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		Uppercase bool
		Lowercase bool
		Numbers   bool
		Specials  bool
	}{
		{
			name:      "all categories",
			size:      16,
			Uppercase: true,
			Lowercase: true,
			Numbers:   true,
			Specials:  true,
		},
		{
			name:      "only uppercase",
			size:      8,
			Uppercase: true,
		},
		{
			name:      "only lowercase",
			size:      8,
			Lowercase: true,
		},
		{
			name:    "only numbers",
			size:    8,
			Numbers: true,
		},
		{
			name:     "only specials",
			size:     8,
			Specials: true,
		},
		{
			name:      "uppercase and numbers",
			size:      12,
			Uppercase: true,
			Numbers:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomString(tt.size, tt.Uppercase, tt.Lowercase, tt.Numbers, tt.Specials)
			if len(got) != tt.size {
				t.Errorf("RandomString() length = [%v], want [%v]", len(got), tt.size)
			}
		})
	}
}

func TestAddSecondDigit(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "00"},
		{1, "01"},
		{5, "05"},
		{9, "09"},
		{10, "10"},
		{23, "23"},
		{99, "99"},
		{100, "100"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=%d", tt.input), func(t *testing.T) {
			got := AddSecondDigit(tt.input)
			if got != tt.want {
				t.Errorf("AddSecondDigit(%v) = [%v], want [%v]", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetNumberDigitsAmmount(t *testing.T) {
	tests := []struct {
		input int64
		want  int
	}{
		{0, 0},
		{1, 1},
		{9, 1},
		{10, 2},
		{99, 2},
		{100, 3},
		{12345, 5},
		{1000000, 7},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input=[%v]", tt.input), func(t *testing.T) {
			got := GetNumberDigitsAmmount(tt.input)
			if got != tt.want {
				t.Errorf("GetNumberDigitsAmmount(%v) = [%v], want [%v]", tt.input, got, tt.want)
			}
		})
	}
}
