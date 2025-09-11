# XJS API Reference

Complete API documentation for the XJS parser.

## Core Packages

- **[lexer](#lexer-package)** - Tokenization and lexical analysis
- **[parser](#parser-package)** - Syntax parsing and AST building  
- **[ast](#ast-package)** - Abstract Syntax Tree definitions
- **[token](#token-package)** - Token types and structures

---

## Lexer Package

### `lexer.New(input string) *Lexer`

Creates a new lexer instance.

**Parameters:**
- `input` - JavaScript source code to tokenize

**Returns:** `*Lexer` instance

**Example:**
```go
l := lexer.New(`let x = 5`)
```

### `(*Lexer) NextToken() token.Token`

Returns the next token from the input.

**Returns:** `token.Token` with type and literal value

**Example:**
```go
tok := l.NextToken()
fmt.Printf("Type: %v, Literal: %s\n", tok.Type, tok.Literal)
```

### `(*Lexer) Position() (line, column int)`

Returns current position in the source code.

**Returns:**
- `line` - Current line number (1-based)
- `column` - Current column number (1-based)

---

## Parser Package

### `parser.New(lexer *lexer.Lexer) *Parser`

Creates a new parser instance.

**Parameters:**
- `lexer` - Lexer instance to consume tokens from

**Returns:** `*Parser` instance

### `(*Parser) ParseProgram() (*ast.Program, error)`

Parses the entire program and returns the AST.

**Returns:**
- `*ast.Program` - Root node of the AST
- `error` - Parse error if any

**Example:**
```go
program, err := p.ParseProgram()
if err != nil {
    log.Fatal(err)
}
```

### `(*Parser) Errors() []ParserError`

Returns all parsing errors encountered.

**Returns:** Slice of `ParserError` with position and message information

**Example:**
```go
if errors := p.Errors(); len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("Line %d: %s\n", err.Position.Line, err.Message)
    }
}
```

### Extension Methods

#### `(*Parser) UseStatementParser(middleware StatementMiddleware)`

Registers a statement parsing middleware.

**Parameters:**
- `middleware` - Function of type `func(*Parser, func() ast.Statement) ast.Statement`

**Example:**
```go
p.UseStatementParser(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
    if p.CurrentToken.Literal == "const" {
        return parseConstStatement(p)
    }
    return next()
})
```

#### `(*Parser) UseExpressionParser(middleware ExpressionMiddleware)`

Registers an expression parsing middleware.

**Parameters:**
- `middleware` - Function of type `func(*Parser, func() ast.Expression) ast.Expression`

#### `(*Parser) RegisterPrefixOperator(operator string, parser PrefixParser)`

Registers a prefix operator.

**Parameters:**
- `operator` - Operator string (e.g., "typeof", "!")
- `parser` - Function of type `func(func() ast.Expression) ast.Expression`

**Example:**
```go
p.RegisterPrefixOperator("typeof", func(right func() ast.Expression) ast.Expression {
    return &TypeofExpression{Right: right()}
})
```

#### `(*Parser) RegisterInfixOperator(operator string, precedence int, parser InfixParser)`

Registers an infix operator with precedence.

**Parameters:**
- `operator` - Operator string (e.g., "^", "+")
- `precedence` - Operator precedence (use package constants)
- `parser` - Function of type `func(ast.Expression, func() ast.Expression) ast.Expression`

**Example:**
```go
p.RegisterInfixOperator("^", parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
    return &PowerExpression{Left: left, Right: right()}
})
```

#### `(*Parser) RegisterOperand(name string, parser OperandParser)`

Registers a custom operand/literal.

**Parameters:**
- `name` - Operand name (e.g., "PI", "RANDOM")
- `parser` - Function of type `func() ast.Expression`

### Precedence Constants

```go
const (
    LOWEST      = 1
    ASSIGNMENT  = 2  // =, +=, -=
    LOGICAL_OR  = 3  // ||
    LOGICAL_AND = 4  // &&
    EQUALITY    = 5  // ===, !==
    COMPARISON  = 6  // >, <, >=, <=
    SUM         = 7  // +, -
    PRODUCT     = 8  // *, /, %
    UNARY       = 9  // -x, !x, ++x, --x
    POSTFIX     = 10 // x++, x--
    CALL        = 11 // func()
    MEMBER      = 12 // obj.prop, obj[prop]
)
```

### Error Types

#### `ParserError`

```go
type ParserError struct {
    Message  string   `json:"message"`
    Position Position `json:"position"`
    Code     string   `json:"code,omitempty"`
}
```

#### `ParserErrors`

```go
type ParserErrors struct {
    Errors []ParserError `json:"errors"`
    Source string        `json:"source,omitempty"`
}

func (pe ParserErrors) Error() string
func (pe ParserErrors) ToJSON() ([]byte, error)
```

---

## AST Package

### Core Interfaces

#### `Node`
Base interface for all AST nodes.

```go
type Node interface {
    WriteTo(b *strings.Builder)
}
```

#### `Statement`
Interface for statement nodes.

```go
type Statement interface {
    Node
}
```

#### `Expression`
Interface for expression nodes.

```go
type Expression interface {
    Node
}
```

### Program Node

#### `Program`
Root node of the AST.

```go
type Program struct {
    Statements []Statement
}

func (p *Program) String() string
func (p *Program) WriteTo(b *strings.Builder)
```

### Statement Nodes

#### `LetStatement`
Variable declaration.

```go
type LetStatement struct {
    Token token.Token
    Name  *Identifier
    Value Expression
}
```

#### `ReturnStatement`
Return statement.

```go
type ReturnStatement struct {
    Token       token.Token
    ReturnValue Expression
}
```

#### `ExpressionStatement`
Expression used as a statement.

```go
type ExpressionStatement struct {
    Token      token.Token
    Expression Expression
}
```

#### `FunctionDeclaration`
Function declaration.

```go
type FunctionDeclaration struct {
    Token      token.Token
    Name       *Identifier
    Parameters []*Identifier
    Body       *BlockStatement
}
```

#### `IfStatement`
If/else conditional.

```go
type IfStatement struct {
    Token       token.Token
    Condition   Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}
```

#### `BlockStatement`
Block of statements.

```go
type BlockStatement struct {
    Token      token.Token
    Statements []Statement
}
```

#### `ForStatement`
For loop.

```go
type ForStatement struct {
    Token     token.Token
    Init      Statement
    Condition Expression
    Update    Expression
    Body      *BlockStatement
}
```

#### `WhileStatement`
While loop.

```go
type WhileStatement struct {
    Token     token.Token
    Condition Expression
    Body      *BlockStatement
}
```

### Expression Nodes

#### `Identifier`
Variable/function names.

```go
type Identifier struct {
    Token token.Token
    Value string
}
```

#### `NumberLiteral`
Numeric literals.

```go
type NumberLiteral struct {
    Token token.Token
    Value float64
}
```

#### `StringLiteral`
String literals.

```go
type StringLiteral struct {
    Token token.Token
    Value string
}
```

#### `BooleanLiteral`
Boolean literals.

```go
type BooleanLiteral struct {
    Token token.Token
    Value bool
}
```

#### `ArrayLiteral`
Array literals.

```go
type ArrayLiteral struct {
    Token    token.Token
    Elements []Expression
}
```

#### `ObjectLiteral`
Object literals.

```go
type ObjectLiteral struct {
    Token token.Token
    Pairs map[Expression]Expression
}
```

#### `FunctionExpression`
Function expressions.

```go
type FunctionExpression struct {
    Token      token.Token
    Parameters []*Identifier
    Body       *BlockStatement
}
```

#### `CallExpression`
Function calls.

```go
type CallExpression struct {
    Token     token.Token
    Function  Expression
    Arguments []Expression
}
```

#### `IndexExpression`
Array/object indexing.

```go
type IndexExpression struct {
    Token token.Token
    Left  Expression
    Index Expression
}
```

#### `DotExpression`
Property access.

```go
type DotExpression struct {
    Token    token.Token
    Left     Expression
    Property *Identifier
}
```

#### `PrefixExpression`
Prefix operators.

```go
type PrefixExpression struct {
    Token    token.Token
    Operator string
    Right    Expression
}
```

#### `InfixExpression`
Infix operators.

```go
type InfixExpression struct {
    Token    token.Token
    Left     Expression
    Operator string
    Right    Expression
}
```

#### `AssignmentExpression`
Variable assignment.

```go
type AssignmentExpression struct {
    Token    token.Token
    Left     Expression
    Operator string // =, +=, -=
    Right    Expression
}
```

#### `UpdateExpression`
Increment/decrement.

```go
type UpdateExpression struct {
    Token    token.Token
    Operator string // ++, --
    Operand  Expression
    Postfix  bool
}
```

---

## Token Package

### Token Type

```go
type Token struct {
    Type     Type
    Literal  string
    Position Position
}

type Position struct {
    Line   int
    Column int
}
```

### Token Types

#### Special Tokens
```go
const (
    ILLEGAL Type = iota
    EOF
)
```

#### Identifiers and Literals
```go
const (
    IDENT      // variables, functions
    INT        // 123
    FLOAT      // 123.45
    STRING     // "hello"
    RAW_STRING // raw string literal
)
```

#### Operators
```go
const (
    // Assignment
    ASSIGN       // =
    PLUS_ASSIGN  // +=
    MINUS_ASSIGN // -=
    
    // Arithmetic
    PLUS     // +
    MINUS    // -
    MULTIPLY // *
    DIVIDE   // /
    MODULO   // %
    
    // Comparison
    EQ     // ==
    NOT_EQ // !=
    LT     // <
    GT     // >
    LTE    // <=
    GTE    // >=
    
    // Logical
    AND // &&
    OR  // ||
    NOT // !
    
    // Increment/Decrement
    INCREMENT // ++
    DECREMENT // --
)
```

#### Delimiters
```go
const (
    COMMA     // ,
    SEMICOLON // ;
    COLON     // :
    LPAREN    // (
    RPAREN    // )
    LBRACE    // {
    RBRACE    // }
    LBRACKET  // [
    RBRACKET  // ]
    DOT       // .
)
```

#### Keywords
```go
const (
    FUNCTION // function
    LET      // let
    RETURN   // return
    IF       // if
    ELSE     // else
    FOR      // for
    WHILE    // while
    TRUE     // true
    FALSE    // false
    NULL     // null
    UNDEFINED // undefined
)
```

### Helper Functions

#### `LookupIdent(ident string) Type`

Determines if an identifier is a keyword.

**Parameters:**
- `ident` - String to check

**Returns:** `Type` - Either `IDENT` or the specific keyword type

**Example:**
```go
tokenType := token.LookupIdent("function") // Returns token.FUNCTION
tokenType := token.LookupIdent("myVar")    // Returns token.IDENT
```

---

## Usage Examples

### Complete Parsing Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/xjslang/xjs/lexer"
    "github.com/xjslang/xjs/parser"
    "github.com/xjslang/xjs/ast"
)

func main() {
    input := `
        let greeting = "Hello"
        function sayHello(name) {
            return greeting + ", " + name + "!"
        }
        console.log(sayHello("World"))
    `
    
    // Create lexer and parser
    l := lexer.New(input)
    p := parser.New(l)
    
    // Parse the program
    program, err := p.ParseProgram()
    if err != nil {
        log.Fatal(err)
    }
    
    // Check for parsing errors
    if errors := p.Errors(); len(errors) > 0 {
        for _, err := range errors {
            fmt.Printf("Parse error at line %d, column %d: %s\n",
                err.Position.Line, err.Position.Column, err.Message)
        }
        return
    }
    
    // Print the parsed result
    fmt.Println("Parsed successfully!")
    fmt.Println("Output:", program.String())
    
    // Walk through the AST
    for i, stmt := range program.Statements {
        fmt.Printf("Statement %d: ", i+1)
        switch s := stmt.(type) {
        case *ast.LetStatement:
            fmt.Printf("Variable declaration: %s\n", s.Name.Value)
        case *ast.FunctionDeclaration:
            fmt.Printf("Function declaration: %s\n", s.Name.Value)
        case *ast.ExpressionStatement:
            fmt.Println("Expression statement")
        }
    }
}
```

### Custom Extension Example

```go
// Adding a custom operator
func addPowerOperator(p *parser.Parser) {
    p.RegisterInfixOperator("**", parser.PRODUCT+1, 
        func(left ast.Expression, right func() ast.Expression) ast.Expression {
            return &PowerExpression{
                Token: p.CurrentToken,
                Left:  left,
                Right: right(),
            }
        })
}

// Custom AST node
type PowerExpression struct {
    Token token.Token
    Left  ast.Expression
    Right ast.Expression
}

func (pe *PowerExpression) WriteTo(b *strings.Builder) {
    b.WriteString("Math.pow(")
    pe.Left.WriteTo(b)
    b.WriteString(",")
    pe.Right.WriteTo(b)
    b.WriteString(")")
}
```

---

## Error Handling

### Comprehensive Error Handling

```go
func parseWithFullErrorHandling(input string) {
    l := lexer.New(input)
    p := parser.New(l)
    
    program, err := p.ParseProgram()
    
    // Handle fatal parse errors
    if err != nil {
        if parserErrors, ok := err.(parser.ParserErrors); ok {
            // Multiple errors
            fmt.Printf("Found %d parsing errors:\n", len(parserErrors.Errors))
            for _, e := range parserErrors.Errors {
                fmt.Printf("  Line %d, Col %d: %s\n", 
                    e.Position.Line, e.Position.Column, e.Message)
            }
            
            // Get JSON representation
            if jsonData, jsonErr := parserErrors.ToJSON(); jsonErr == nil {
                fmt.Printf("JSON: %s\n", string(jsonData))
            }
        } else {
            // Single error
            fmt.Printf("Parse error: %v\n", err)
        }
        return
    }
    
    // Handle non-fatal parser errors
    if errors := p.Errors(); len(errors) > 0 {
        fmt.Println("Parser warnings:")
        for _, e := range errors {
            fmt.Printf("  %s\n", e.Message)
        }
    }
    
    fmt.Println("Success:", program.String())
}
```

---

This completes the XJS API reference. For more examples and advanced usage, see the [examples directory](examples/).