package parser

import (
	"github.com/Stremax-Team/stremax-lang/pkg/lexer"
	"testing"
)

// Sample code for benchmarking
const benchmarkInput = `
let x = 5;
let y = 10;
let result = x + y;

if (result > 10) {
	return true;
} else {
	return false;
}

// This is a comment
/* This is a block comment */

contract SimpleContract {
	state {
		owner: Address
		balance: Int
	}
	
	constructor() {
		owner = msg.sender;
		balance = 0;
	}
	
	function deposit(amount: Int) {
		balance += amount;
	}
	
	function withdraw(amount: Int) {
		require(amount <= balance, "Insufficient balance");
		balance -= amount;
	}
}
`

// BenchmarkParser benchmarks the parser's performance
func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := lexer.New(benchmarkInput)
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkParseExpression benchmarks the expression parsing
func BenchmarkParseExpression(b *testing.B) {
	input := "5 + 10 * 2 + 20 / 4 - 8"
	l := lexer.New(input)
	p := New(l)

	// Skip to the expression
	p.nextToken()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset the parser for each iteration
		l = lexer.New(input)
		p = New(l)
		p.nextToken() // Skip to the expression
		p.parseExpression(LOWEST)
	}
}
