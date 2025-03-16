# Stremax-Lang

Stremax-Lang is a simple programming language designed specifically for blockchain development and smart contract execution. It provides a clean syntax and built-in primitives for common blockchain operations.

## Features

- Simple, expressive syntax
- Built-in blockchain primitives (addresses, transactions, etc.)
- Smart contract support
- Type safety for blockchain operations
- Easy integration with existing blockchain platforms

## Current Implementation Status

- ✅ Lexer: Tokenizes source code into tokens
- ⏳ Parser: Partially implemented, needs completion
- ⏳ Interpreter: Basic structure implemented, needs completion
- ✅ Blockchain: Basic blockchain implementation with blocks, transactions, and mining
- ✅ Smart Contracts: Basic smart contract implementation with state, functions, and events
- ✅ CLI: Command-line interface for running Stremax-Lang programs

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

### Token Contract

See [examples/token.sx](examples/token.sx) for a complete ERC20-like token implementation.

### Voting Contract

See [examples/voting.sx](examples/voting.sx) for a voting contract implementation.

### Auction Contract

See [examples/auction.sx](examples/auction.sx) for an auction contract implementation.

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

## License

MIT 