# XJS Architecture Overview

Understanding the internal structure and design principles of XJS.

## High-Level Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Source Code   │───▶│      Lexer       │───▶│     Tokens      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
                                                         ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Generated JS   │◀───│      Parser      │◀───│  Token Stream   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                               │
                               ▼
                       ┌─────────────────┐
                       │       AST       │
                       └─────────────────┘
```

## Core Components

### 1. Lexer (Tokenizer)

**Location**: `lexer/lexer.go`  
**Purpose**: Converts source code into tokens

#### Key Features:
- **Position tracking**: Line and column numbers
- **Token identification**: Keywords, operators, literals
- **Error handling**: Invalid characters and malformed tokens

#### Design Patterns:
- **State machine**: Current character → Next character
- **Lookahead**: Peek at next character without consuming
- **Position tracking**: Maintains current line/column

```go
type Lexer struct {
    input        string  // Source code
    position     int     // Current position
    readPosition int     // Next position
    ch           byte    // Current character
    line         int     // Current line
    column       int     // Current column
}
```

### 2. Parser (Syntax Analyzer)

**Location**: `parser/parser.go`  
**Purpose**: Builds AST from token stream

#### Architecture:
- **Pratt Parser**: Operator precedence parsing
- **Recursive Descent**: Statement and expression parsing
- **Middleware System**: Extensible parsing hooks

#### Key Methods:
```go
func (p *Parser) ParseProgram() (*ast.Program, error)
func (p *Parser) parseStatement() ast.Statement
func (p *Parser) parseExpression() ast.Expression
```

### 3. AST (Abstract Syntax Tree)

**Location**: `ast/ast.go`  
**Purpose**: Represents parsed code structure

#### Node Hierarchy:
```
Node (interface)
├── Statement (interface)
│   ├── LetStatement
│   ├── FunctionDeclaration
│   ├── ReturnStatement
│   └── ...
└── Expression (interface)
    ├── Identifier
    ├── NumberLiteral
    ├── InfixExpression
    └── ...
```

#### Code Generation:
Every node implements `WriteTo(b *strings.Builder)` for output generation.

### 4. Token System

**Location**: `token/token.go`  
**Purpose**: Defines all language tokens

#### Token Structure:
```go
type Token struct {
    Type     Type
    Literal  string
    Position Position
}
```

## Parsing Strategy

### Pratt Parser for Expressions

XJS uses Vaughan Pratt's "Top Down Operator Precedence" parsing technique:

1. **Precedence-driven**: Each operator has a precedence level
2. **Left-associative**: Handles operator associativity naturally  
3. **Extensible**: Easy to add new operators

#### Precedence Levels:
```go
const (
    LOWEST      = 1   // Base level
    ASSIGNMENT  = 2   // = += -=
    LOGICAL_OR  = 3   // ||
    LOGICAL_AND = 4   // &&
    EQUALITY    = 5   // === !==
    COMPARISON  = 6   // > < >= <=
    SUM         = 7   // + -
    PRODUCT     = 8   // * / %
    UNARY       = 9   // -x !x ++x --x
    POSTFIX     = 10  // x++ x--
    CALL        = 11  // func()
    MEMBER      = 12  // obj.prop
)
```

### Recursive Descent for Statements

Statements use traditional recursive descent parsing:

```go
func (p *Parser) parseStatement() ast.Statement {
    switch p.CurrentToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    case token.FUNCTION:
        return p.parseFunctionDeclaration()
    default:
        return p.parseExpressionStatement()
    }
}
```

## Middleware System

### Architecture

The middleware system allows plugins to intercept and modify parsing:

```
Input → [Middleware N] → [Middleware N-1] → ... → [Middleware 1] → Default Parser → Output
```

### Types of Middleware

#### 1. Statement Middleware
```go
type StatementMiddleware func(
    parser *Parser, 
    next func() ast.Statement
) ast.Statement
```

#### 2. Expression Middleware
```go
type ExpressionMiddleware func(
    parser *Parser, 
    next func() ast.Expression
) ast.Expression
```

### Execution Order

Middlewares execute in **LIFO** (Last-In, First-Out) order:

```go
p.UseStatementParser(middleware1) // Executes last
p.UseStatementParser(middleware2) // Executes middle
p.UseStatementParser(middleware3) // Executes first
```

This allows:
- **Overriding**: Later middlewares can override earlier ones
- **Composition**: Chain multiple transformations
- **Fallback**: Use `next()` to delegate to previous middleware

## Error Handling Strategy

### Multi-Level Error Handling

1. **Lexer Errors**: Invalid characters, malformed tokens
2. **Parser Errors**: Syntax errors, unexpected tokens
3. **Semantic Errors**: Type mismatches, undefined variables (future)

### Error Collection

XJS collects multiple errors instead of stopping at the first:

```go
type ParserError struct {
    Message  string
    Position Position
    Code     string
}

