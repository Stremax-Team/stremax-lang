# Stremax-Lang Testing

This document summarizes the testing of the Stremax-Lang interpreter.

## Continuous Integration

Stremax-Lang uses GitHub Actions for continuous integration. The CI pipeline includes:

1. **Building**: Compiles the project and uploads the binary as an artifact
2. **Testing**: Runs unit tests with race detection and coverage analysis
3. **Linting**: Ensures code quality with gofmt, golint, and go vet
4. **Integration Testing**: Runs example programs to verify end-to-end functionality
5. **Benchmarking**: Measures performance and uploads results as artifacts

You can view the CI configuration in `.github/workflows/ci.yml`.

## Running Tests Locally

The easiest way to run all tests locally is to use the provided script:

```bash
./run_tests.sh
```

This script will:
1. Check your Go version and adjust commands accordingly
2. Check and fix code formatting
3. Run linting and static analysis
4. Run unit tests with race detection and coverage
5. Run benchmarks
6. Build the project
7. Run example programs

### Unit Tests

Run all unit tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View coverage in browser
```

### Benchmarks

Run benchmarks:

```bash
go test -bench=. -benchmem ./...
```

Run specific benchmarks:

```bash
go test -bench=BenchmarkLexer -benchmem ./pkg/lexer
```

### Linting

Format code:

```bash
gofmt -w .
```

Run linter:

```bash
golint ./...
# or if golint is installed via go install
$(go env GOPATH)/bin/golint ./...
```

Run Go's built-in code analyzer:

```bash
go vet ./...
```

## Interpreter Tests

The interpreter has been tested with a suite of unit tests in `pkg/interpreter/interpreter_test.go`. These tests cover:

- Integer expressions and arithmetic operations
- Boolean expressions and comparison operations
- Conditional statements (if-else)
- Variable declarations and references
- String literals and concatenation
- Logical operators (AND, OR) with short-circuit evaluation
- Block scoping for variables
- Error handling

All tests are passing, which indicates that the core functionality of the interpreter is working correctly.

## Benchmark Tests

We've added benchmark tests for key components:

- **Lexer Benchmarks** (`pkg/lexer/lexer_benchmark_test.go`): Tests the performance of tokenization
- **Parser Benchmarks** (`pkg/parser/parser_benchmark_test.go`): Tests the performance of parsing
- **Interpreter Benchmarks** (`pkg/interpreter/interpreter_benchmark_test.go`): Tests the performance of code execution

These benchmarks help us track performance over time and detect regressions.

## Example Programs

We've created several example programs to demonstrate the features of Stremax-Lang:

### Basic Examples (Working)

- **simple.sx**: A very basic example with variable assignment and arithmetic.
- **arithmetic.sx**: Demonstrates various arithmetic operations.
- **strings.sx**: Shows string operations like concatenation.
- **conditionals.sx**: Illustrates conditional expressions with if-else statements.
- **boolean.sx**: Tests boolean operations and comparisons.
- **logical.sx**: Demonstrates logical operators (&&, ||) with short-circuit evaluation.
- **simple_logical.sx**: A minimal example of logical operators.
- **combined.sx**: Combines various features in one example.
- **errors.sx**: Demonstrates error handling (division by zero).
- **scoping.sx**: Tests variable scoping.
- **scoping_error.sx**: Tests variable scoping errors.
- **debug.sx**: Tests interpreter output with multiple expressions.

### Advanced Examples (Not Working Yet)

- **token.sx**: An ERC20-like token contract implementation.
- **voting.sx**: A voting contract implementation.
- **auction.sx**: An auction contract implementation.

## Test Results

| Example | Status | Result |
|---------|--------|--------|
| simple.sx | ✅ | Returns 15 |
| arithmetic.sx | ✅ | Returns 37 (last expression result) |
| strings.sx | ✅ | Returns "John Doe" |
| conditionals.sx | ✅ | Returns 50 |
| boolean.sx | ✅ | Returns 1 (true) |
| logical.sx | ✅ | Returns true (last expression result) |
| simple_logical.sx | ✅ | Returns false (a && b where a=true, b=false) |
| combined.sx | ✅ | Returns 100 |
| errors.sx | ✅ | Throws "Division by zero" error |
| scoping.sx | ✅ | Returns 30 |
| scoping_error.sx | ✅ | Throws "Identifier not found" error |
| debug.sx | ✅ | Returns 65 (only shows final result) |
| token.sx | ❌ | Parser errors |
| voting.sx | ❌ | Parser errors |
| auction.sx | ❌ | Parser errors |

## Observations

1. The interpreter correctly handles basic arithmetic operations, string concatenation, and conditional expressions.
2. Error handling works as expected, with appropriate error messages for runtime errors like division by zero.
3. Block scoping is now properly implemented - variables defined in blocks are not accessible outside their blocks.
4. Boolean literals (true/false) are now directly supported.
5. Logical operators (&&, ||) with short-circuit evaluation are now supported.
6. Only the last expression's result is returned, which is consistent with the interpreter's design.
7. The advanced examples (token.sx, voting.sx, auction.sx) fail with parser errors because the parser doesn't support contract syntax and blockchain-specific features yet.
8. The CLI is very simple and only supports running a file with the `run -file` command. There are no debug or verbose flags available.

## Current Limitations

The current implementation has the following limitations:

1. No support for contract syntax and blockchain-specific features.
2. Limited type system (no custom types or type annotations).
3. No support for functions or function calls.
4. No support for arrays or maps.
5. String and integer concatenation is not supported.
6. Only the last expression's result is returned.
7. No debugging or verbose output options in the CLI.

## Recent Improvements

1. **Block Scoping**: Variables defined in blocks are now properly scoped and not accessible outside their blocks.
2. **Boolean Literals**: Direct support for boolean literals (true/false) has been added.
3. **Logical Operators**: Support for logical AND (&&) and OR (||) operators with short-circuit evaluation has been implemented.

## Next Steps

1. Implement support for string and integer concatenation.
2. Add support for arrays and maps.
3. Implement function declarations and function calls.
4. Add support for contract syntax and blockchain-specific features.
5. Enhance the type system to support custom types and type annotations.
6. Add debugging and verbose output options to the CLI.
7. Implement a REPL (Read-Eval-Print Loop) for interactive development. 