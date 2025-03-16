package lexer

import (
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

// BenchmarkLexer benchmarks the lexer's performance
func BenchmarkLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := New(benchmarkInput)
		for {
			tok := l.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkNextToken benchmarks just the NextToken method
func BenchmarkNextToken(b *testing.B) {
	l := New(benchmarkInput)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.NextToken()
		// Reset the lexer when we reach EOF
		if l.ch == 0 {
			l = New(benchmarkInput)
		}
	}
}
