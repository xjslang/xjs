// Package lexer provides lexical analysis functionality for the XJS language.
// It tokenizes source code into a sequence of tokens that can be consumed by the parser.
package lexer

import (
	"strings"

	"github.com/xjslang/xjs/token"
)

// Interceptor allows middleware-style interception of token generation.
type Interceptor func(l *Lexer, next func() token.Token) token.Token

// Builder provides a fluent interface for constructing lexers with interceptors and dynamic tokens.
type Builder struct {
	interceptors  []Interceptor
	dynamicTokens map[string]token.Type
	nextTokenID   token.Type
}

// Lexer performs lexical analysis on XJS source code, converting it into a stream of tokens.
type Lexer struct {
	input            string
	position         int  // current position in input (points to current char)
	readPosition     int  // current reading position in input (after current char)
	CurrentChar      byte // current char under examination
	Line             int  // current line
	Column           int  // current column
	nextToken        func(*Lexer) token.Token
	hadNewlineBefore bool     // tracks if we just consumed a newline
	leadingComments  []string // leading comments before the token
}

func newWithOptions(input string, interceptors ...Interceptor) *Lexer {
	l := &Lexer{
		input:     input,
		Column:    -1,
		nextToken: baseNextToken,
	}
	for _, reader := range interceptors {
		l.useTokenInterceptor(reader)
	}
	l.ReadChar()
	return l
}

// new creates a new Lexer instance for the given input string.
func new(input string) *Lexer {
	return newWithOptions(input)
}

// ReadChar advances the lexer to the next character in the input.
func (l *Lexer) ReadChar() {
	// If the previous character was a newline, reset column
	if l.CurrentChar == '\n' {
		l.Line++
		l.Column = -1 // Will become 0 after increment below
	}

	// Increment column for the new character position
	l.Column++

	if l.readPosition >= len(l.input) {
		l.CurrentChar = 0 // ASCII NUL character represents "EOF"
	} else {
		l.CurrentChar = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// PeekChar returns the next character without advancing the lexer position.
func (l *Lexer) PeekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readLeadingComments reads leading comments (it includes collapsed blank lines)
func (l *Lexer) readLeadingComments() {
	l.hadNewlineBefore = false
	l.leadingComments = nil
	for {
		// read whitespaces
		for isWhitespace(l.CurrentChar) {
			if l.CurrentChar == '\n' {
				l.hadNewlineBefore = true
				l.leadingComments = append(l.leadingComments, "")
			}
			l.ReadChar()
		}

		// read comments
		for l.CurrentChar == '/' && l.PeekChar() == '/' {
			// consume "//"
			l.ReadChar()
			l.ReadChar()

			var comment strings.Builder
			for l.CurrentChar != '\n' && l.CurrentChar != 0 {
				comment.WriteByte(l.CurrentChar)
				l.ReadChar()
			}
			// omits the last newline
			if l.CurrentChar == '\n' {
				l.hadNewlineBefore = true
				l.ReadChar()
			}
			l.leadingComments = append(l.leadingComments, comment.String())
		}

		if !isWhitespace(l.CurrentChar) {
			break
		}
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

// readNumber reads a number (integer, decimal, scientific notation, hexadecimal, binary, or octal)
func (l *Lexer) readNumber() (string, token.Type) {
	position := l.position
	tokenType := token.INT

	// Check for hexadecimal numbers (0x or 0X)
	if l.CurrentChar == '0' && (l.PeekChar() == 'x' || l.PeekChar() == 'X') {
		return l.readHexNumber()
	}

	// Check for binary numbers (0b or 0B)
	if l.CurrentChar == '0' && (l.PeekChar() == 'b' || l.PeekChar() == 'B') {
		return l.readBinaryNumber()
	}

	// Check for octal numbers (0o or 0O)
	if l.CurrentChar == '0' && (l.PeekChar() == 'o' || l.PeekChar() == 'O') {
		return l.readOctalNumber()
	}

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

	// Check for scientific notation (e or E)
	if l.CurrentChar == 'e' || l.CurrentChar == 'E' {
		tokenType = token.FLOAT
		l.ReadChar() // consume the 'e' or 'E'

		// Check for optional sign
		if l.CurrentChar == '+' || l.CurrentChar == '-' {
			l.ReadChar() // consume the sign
		}

		// Read the exponent digits
		if !isDigit(l.CurrentChar) {
			// Invalid scientific notation - return what we have so far
			return l.input[position:l.position], tokenType
		}

		for isDigit(l.CurrentChar) {
			l.ReadChar()
		}
	}

	return l.input[position:l.position], tokenType
}

// readHexNumber reads a hexadecimal number (0x or 0X followed by hex digits)
func (l *Lexer) readHexNumber() (string, token.Type) {
	position := l.position

	// Consume '0'
	l.ReadChar()
	// Consume 'x' or 'X'
	l.ReadChar()

	// Read hex digits
	if !isHexDigit(l.CurrentChar) {
		// Invalid hex number - return what we have so far
		return l.input[position:l.position], token.INT
	}

	for isHexDigit(l.CurrentChar) {
		l.ReadChar()
	}

	return l.input[position:l.position], token.INT
}

// readBinaryNumber reads a binary number (0b or 0B followed by binary digits)
func (l *Lexer) readBinaryNumber() (string, token.Type) {
	position := l.position

	// Consume '0'
	l.ReadChar()
	// Consume 'b' or 'B'
	l.ReadChar()

	// Read binary digits
	if !isBinaryDigit(l.CurrentChar) {
		// Invalid binary number - return what we have so far
		return l.input[position:l.position], token.INT
	}

	for isBinaryDigit(l.CurrentChar) {
		l.ReadChar()
	}

	return l.input[position:l.position], token.INT
}

// readOctalNumber reads an octal number (0o or 0O followed by octal digits)
func (l *Lexer) readOctalNumber() (string, token.Type) {
	position := l.position

	// Consume '0'
	l.ReadChar()
	// Consume 'o' or 'O'
	l.ReadChar()

	// Read octal digits
	if !isOctalDigit(l.CurrentChar) {
		// Invalid octal number - return what we have so far
		return l.input[position:l.position], token.INT
	}

	for isOctalDigit(l.CurrentChar) {
		l.ReadChar()
	}

	return l.input[position:l.position], token.INT
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

// NextToken generates and returns the next token from the input stream.
func (l *Lexer) NextToken() token.Token {
	l.readLeadingComments()
	return l.nextToken(l)
}

func (l *Lexer) useTokenInterceptor(interceptor Interceptor) {
	next := l.nextToken
	l.nextToken = func(l *Lexer) token.Token {
		return interceptor(l, func() token.Token {
			return next(l)
		})
	}
}
