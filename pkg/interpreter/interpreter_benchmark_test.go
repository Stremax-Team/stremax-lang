package interpreter

import (
	"github.com/Stremax-Team/stremax-lang/pkg/lexer"
	"github.com/Stremax-Team/stremax-lang/pkg/parser"
	"testing"
)

// Sample code for benchmarking
const benchmarkInput = `
let x = 5;
let y = 10;
let result = x + y;

if (result > 10) {
	result = result * 2;
} else {
	result = result / 2;
}

result;
`

// BenchmarkInterpreter benchmarks the interpreter's performance
func BenchmarkInterpreter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		interpreter := New(benchmarkInput)
		interpreter.Run()
	}
}

// BenchmarkEvalExpression benchmarks expression evaluation
func BenchmarkEvalExpression(b *testing.B) {
	input := "5 + 10 * 2 + 20 / 4 - 8"
	l := lexer.New(input)
	p := parser.New(l)
	expr := p.ParseProgram().Statements[0].(*parser.ExpressionStatement).Expression

	interpreter := New("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		interpreter.evalExpression(expr)
	}
}

// BenchmarkArithmeticOperations benchmarks arithmetic operations
func BenchmarkArithmeticOperations(b *testing.B) {
	input := `
	let a = 10;
	let b = 20;
	let c = 30;
	let d = 40;
	let result = a + b * c - d / 2;
	result;
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		interpreter := New(input)
		interpreter.Run()
	}
}
