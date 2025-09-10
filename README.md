# XJS (eXtensible JavaScript parser)

**XJS** is a highly customizable JavaScript parser. Our goal is to create a JavaScript compiler that includes only the essential, proven features while enabling users to extend the language through dynamic plugins.

## Installation

```bash
go get github.com/xjslang/xjs@latest
```

## Minimalism and Sufficiency

Rather than accumulating features over time, **XJS** starts with a carefully curated set of **necessary and sufficient** language constructs. We have deliberately excluded redundant and confusing features:

- **No classes** - Functions provide sufficient abstraction capabilities
- **No arrow functions** - Regular function syntax is adequate
- **No `const/var`** - A single variable declaration mechanism suffices
- **No `try/catch`** - Alternative error handling patterns are preferred
- **No weak equality** - The `==/!=` operators are automatically translated to `===/!==`
- **No redundant syntactic sugar** - Focus on core functionality

This approach ensures that every included feature has demonstrated genuine utility and necessity over the years. However, you can always create a plugin to implement any of the discarded features!

## Extensible Architecture

Everything revolves around the middlewares `UseStatementParser` and `UseExpressionParser`. With these two methods, you can customize the syntax as you wish, adding new features to the language or modifying existing ones.

For convenience, we have also included the methods `RegisterPrefixOperator`, `RegisterInfixOperator`, and `RegisterOperand`, which internally use the middlewares mentioned above.

Additionally, you can concatenate different parsers, further enriching the syntax to suit your preferences. Parsers are executed in LIFO order (Last-In, First-Out).

<details>
	<summary>UseStatementParser example</summary>

```go
import (
	"fmt"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

// Represents a `const` node
type ConstStatement struct {
	Token token.Token
	Name  *ast.Identifier
	Value ast.Expression
}

// Tells the parser how to write a node
func (ls *ConstStatement) WriteTo(b *strings.Builder) {
	b.WriteString("const ")
	ls.Name.WriteTo(b)
	if ls.Value != nil {
		b.WriteRune('=')
		ls.Value.WriteTo(b)
	}
}

func main() {
	input := "const x = 42"
	l := lexer.New(input)
	p := New(l)
	// adds support for the `const` keyword!
	p.UseStatementParser(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "const" {
			stmt := &ConstStatement{Token: p.CurrentToken}
			p.NextToken() // moves to identifier token
			stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
			if !p.ExpectToken(token.ASSIGN) { // expects "="
				return nil
			}
			p.NextToken() // moves to value expression
			stmt.Value = p.ParseExpression()
			return stmt
		}
		return next() // otherwise, next!
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: const x=42
}
```
</details>

<details>
	<summary>UseExpressionParser example</summary>

```go
// ...
```
</details>

<details>
	<summary>RegisterPrefixOperator example</summary>

```go
// ...
```
</details>

<details>
	<summary>RegisterInfixOperator example</summary>

```go
// ...
```
</details>

<details>
	<summary>RegisterOperand example</summary>

```go
import (
	"fmt"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
)

func main() {
	input := "let area = PI * r * r"
	l := lexer.New(input)
	p := New(l)
	// adds support for the PI constant!
	p.RegisterOperand("PI", func() ast.Expression {
		return &PiLiteral{Token: p.CurrentToken}
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let area=((Math.PI*r)*r)
}
```
</details>

<details>
	<summary>Concatenate multiple parsers</summary>

```go
// ...
```
</details>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
