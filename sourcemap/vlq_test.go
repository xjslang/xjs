package sourcemap

import (
	"strings"
	"testing"
)

// encodeVLQSegment encodes a source map segment (up to 5 fields) as VLQ.
// Fields: generatedColumn, sourceIndex, sourceLine, sourceColumn, [nameIndex]
func encodeVLQSegment(fields ...int) string {
	var result strings.Builder
	for _, field := range fields {
		result.WriteString(encodeVLQ(field))
	}
	return result.String()
}

func TestEncodeVLQ(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "A"},
		{1, "C"},
		{-1, "D"},
		{2, "E"},
		{-2, "F"},
		{15, "e"},
		{16, "gB"},
		{-16, "hB"},
		{100, "oG"},
		{-100, "pG"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := encodeVLQ(tt.input)
			if result != tt.expected {
				t.Errorf("encodeVLQ(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeVLQSegment(t *testing.T) {
	// Test encoding a typical segment: column=0, source=0, line=0, column=0
	result := encodeVLQSegment(0, 0, 0, 0)
	if result != "AAAA" {
		t.Errorf("encodeVLQSegment(0,0,0,0) = %q, want %q", result, "AAAA")
	}

	// Test encoding with some offsets
	result = encodeVLQSegment(5, 0, 10, 3)
	if result == "" {
		t.Error("encodeVLQSegment should not return empty string")
	}
}
