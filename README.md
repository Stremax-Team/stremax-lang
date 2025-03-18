# Stremax-Lang

[![Go Reference](https://pkg.go.dev/badge/github.com/Stremax-Team/stremax-lang.svg)](https://pkg.go.dev/github.com/Stremax-Team/stremax-lang)
[![Go Report Card](https://goreportcard.com/badge/github.com/Stremax-Team/stremax-lang)](https://goreportcard.com/report/github.com/Stremax-Team/stremax-lang)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/release/Stremax-Team/stremax-lang.svg)](https://github.com/Stremax-Team/stremax-lang/releases)
[![Build Status](https://github.com/Stremax-Team/stremax-lang/workflows/Go/badge.svg)](https://github.com/Stremax-Team/stremax-lang/actions)

Stremax-Lang is a simple programming language designed specifically for blockchain development and smart contract execution. It provides a clean syntax and built-in primitives for common blockchain operations.

## Features

- Simple, expressive syntax
- Built-in blockchain primitives (addresses, transactions, etc.)
- Smart contract support
- Type safety for blockchain operations
- Easy integration with existing blockchain platforms
- Proper block scoping for variables
- Logical operators with short-circuit evaluation

## Current Implementation Status

- ✅ Lexer: Tokenizes source code into tokens
- ⏳ Parser: Partially implemented, needs completion
- ⏳ Interpreter: Basic structure implemented, needs completion
- ✅ Blockchain: Basic blockchain implementation with blocks, transactions, and mining
- ✅ Smart Contracts: Basic smart contract implementation with state, functions, and events
- ✅ CLI: Command-line interface for running Stremax-Lang programs

## Recent Improvements

- ✅ Block Scoping: Variables defined in blocks are now properly scoped
- ✅ Boolean Literals: Direct support for boolean literals (true/false)
- ✅ Logical Operators: Support for logical AND (&&) and OR (||) with short-circuit evaluation
- ✅ Functions: Support for function declarations, calls, closures, and recursion
- ✅ String/Integer Concatenation: Enhanced support for string concatenation with different types
- ✅ Arrays: Support for array literals, array access, and nested arrays
- ✅ Maps: Support for map/hash/dictionary literals with string, integer, or boolean keys

## Project Structure

- `cmd/stremax`: Command-line tools for the language
- `pkg/lexer`: Lexical analyzer for tokenizing source code
- `pkg/parser`: Parser for building abstract syntax trees
- `pkg/interpreter`: Interpreter for executing Stremax-Lang code
- `pkg/blockchain`: Blockchain-specific functionality
- `examples`: Example programs written in Stremax-Lang

## Language Syntax (Draft)

```
// Contract definition
contract TokenContract {
    // State variables
    state {
        owner: Address
        totalSupply: Int
        balances: Map<Address, Int>
    }

    // Constructor
    constructor(initialSupply: Int) {
        owner = msg.sender
        totalSupply = initialSupply
        balances[owner] = initialSupply
    }

    // Functions
    function transfer(to: Address, amount: Int) {
        require(balances[msg.sender] >= amount, "Insufficient balance")
        
        balances[msg.sender] -= amount
        balances[to] += amount
        
        emit Transfer(msg.sender, to, amount)
    }

    // Events
    event Transfer(from: Address, to: Address, amount: Int)
}
```

## Language Features

### Basic Types

- `Int`: Integer values
- `String`: String values
- `Bool`: Boolean values (true/false)
- `Address`: Blockchain addresses

### Operators

- **Arithmetic**: `+`, `-`, `*`, `/`
- **Comparison**: `==`, `!=`, `<`, `>`
- **Logical**: `&&` (AND), `||` (OR), `!` (NOT)

### Blockchain-Specific Types

- `Address`: Represents a blockchain address
- `Map<K, V>`: Key-value mapping (e.g., `Map<Address, Int>` for balances)

### Contract Structure

```
contract ContractName {
    // State variables
    state {
        variable1: Type
        variable2: Type
        // ...
    }

    // Constructor
    constructor(param1: Type, param2: Type) {
        // Initialization code
    }

    // Functions
    function functionName(param1: Type, param2: Type): ReturnType {
        // Function body
    }

    // Events
    event EventName(param1: Type, param2: Type)
}
```

### Special Variables

- `msg.sender`: The address that called the current function
- `msg.value`: The amount of cryptocurrency sent with the function call
- `now()`: The current timestamp

### Control Flow

- `if (condition) { ... } else { ... }`: Conditional execution
- `require(condition, "error message")`: Assert a condition or revert the transaction

### Blockchain Operations

- `emit EventName(arg1, arg2)`: Emit an event
- `address.transfer(amount)`: Transfer cryptocurrency to an address
- `address.send(amount)`: Send cryptocurrency to an address (returns success/failure)

## Examples

### Basic Examples

See the [examples](examples/) directory for basic examples demonstrating language features:
- Variable assignment and arithmetic
- String operations
- Conditional expressions
- Boolean operations and logical operators
- Error handling
- Block scoping

### Advanced Examples

See these examples for more complex use cases (not fully supported yet):
- [examples/token.sx](examples/token.sx): An ERC20-like token implementation
- [examples/voting.sx](examples/voting.sx): A voting contract implementation
- [examples/auction.sx](examples/auction.sx): An auction contract implementation

## Getting Started

### Installation

```bash
# Clone the repository
git clone https://github.com/Stremax-Team/stremax-lang.git
cd stremax-lang

# Build the compiler
go build -o stremax ./cmd/stremax
```

### Running a Stremax-Lang Program

```bash
# Run a program
./stremax run -file ./examples/simple.sx
```

### Deploying a Contract

```bash
# Deploy a contract to a blockchain (not implemented yet) (result with error)
./stremax deploy -file ./examples/token.sx
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Documentation

The API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/Stremax-Team/stremax-lang). This documentation is automatically generated from the code comments.

To write good documentation comments:

1. Every exported function, type, and variable should have a comment
2. Comments for functions should explain what the function does, its parameters, and its return values
3. Use complete sentences that start with the name of the thing being described
4. Follow the [Go Documentation Comments](https://go.dev/doc/comment) guidelines

Example of a well-documented function:

```go
// ParseProgram parses the input source code and returns an AST representation.
// It returns nil and populates the Errors slice if any parsing errors occur.
func (p *Parser) ParseProgram() *Program {
    // Implementation...
}
```

## License

MIT 