package sourcemap

import "strings"

const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// encodeVLQ encodes a single integer using Base64 VLQ encoding.
// VLQ (Variable Length Quantity) allows encoding arbitrarily large integers
// in a compact format suitable for source maps.
func encodeVLQ(n int) string {
	var result strings.Builder

	// Convert to VLQ signed format: if negative, set LSB to 1
	if n < 0 {
		n = (-n << 1) | 1
	} else {
		n = n << 1
	}

	// Encode in 5-bit chunks with continuation bit
	for {
		digit := n & 0x1F // Take 5 bits
		n >>= 5
		if n > 0 {
			digit |= 0x20 // Set continuation bit
		}
		result.WriteByte(base64Chars[digit])
		if n == 0 {
			break
		}
	}

	return result.String()
}