type ParserErrors struct {
    Errors []ParserError
    Source string
}
```

### Recovery Strategy

When an error occurs:
1. **Record the error**
2. **Skip to next statement boundary**
3. **Continue parsing**
4. **Report all errors at end**

## Memory Management

### AST Node Lifecycle

1. **Creation**: Nodes created during parsing
2. **Linking**: Parent-child relationships established
3. **Transformation**: Middleware can modify nodes
4. **Generation**: `WriteTo()` produces output
5. **Cleanup**: Go garbage collector handles cleanup

### Optimization Opportunities

- **Node pooling**: Reuse common node types
- **Lazy evaluation**: Parse functions only when needed
- **Incremental parsing**: Parse only changed sections

## Extension Points

### 1. Lexer Extensions

Add new token types:

```go
const (
    // Custom tokens
    TEMPLATE_LITERAL  // `template ${expr}`
    REGEX_LITERAL     // /pattern/flags
    BIG_INT          // 123n
)
```

### 2. Parser Extensions

#### Statement Parsers
```go
p.UseStatementParser(func(p *Parser, next func() ast.Statement) ast.Statement {
    // Custom statement parsing logic
    return next()
})
```

#### Expression Parsers
```go
p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
    // Custom expression parsing logic
    return next()
})
```

#### Operator Registration
```go
// Prefix operators (typeof, !, -)
p.RegisterPrefixOperator("typeof", prefixParser)

// Infix operators (+, -, *, etc.)
p.RegisterInfixOperator("**", precedence, infixParser)

// Operands (identifiers, literals)
p.RegisterOperand("PI", operandParser)
```

## Performance Characteristics

### Time Complexity
- **Lexing**: O(n) where n is source code length
- **Parsing**: O(n) for most constructs, O(n²) worst case for deep nesting
- **Code Generation**: O(n) where n is number of AST nodes

### Space Complexity
- **Tokens**: O(n) temporary storage during parsing
- **AST**: O(n) permanent storage until generation
- **Output**: O(n) generated code size

### Optimization Strategies

1. **Single-pass parsing**: No separate semantic analysis phase
2. **Minimal AST**: Only essential information stored
3. **Streaming generation**: Output produced incrementally
4. **Memory pooling**: Reuse token and node objects

## Design Principles

### 1. Minimalism
- **Essential features only**: No redundant syntax
- **Single way to do things**: Reduces cognitive load
- **Clean AST**: Simple node hierarchy

### 2. Extensibility
- **Middleware system**: Powerful extension mechanism
- **Plugin architecture**: Easy to add features
- **Composability**: Plugins work together

### 3. Performance
- **Linear parsing**: O(n) performance
- **Minimal allocations**: Reuse objects where possible
- **Fast compilation**: Quick turnaround during development

### 4. Correctness
- **Comprehensive testing**: High test coverage
- **Error recovery**: Continue parsing after errors
- **Position tracking**: Accurate error reporting

## Future Architecture Enhancements

### 1. Incremental Parsing
- **Change detection**: Parse only modified sections
- **AST caching**: Reuse unchanged subtrees
- **IDE integration**: Real-time syntax checking

### 2. Parallel Processing
- **Function-level parallelism**: Parse functions in parallel
- **Pipeline stages**: Overlap lexing, parsing, and generation
- **Worker pools**: Distribute parsing across cores

### 3. Advanced Optimizations
- **Constant folding**: Evaluate constant expressions at parse time
- **Dead code elimination**: Remove unreachable code
- **Inlining**: Expand simple function calls

### 4. Language Server Protocol
- **LSP implementation**: IDE features like autocomplete
- **Semantic analysis**: Type checking and validation
- **Refactoring support**: Automated code transformations

## Code Organization

```
xjs/
├── lexer/           # Tokenization
│   ├── lexer.go     # Main lexer implementation
│   └── lexer_test.go
├── parser/          # Syntax analysis
│   ├── parser.go    # Core parser
│   ├── parser_functions.go      # Statement/expression parsing
│   ├── parser_middlewares.go    # Extension system
│   ├── parser_registry.go       # Operator registration
│   ├── errors.go    # Error handling
│   └── *_test.go
├── ast/             # Abstract Syntax Tree
│   ├── ast.go       # Node definitions
│   └── ast_test.go
├── token/           # Token system
│   ├── token.go     # Token types and utilities
│   └── token_test.go
├── internal/        # Internal utilities
│   └── debug.go     # Debugging helpers
└── test/            # Integration tests
    ├── integration/ # End-to-end tests
    └── testdata/    # Test fixtures
```

This architecture provides:
- **Separation of concerns**: Each package has a clear responsibility
- **Testability**: Each component can be tested independently
- **Maintainability**: Clean interfaces between components
- **Extensibility**: Plugin system allows external enhancements

---

Understanding this architecture will help you contribute effectively to XJS or build powerful plugins. For implementation details, see the [API Reference](api-reference.md).