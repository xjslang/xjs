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

This approach ensures that every included feature has demonstrated genuine utility and necessity over the years.

## Extensible Architecture

Everything revolves around the middlewares `UseStatementParser` and `UseExpressionParser`. With these two methods, you can customize the syntax as you wish, adding new features to the language or modifying existing ones.

For convenience, we have also included the methods `RegisterPrefixOperator`, `RegisterInfixOperator`, and `RegisterOperand`, which internally use the middlewares mentioned above.

Additionally, you can concatenate different parsers, further enriching the syntax to suit your preferences. Parsers are executed in LIFO order (Last-In, First-Out).

<details>
	<summary>UseStatementParser example</summary>

```go
// ...
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
// ...
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
