// Package lexer provides lexical analysis functionality for the XJS language.
// It tokenizes source code into a sequence of tokens that can be consumed by the parser.
package lexer

import (
	"strings"

	"github.com/xjslang/xjs/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	CurrentChar  byte // current char under examination
	Line         int  // current line
	Column       int  // current column

	nextToken func(*Lexer) token.Token
}

// New creates a new lexer instance
func New(input string, middlewares ...func(l *Lexer, next func() token.Token) token.Token) *Lexer {
	l := &Lexer{
		input:  input,
		Line:   1,
		Column: 0,

		nextToken: baseNextToken,
	}
	for _, md := range middlewares {
		l.UseTokenReader(md)
	}
	l.ReadChar()
	return l
}

// ReadChar reads the next character and advances position in the input
func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.input) {
		l.CurrentChar = 0 // ASCII NUL character represents "EOF"
	} else {
		l.CurrentChar = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.CurrentChar == '\n' {
		l.Line++
		l.Column = 0
	} else {
		l.Column++
	}
}

func (l *Lexer) PeekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.CurrentChar == ' ' || l.CurrentChar == '\t' || l.CurrentChar == '\n' || l.CurrentChar == '\r' {
		l.ReadChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.CurrentChar) || isDigit(l.CurrentChar) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or decimal)
func (l *Lexer) readNumber() (string, token.Type) {
	position := l.position
	tokenType := token.INT

	for isDigit(l.CurrentChar) {
		l.ReadChar()
	}

	// Check if it's a decimal number
	if l.CurrentChar == '.' && isDigit(l.PeekChar()) {
		tokenType = token.FLOAT
		l.ReadChar() // consume the '.'
		for isDigit(l.CurrentChar) {
			l.ReadChar()
		}
	}

	return l.input[position:l.position], tokenType
}

// readString reads a string literal
func (l *Lexer) readString(delimiter byte) string {
	var result strings.Builder

	for {
		l.ReadChar()
		if l.CurrentChar == 0 {
			break
		}
		// Handle escape sequences
		if l.CurrentChar == '\\' {
			l.ReadChar() // Move to the character after backslash
			if l.CurrentChar == 'x' {
				// Handle hexadecimal escape sequence \xHH
				hex1 := l.PeekChar()
				if isHexDigit(hex1) {
					l.ReadChar() // consume first hex digit
					hex2 := l.PeekChar()
					if isHexDigit(hex2) {
						l.ReadChar() // consume second hex digit
						// Convert hex digits to byte value
						value := hexDigitValue(hex1)*16 + hexDigitValue(hex2)
						result.WriteByte(byte(value))
						continue
					}
				}
				// If not a valid hex sequence, treat as literal \x
				result.WriteByte('\\')
				result.WriteByte('x')
				continue
			} else if l.CurrentChar == 'u' {
				// Check if it's extended Unicode \u{...}
				if l.PeekChar() == '{' {
					// Handle extended Unicode escape sequence \u{H...}
					l.ReadChar() // consume the '{'
					var hexDigits []byte
					isValid := true

					// Read hex digits until we find '}'
					for {
						nextChar := l.PeekChar()
						if nextChar == '}' {
							l.ReadChar() // consume the '}'
							break
						}
						if !isHexDigit(nextChar) || len(hexDigits) >= 6 {
							// Invalid sequence - mark as invalid and break
							isValid = false
							break
						}
						l.ReadChar()
						hexDigits = append(hexDigits, l.CurrentChar)
					}

					// Validate the sequence
					if !isValid || len(hexDigits) == 0 || len(hexDigits) > 6 {
						// Treat as literal \u{...
						result.WriteByte('\\')
						result.WriteByte('u')
						result.WriteByte('{')
						for _, digit := range hexDigits {
							result.WriteByte(digit)
						}
						if isValid {
							result.WriteByte('}')
						}
						continue
					}

					// Convert hex digits to value
					value := 0
					for _, digit := range hexDigits {
						value = value*16 + hexDigitValue(digit)
					}

					// Validate Unicode range
					if value > 0x10FFFF {
						// Invalid Unicode code point - treat as literal
						result.WriteByte('\\')
						result.WriteByte('u')
						result.WriteByte('{')
						for _, digit := range hexDigits {
							result.WriteByte(digit)
						}
						result.WriteByte('}')
						continue
					}

					// Convert to UTF-8 and add to result
					utf8Bytes := encodeUTF8(value)
					for _, b := range utf8Bytes {
						result.WriteByte(b)
					}
					continue
				} else {
					// Handle regular Unicode escape sequence \uHHHH
					hex1 := l.PeekChar()
					if isHexDigit(hex1) {
						l.ReadChar() // consume first hex digit
						hex2 := l.PeekChar()
						if isHexDigit(hex2) {
							l.ReadChar() // consume second hex digit
							hex3 := l.PeekChar()
							if isHexDigit(hex3) {
								l.ReadChar() // consume third hex digit
								hex4 := l.PeekChar()
								if isHexDigit(hex4) {
									l.ReadChar() // consume fourth hex digit
									// Convert 4 hex digits to Unicode value
									value := hexDigitValue(hex1)*4096 + hexDigitValue(hex2)*256 + hexDigitValue(hex3)*16 + hexDigitValue(hex4)
									// Convert to UTF-8 and write the bytes
									utf8Bytes := encodeUTF8(value)
									for _, b := range utf8Bytes {
										result.WriteByte(b)
									}
									continue
								}
							}
						}
					}
					// If not a valid Unicode sequence, treat as literal \u
					result.WriteByte('\\')
					result.WriteByte('u')
					continue
				}
			} else {
				// Keep escape sequences as-is for valid JavaScript output
				switch l.CurrentChar {
				case 'n', 't', 'r', '\\', '"', '\'':
					result.WriteByte('\\')
					result.WriteByte(l.CurrentChar)
				default:
					// For any other character, include both \ and the character
					result.WriteByte('\\')
					result.WriteByte(l.CurrentChar)
				}
				continue
			}
		}
		if l.CurrentChar == delimiter {
			break
		}
		result.WriteByte(l.CurrentChar)
	}
	return result.String()
}

func (l *Lexer) readRawString() string {
	var result strings.Builder
	for {
		l.ReadChar()
		if l.CurrentChar == 0 {
			break
		}
		// Handle escaped backticks
		if l.CurrentChar == '\\' {
			nextChar := l.PeekChar()
			if nextChar == '`' {
				l.ReadChar() // consume the backtick
				result.WriteByte('`')
				continue
			}
		}
		if l.CurrentChar == '`' {
			break
		}
		result.WriteByte(l.CurrentChar)
	}
	return result.String()
}

func (l *Lexer) NextToken() token.Token {
	return l.nextToken(l)
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skipLineComment skips characters until the end of line for line comments (//)
func (l *Lexer) skipLineComment() {
	for l.CurrentChar != '\n' && l.CurrentChar != 0 {
		l.ReadChar()
	}
}
