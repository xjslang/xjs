# AST Package

## Design Decision: AST and Code Generation Coupling

For practicality, the Abstract Syntax Tree (AST) and the JavaScript code generation are **intentionally coupled** in the same package.

### Why This Approach?

XJS is designed to transpile to **JavaScript only**. While separating the code generation into its own package would be architecturally "pure", it would introduce significant overhead with no practical benefit:

**Current approach (coupled):**
```go
// ast/ast.go
type BinaryExpression struct {
    Token    token.Token
    Left     Expression
    Operator string
    Right    Expression
}

func (be *BinaryExpression) WriteTo(cw *CodeWriter) {
    cw.WriteRune('(')
    be.Left.WriteTo(cw)
    cw.WriteString(be.Operator)
    be.Right.WriteTo(cw)
    cw.WriteRune(')')
}
```

**Decoupled approach (hypothetical):**
```go
// ast/ast.go
type BinaryExpression struct {
    Token    token.Token
    Left     Expression
    Operator string
    Right    Expression
}

// compiler/js.go
type JsNode interface {
    WriteTo(cw *CodeWriter)
}

type JsBinaryExpression struct {
    Left     JsNode
    Operator string
    Right    JsNode
}

func NewJsBinaryExpression(be *ast.BinaryExpression) *JsBinaryExpression {
    return &JsBinaryExpression{
        Left:     convertExpression(be.Left),
        Operator: be.Operator,
        Right:    convertExpression(be.Right),
    }
}

func (jbe *JsBinaryExpression) WriteTo(cw *CodeWriter) {
    // Same implementation as above
}
```

**Cost of decoupling:**
- Duplicate all 20+ AST node types in the `compiler` package
- Write conversion functions for each type (`NewJsProgram`, `NewJsBinaryExpression`, etc.)
- Maintain parallel type hierarchies
- Add complexity with no practical return on investment

### Files in This Package

- **`ast.go`** - All AST node definitions and their `WriteTo` methods
- **`code_writer.go`** - Helper for generating JavaScript code with source maps

Since we only target JavaScript, keeping the code generation logic (`WriteTo` methods) directly on AST nodes follows the **YAGNI principle** (You Aren't Gonna Need It) and keeps the codebase maintainable.

If in the future we need to support multiple target languages (Python, C, etc.), the architecture can be refactored. Until then, simplicity wins.
