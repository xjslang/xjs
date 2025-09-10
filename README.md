# XJS (eXtensible JavaScript parser)

**XJS** is a highly customizable JavaScript parser. Our goal is to create a JavaScript compiler that includes only the essential, proven features while enabling users to extend the language through dynamic plugins.

## Installation

```bash
go get github.com/xjslang/xjs@latest
```

## Minimalism and Sufficiency

Rather than accumulating features over time, **XJS** starts with a carefully curated set of **necessary and sufficient** language constructs. We have deliberately excluded redundant features:

- **No classes** - Functions provide sufficient abstraction capabilities
- **No arrow functions** - Regular function syntax is adequate
- **No `const/var`** - A single variable declaration mechanism suffices
- **No `try/catch`** - Alternative error handling patterns are preferred
- **No redundant syntactic sugar** - Focus on core functionality

This approach ensures that every included feature has demonstrated genuine utility and necessity over the years.

## Extensible Architecture

Todo gira en torno a los middlewares `UseStatementParser` y `UseExpressionParser`. Mediante estos dos métodos podemos personalizar la sintaxis a nuestro gusto, añadiendo nuevas características al lenguage o modificando características existentes.

No obstante, y por conveniencia, hemos añadido los métodos `RegisterPrefixOperator`, `RegisterInfixOperator` y `RegisterOperand`, que internamente utilizan los middlewares anteriores.

Además, podemos concatenar diferentes parsers, enriqueciendo de esta forma la sintaxis a nuestras preferencias. Los parsers se ejecutan siguiendo el orden LIFO (Last-In, First-Out).

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
