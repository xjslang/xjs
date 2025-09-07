## Introduction

XJS is a highly configurable JavaScript parser. The idea is to keep the core minimal, excluding redundant structures such as `const`, `var`, or `arrow functions`, and allowing users to add their own structures through the `UseStatementParser` and `UseExpressionParser` methods, which follow a "middleware" design pattern similar to Express.js.

> [!NOTE]  
> As an interesting fact, XJS always interprets the `==` operator as `===`, thus ending the eternal debate between loose equality and strict equality. Otherwise, XJS could be considered a subset of JavaScript.

## Style Guide

**This project has an international scope. Therefore, all source code, comments, and documentation MUST ALWAYS be written in English**, regardless of the language used in conversations, such as Spanish.

**We will NOT use blank lines or comments inside function bodies.** However, blank lines may be used to separate structures, functions, etc.

For example, the following code is incorrect because the comments are not written in English and comments and blank lines are used inside the function body:

```go
// Test ejecuta todos los tests de transpilación (equivalente a test_transpilation.sh)
func Test() error {
	fmt.Println("🚀 Starting XJS Transpilation Tests")
	fmt.Println("====================================")

	// Verificar que Node.js esté disponible
	if err := checkNodeJS(); err != nil {
		return err
	}

	fmt.Println()

	// Ejecutar todos los tests de integración en secuencia
	mg.SerialDeps(TestTranspilation, TestInline, TestErrors)

	fmt.Println()
	fmt.Println("⚡ Running benchmarks...")
	if err := sh.RunV("go", "test", "-run=^$", "-bench=BenchmarkTranspilation", "-benchmem", "./test/integration"); err != nil {
		// Los benchmarks pueden fallar pero no queremos que falle todo el test
		fmt.Println("⚠️  Some benchmarks failed, but continuing...")
	}

	fmt.Println()
	fmt.Println("🎉 All tests completed successfully!")
	return nil
}
```

**We will NOT use redundant comments that do not provide more information than the code itself.** In general, we will avoid writing obvious comments.

For example, the following code is incorrect, since the comment seems redundant:

```go
// Node represents any node in the AST
type Node interface {
	WriteTo(b *strings.Builder)
}
```