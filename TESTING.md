# Stremax-Lang Testing

This document summarizes the testing of the Stremax-Lang interpreter.

## Test Coverage

The Stremax-Lang interpreter includes comprehensive tests for the following features:

1. **Lexical Analysis**: Tests for tokenizing all language constructs
2. **Parsing**: Tests for creating correct AST nodes for all syntax elements
3. **Evaluation**: Tests for correctly evaluating expressions and statements
4. **Error Handling**: Tests for appropriate error messages
5. **Functions**: Tests for function declarations, calls, closures, and recursion
6. **Logical Operations**: Tests for logical operators with short-circuit evaluation
7. **Block Scoping**: Tests for variable scoping rules
8. **String Concatenation**: Tests for string and mixed-type concatenation
9. **Collections**: Tests for arrays and maps including creation and element access

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

Or you can run tests for specific packages:

```bash
go test ./pkg/lexer
go test ./pkg/parser
go test ./pkg/interpreter
```

To run tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Example Programs

The `examples/` directory contains programs that demonstrate various language features:

- Basic features: `simple.sx`, `arithmetic.sx`, `strings.sx`, etc.
- Logical operations: `logical.sx`, `simple_logical.sx`
- Functions: `functions.sx`, `functions_simple.sx`, `functions_advanced.sx`, `factorial.sx`
- String concatenation: `concatenation.sx`, `string_concat.sx`
- Collections: `arrays.sx`, `maps.sx`, `collections.sx`

Each example program is tested to ensure it produces the expected result.

## Benchmarks

Performance benchmarks are included for critical components:

- Lexer: `pkg/lexer/lexer_benchmark_test.go`
- Parser: `pkg/parser/parser_benchmark_test.go`
- Interpreter: `pkg/interpreter/interpreter_benchmark_test.go`

Run benchmarks with:

```bash
go test -bench=. ./...
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