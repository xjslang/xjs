# XJS Transpilation Test Data

This directory contains test fixtures for the XJS transpiler. Following Go conventions, files in `testdata/` are automatically excluded from builds and are only used during testing.

## Testing Approach

Instead of comparing transpiled JavaScript strings directly, our tests:

1. **Transpile** XJS code to JavaScript using the parser
2. **Execute** the transpiled JavaScript using Node.js
3. **Compare** the execution output with expected results

This approach ensures that:
- The transpiled code is syntactically correct
- The transpiled code produces the expected behavior
- Tests are independent of formatting and stylistic changes
- We test the actual functionality, not just string matching

## Test Fixtures

### Core Test Files

- `transpilation_test.go` - Main transpilation tests using fixture files
- `middleware_test.go` - Tests for custom middleware functionality
- `test_transpilation.sh` - Test runner script

### Test Fixtures

Located in `testdata/`:

- `basic.js` / `basic.output` - Simple console.log test
- `function.js` / `function.output` - Function declaration and execution
- `array_loop.js` / `array_loop.output` - Array iteration with loops
- `conditional.js` / `conditional.output` - If/else conditionals
- `object.js` / `object.output` - Object creation and property access

Each test case consists of:
- `.js` file - The XJS source code to transpile
- `.output` file - The expected output when executed

## Running Tests

### Prerequisites

- Go 1.21+ installed
- Node.js installed (for executing transpiled JavaScript)

### Commands

```bash
# Run all transpilation tests
go test -v ./...

# Run specific test groups
go test -v -run TestTranspilation
go test -v -run TestMiddlewareHandlers
go test -v -run TestComplexTranspilation

# Run with the convenience script
./test_transpilation.sh

# Run benchmarks
go test -bench=BenchmarkTranspilation -benchmem
```

## Test Categories

### 1. Fixture-Based Tests (`TestTranspilation`)

Tests using predefined input/output file pairs:
- Loads XJS code from `.js` fixtures
- Compares execution output with `.out` files
- Covers common JavaScript patterns

### 2. Inline Tests (`TestTranspilationBasicInline`)

Programmatic tests with inline test data:
- Simple expressions and statements
- Variable declarations
- Basic arithmetic
- Function calls

### 3. Middleware Tests (`TestMiddlewareHandlers`)

Tests for XJS's custom middleware system:
- Expression handlers
- Statement handlers
- Multiple middleware chains
- Custom language features

### 4. Complex Transpilation Tests (`TestComplexTranspilation`)

Advanced scenarios:
- Nested functions
- Array methods
- String operations
- Boolean logic
- Type checking

### 5. Error Handling Tests (`TestTranspilationErrors`)

Negative test cases:
- Invalid syntax
- Malformed input
- Parser error handling

### 6. Performance Tests (`BenchmarkTranspilation`, `TestPerformanceRegression`)

Performance validation:
- Transpilation speed benchmarks
- Large input handling
- Memory usage verification

## Adding New Tests

### Adding Fixture Tests

1. Create a new `.js` file in `testdata/` with XJS code
2. Create a corresponding `.out` file with expected output
3. Add the base filename (without extension) to the `testCases` slice in `transpilation_test.go`

Example:
```javascript
// testdata/my_feature.js
let x = 5;
console.log(x * 2);
```

```
// testdata/my_feature.out
10
```

### Adding Inline Tests

Add a new test case to the `tests` slice in `TestTranspilationBasicInline`:

```go
{
    name:           "my_test_case",
    inputFile:      `console.log('Hello, XJS!')`,
    expectedOutput: "Hello, XJS!",
},
```

### Adding Middleware Tests

Create new test functions in `middleware_test.go` that use custom expression or statement handlers.

## Best Practices

1. **Test Output Determinism**: Ensure test outputs are deterministic and don't depend on system time, random values, etc.

2. **Error Messages**: Include the transpiled JavaScript in error messages for debugging

3. **Node.js Availability**: Tests automatically skip if Node.js is not available

4. **Isolated Tests**: Each test should be independent and not rely on state from other tests

5. **Comprehensive Coverage**: Test both success and failure cases

## Troubleshooting

### Common Issues

1. **Node.js Not Found**: Install Node.js or ensure it's in your PATH
2. **Test Output Mismatch**: Check for extra whitespace, newlines, or formatting differences
3. **Transpilation Errors**: Verify the XJS syntax is correct

### Debugging

To debug a failing test:

1. Run with verbose output: `go test -v`
2. Check the transpiled JavaScript in the error message
3. Manually run the transpiled JS with Node.js to see the actual output
4. Compare character-by-character for subtle differences

## Integration with CI/CD

The test suite is designed to work in CI/CD environments:

- Automatic Node.js detection and graceful skipping
- Clear exit codes for success/failure
- Comprehensive error reporting
- Performance regression detection

Example GitHub Actions integration:

```yaml
- name: Setup Node.js
  uses: actions/setup-node@v3
  with:
    node-version: '18'

- name: Run Transpilation Tests
  run: |
    cd xjs
    go test -v ./...
```
