# Stremax-Lang Examples

This directory contains example programs written in Stremax-Lang.

## Simple Examples

These examples demonstrate the basic features of Stremax-Lang that are currently implemented:

- **simple.sx**: A very basic example with variable assignment and arithmetic.
- **arithmetic.sx**: Demonstrates various arithmetic operations.
- **strings.sx**: Shows string operations like concatenation.
- **conditionals.sx**: Illustrates conditional expressions with if-else statements.
- **boolean.sx**: Tests boolean operations and comparisons.
- **logical.sx**: Demonstrates logical operators (&&, ||) with short-circuit evaluation.
- **simple_logical.sx**: A minimal example of logical operators.
- **functions_simple.sx**: Basic function declarations and calls.
- **factorial.sx**: Recursive functions with the factorial example.
- **functions_advanced.sx**: Advanced functions including closures.
- **combined.sx**: Combines various features in one example.
- **errors.sx**: Demonstrates error handling (division by zero).
- **scoping.sx**: Tests variable scoping.
- **scoping_error.sx**: Tests variable scoping errors.
- **debug.sx**: Tests interpreter output with multiple expressions.
- **concatenation.sx**: Demonstrates string and integer concatenation features.
- **arrays.sx**: Shows array creation, access, and operations.
- **maps.sx**: Shows map/dictionary creation, access, and operations.
- **collections.sx**: Comprehensive example of arrays and maps working together.

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
3. ~~No support for arrays or maps.~~
4. ~~String and integer concatenation is not fully supported.~~
5. Only the last expression's result is returned.
6. No debugging or verbose output options in the CLI.

## Recent Improvements

1. **Block Scoping**: Variables defined in blocks are now properly scoped and not accessible outside their blocks.
2. **Boolean Literals**: Direct support for boolean literals (true/false) has been added.
3. **Logical Operators**: Support for logical AND (&&) and OR (||) operators with short-circuit evaluation has been implemented.
4. **Functions**: Full support for function declarations, function calls, closures, and recursion has been added.
5. **String Concatenation**: Enhanced support for string concatenation with different types.
6. **Arrays**: Support for array literals, array access, and nested arrays.
7. **Maps**: Support for map/hash/dictionary literals with string, integer, or boolean keys.

These limitations will be addressed in future versions of Stremax-Lang. 