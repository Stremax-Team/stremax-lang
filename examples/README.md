# Stremax-Lang Examples

This directory contains example programs written in Stremax-Lang.

## Simple Examples

These examples demonstrate the basic features of Stremax-Lang that are currently implemented:

- **simple.sx**: A very basic example with variable assignment and arithmetic.
- **arithmetic.sx**: Demonstrates various arithmetic operations.
- **strings.sx**: Shows string operations like concatenation.
- **conditionals.sx**: Illustrates conditional expressions with if-else statements.
- **boolean.sx**: Tests boolean operations and comparisons.
- **combined.sx**: Combines various features in one example.
- **errors.sx**: Demonstrates error handling (division by zero).
- **scoping.sx**: Tests variable scoping.
- **scoping_error.sx**: Tests variable scoping errors.
- **debug.sx**: Tests interpreter output with multiple expressions.

## Advanced Examples (Not Fully Supported Yet)

These examples showcase the planned features of Stremax-Lang, but they may not work with the current implementation:

- **token.sx**: An ERC20-like token contract implementation.
- **voting.sx**: A voting contract implementation.
- **auction.sx**: An auction contract implementation.

## Running the Examples

To run an example, use the `stremax run` command:

```bash
./stremax run -file ./examples/simple.sx
```

## Current Limitations

The current implementation has the following limitations:

1. No support for contract syntax and blockchain-specific features.
2. Limited type system (no custom types or type annotations).
3. No support for functions or function calls.
4. Block scoping is not fully implemented (variables defined in blocks are accessible outside).
5. No support for arrays or maps.
6. String and integer concatenation is not supported.
7. No direct support for boolean literals (true/false).
8. Only the last expression's result is returned.
9. No debugging or verbose output options in the CLI.

These limitations will be addressed in future versions of Stremax-Lang. 