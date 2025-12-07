package utils

import (
	"testing"
)

func TestConvertDateLayout(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{
			format:   "Y-M-D",
			expected: "2006-01-02",
		},
		{
			format:   "Y-M-D h:m:s",
			expected: "2006-01-02 15:04:05",
		},
		{
			format:   "YMD_hms",
			expected: "20060102_150405",
		},
		{
			format:   "YMD_hms.f",
			expected: "20060102_150405.000",
		},
	}

	for _, test := range tests {
		result := ConvertDateLayout(test.format)
		if result != test.expected {
			t.Errorf("ConvertDateLayout(%s) = %s, expected %s", test.format, result, test.expected)
		}
	}
}
