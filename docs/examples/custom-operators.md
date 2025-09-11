# Custom Operators Examples

Learn how to add custom operators to XJS with proper precedence and behavior.

## Table of Contents

- [Prefix Operators](#prefix-operators)
- [Infix Operators](#infix-operators)  
- [Postfix Operators](#postfix-operators)
- [Complex Operators](#complex-operators)
- [Operator Precedence](#operator-precedence)

## Prefix Operators

Prefix operators appear before their operand (e.g., `!x`, `-y`, `typeof z`).

### Example: `typeof` Operator

```go
package main

import (
	"fmt"
	"strings"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// TypeofExpression represents typeof x
type TypeofExpression struct {
	Token token.Token
	Right ast.Expression
}

func (te *TypeofExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(typeof ")
	te.Right.WriteTo(b)
	b.WriteString(")")
}

func main() {
	input := `if (typeof x === "undefined") { console.log("x is undefined") }`
	
	l := lexer.New(input)
	p := parser.New(l)
	
	p.RegisterPrefixOperator("typeof", func(right func() ast.Expression) ast.Expression {
		return &TypeofExpression{
			Token: p.CurrentToken,
			Right: right(),
		}
	})
	
	program, _ := p.ParseProgram()
	fmt.Println(program.String())
	// Output: if (((typeof x)==="undefined")){console.log("x is undefined")}
}
```

### Example: `delete` Operator

```go
type DeleteExpression struct {
	Token token.Token
	Right ast.Expression
}

func (de *DeleteExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(delete ")
	de.Right.WriteTo(b)
	b.WriteString(")")
}

p.RegisterPrefixOperator("delete", func(right func() ast.Expression) ast.Expression {
	return &DeleteExpression{
		Token: p.CurrentToken,
		Right: right(),
	}
})

// Usage: delete obj.prop
// Output: (delete obj.prop)
```

### Example: `void` Operator

```go
type VoidExpression struct {
	Token token.Token
	Right ast.Expression
}

func (ve *VoidExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(void ")
	ve.Right.WriteTo(b)
	b.WriteString(")")
}

p.RegisterPrefixOperator("void", func(right func() ast.Expression) ast.Expression {
	return &VoidExpression{
		Token: p.CurrentToken,
		Right: right(),
	}
})

// Usage: void 0
// Output: (void 0)
```

## Infix Operators

Infix operators appear between two operands (e.g., `x + y`, `a ** b`).

### Example: Power Operator (`**`)

```go
package main

import (
	"fmt"
	"strings"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

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

func main() {
	input := `let result = base ** exponent`
	
	l := lexer.New(input)
	p := parser.New(l)
	
	// Right associative power operator (higher precedence than multiplication)
	p.RegisterInfixOperator("**", parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowerExpression{
			Token: p.CurrentToken,
			Left:  left,
			Right: right(),
		}
	})
	
	program, _ := p.ParseProgram()
	fmt.Println(program.String())
	// Output: let result=Math.pow(base,exponent)
}
```

### Example: Nullish Coalescing (`??`)

```go
type NullishCoalescingExpression struct {
	Token token.Token
	Left  ast.Expression
	Right ast.Expression
}

func (nce *NullishCoalescingExpression) WriteTo(b *strings.Builder) {
	b.WriteString("((")
	nce.Left.WriteTo(b)
	b.WriteString("!==null&&")
	nce.Left.WriteTo(b)
	b.WriteString("!==undefined)?")
	nce.Left.WriteTo(b)
	b.WriteString(":")
	nce.Right.WriteTo(b)
	b.WriteString(")")
}

p.RegisterInfixOperator("??", parser.LOGICAL_OR, func(left ast.Expression, right func() ast.Expression) ast.Expression {
	return &NullishCoalescingExpression{
		Token: p.CurrentToken,
		Left:  left,
		Right: right(),
	}
})

// Usage: value ?? defaultValue
// Output: ((value!==null&&value!==undefined)?value:defaultValue)
```

### Example: Pipeline Operator (`|>`)

```go
type PipelineExpression struct {
	Token token.Token
	Left  ast.Expression
	Right ast.Expression
}

func (pe *PipelineExpression) WriteTo(b *strings.Builder) {
	// Transform: value |> func becomes func(value)
	pe.Right.WriteTo(b)
	b.WriteString("(")
	pe.Left.WriteTo(b)
	b.WriteString(")")
}

p.RegisterInfixOperator("|>", parser.LOWEST+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
	return &PipelineExpression{
		Token: p.CurrentToken,
		Left:  left,
		Right: right(),
	}
})

// Usage: value |> transform |> display
// Output: display(transform(value))
```

### Example: Range Operator (`..`)

```go
type RangeExpression struct {
	Token token.Token
	Start ast.Expression
	End   ast.Expression
}

func (re *RangeExpression) WriteTo(b *strings.Builder) {
	b.WriteString("Array.from({length:(")
	re.End.WriteTo(b)
	b.WriteString("-")
	re.Start.WriteTo(b)
	b.WriteString("+1)},(_,i)=>(")
	re.Start.WriteTo(b)
	b.WriteString("+i))")
}

p.RegisterInfixOperator("..", parser.COMPARISON-1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
	return &RangeExpression{
		Token: p.CurrentToken,
		Start: left,
		End:   right(),
	}
})

// Usage: 1..10
// Output: Array.from({length:(10-1+1)},(_,i)=>(1+i))
```

## Postfix Operators

Postfix operators are typically handled in expression parsing, but here's how you might implement them:

### Example: Factorial Operator (`!`)

```go
type FactorialExpression struct {
	Token   token.Token
	Operand ast.Expression
}

func (fe *FactorialExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(function(n){let r=1;for(let i=2;i<=n;i++)r*=i;return r})(")
	fe.Operand.WriteTo(b)
	b.WriteString(")")
}

// This would require modifying the expression parser to handle postfix operators
p.UseExpressionParser(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
	expr := next()
	
	// Check if next token is our postfix operator
	if p.CurrentToken.Literal == "!" {
		return &FactorialExpression{
			Token:   p.CurrentToken,
			Operand: expr,
		}
	}
	
	return expr
})
```

## Complex Operators

### Example: Ternary Operator (`? :`)

```go
type TernaryExpression struct {
	Token     token.Token
	Condition ast.Expression
	TrueExpr  ast.Expression
	FalseExpr ast.Expression
}

func (te *TernaryExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(")
	te.Condition.WriteTo(b)
	b.WriteString("?")
	te.TrueExpr.WriteTo(b)
	b.WriteString(":")
	te.FalseExpr.WriteTo(b)
	b.WriteString(")")
}

p.UseExpressionParser(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
	expr := next()
	
	// Check for ternary operator
	if p.CurrentToken.Type == token.QUESTION {
		ternary := &TernaryExpression{
			Token:     p.CurrentToken,
			Condition: expr,
		}
		
		p.NextToken() // skip ?
		ternary.TrueExpr = p.ParseExpression()
		
		if !p.ExpectToken(token.COLON) {
			return nil
		}
		
		p.NextToken() // skip :
		ternary.FalseExpr = p.ParseExpression()
		
		return ternary
	}
	
	return expr
})
```

### Example: Optional Chaining (`?.`)

```go
type OptionalChainingExpression struct {
	Token    token.Token
	Left     ast.Expression
	Property ast.Expression
}

func (oce *OptionalChainingExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(")
	oce.Left.WriteTo(b)
	b.WriteString("==null?undefined:")
	oce.Left.WriteTo(b)
	b.WriteString(".")
	oce.Property.WriteTo(b)
	b.WriteString(")")
}

p.RegisterInfixOperator("?.", parser.MEMBER, func(left ast.Expression, right func() ast.Expression) ast.Expression {
	return &OptionalChainingExpression{
		Token:    p.CurrentToken,
		Left:     left,
		Property: right(),
	}
})

// Usage: obj?.prop
// Output: (obj==null?undefined:obj.prop)
```

## Operator Precedence

Understanding precedence is crucial for correct parsing:

### Precedence Levels (Low to High)

```go
const (
    LOWEST      = 1   // Lowest precedence
    ASSIGNMENT  = 2   // = += -= *= /= %=
    TERNARY     = 3   // ? :
    LOGICAL_OR  = 4   // ||
    LOGICAL_AND = 5   // &&
    BITWISE_OR  = 6   // |
    BITWISE_XOR = 7   // ^
    BITWISE_AND = 8   // &
    EQUALITY    = 9   // === !== == !=
    COMPARISON  = 10  // < > <= >= instanceof in
    SHIFT       = 11  // << >> >>>
    SUM         = 12  // + -
    PRODUCT     = 13  // * / %
    POWER       = 14  // **
    UNARY       = 15  // ! ~ + - typeof void delete
    POSTFIX     = 16  // ++ --
    CALL        = 17  // () []
    MEMBER      = 18  // . ?.
)
```

### Precedence Examples

```go
// Wrong precedence (addition has higher precedence than multiplication)
p.RegisterInfixOperator("myop", parser.PRODUCT+1, ...)

// Correct precedence (between multiplication and power)
p.RegisterInfixOperator("myop", parser.PRODUCT+1, ...)

// Same precedence as multiplication
p.RegisterInfixOperator("myop", parser.PRODUCT, ...)
```

### Associativity

Most operators are left-associative by default. For right-associative operators (like `**`), you need special handling:

```go
// Right-associative power operator
p.RegisterInfixOperator("**", parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
    // The parser automatically handles right-associativity for operators
    // with the same precedence level
    return &PowerExpression{Left: left, Right: right()}
})

// For explicit right-associativity, you might need custom logic in the parser
```

## Complete Example: Mathematical Extensions

```go
package main

import (
	"fmt"
	"strings"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// Power operator
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

// Modular exponentiation
type ModPowExpression struct {
	Token    token.Token
	Base     ast.Expression
	Exponent ast.Expression
	Modulus  ast.Expression
}

func (mpe *ModPowExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(function(b,e,m){let r=1;b=b%m;while(e>0){if(e%2===1)r=(r*b)%m;e=Math.floor(e/2);b=(b*b)%m}return r})(")
	mpe.Base.WriteTo(b)
	b.WriteString(",")
	mpe.Exponent.WriteTo(b)
	b.WriteString(",")
	mpe.Modulus.WriteTo(b)
	b.WriteString(")")
}

// Mathematical constants
type MathConstant struct {
	Token token.Token
	Name  string
}

func (mc *MathConstant) WriteTo(b *strings.Builder) {
	switch mc.Name {
	case "PI":
		b.WriteString("Math.PI")
	case "E":
		b.WriteString("Math.E")
	case "TAU":
		b.WriteString("(2*Math.PI)")
	case "PHI":
		b.WriteString("((1+Math.sqrt(5))/2)")
	}
}

func setupMathExtensions(p *parser.Parser) {
	// Power operator
	p.RegisterInfixOperator("**", parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowerExpression{Left: left, Right: right()}
	})
	
	// Mathematical constants
	constants := []string{"PI", "E", "TAU", "PHI"}
	for _, name := range constants {
		p.RegisterOperand(name, func() ast.Expression {
			return &MathConstant{Token: p.CurrentToken, Name: p.CurrentToken.Literal}
		})
	}
	
	// Modular exponentiation: base^exp mod m
	p.UseExpressionParser(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		expr := next()
		
		// Look for pattern: expr^exp mod m
		if p.CurrentToken.Literal == "mod" {
			if powerExpr, ok := expr.(*PowerExpression); ok {
				modExpr := &ModPowExpression{
					Token:    p.CurrentToken,
					Base:     powerExpr.Left,
					Exponent: powerExpr.Right,
				}
				p.NextToken()
				modExpr.Modulus = p.ParseExpression()
				return modExpr
			}
		}
		
		return expr
	})
}

func main() {
	input := `
		let circle_area = PI * radius**2
		let golden_ratio = PHI
		let crypto_result = base**exp mod prime
	`
	
	l := lexer.New(input)
	p := parser.New(l)
	
	setupMathExtensions(p)
	
	program, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	
	fmt.Println(program.String())
}
```

## Best Practices for Custom Operators

### 1. Choose Appropriate Precedence

```go
// Consider how your operator should interact with others
// Higher number = higher precedence (binds tighter)

p.RegisterInfixOperator("custom", parser.PRODUCT, ...) // Same as * /
p.RegisterInfixOperator("custom", parser.PRODUCT+1, ...) // Higher than * /
p.RegisterInfixOperator("custom", parser.SUM-1, ...) // Lower than + -
```

### 2. Handle Edge Cases

```go
p.RegisterPrefixOperator("myop", func(right func() ast.Expression) ast.Expression {
	rightExpr := right()
	if rightExpr == nil {
		// Handle parse error
		return nil
	}
	
	return &MyOperatorExpression{
		Token: p.CurrentToken,
		Right: rightExpr,
	}
})
```

### 3. Provide Meaningful Output

```go
func (expr *MyExpression) WriteTo(b *strings.Builder) {
	// Generate readable, valid JavaScript
	b.WriteString("/* MyOperator */ (")
	expr.Left.WriteTo(b)
	b.WriteString(" custom_op ")
	expr.Right.WriteTo(b)
	b.WriteString(")")
}
```

### 4. Test Thoroughly

```go
func TestCustomOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a ** b", "Math.pow(a,b)"},
		{"2 ** 3 * 4", "(Math.pow(2,3)*4)"}, // Test precedence
		{"(2 + 3) ** 4", "Math.pow((2+3),4)"}, // Test grouping
	}
	
	// Run tests...
}
```

---

These examples show how to extend XJS with powerful custom operators. Remember to consider precedence, associativity, and edge cases when implementing your operators.