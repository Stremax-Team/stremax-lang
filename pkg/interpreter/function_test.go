package interpreter

import (
	"strings"
	"testing"
)

func TestSimpleFunctionCall(t *testing.T) {
	input := `
	let adder = function(a, b) { a + b; };
	adder(2, 3);
	`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 5)
}

func TestFunctionLetStatement(t *testing.T) {
	input := `
	let add = function(a, b) {
		return a + b;
	};
	let result = add(5, 5);
	result;
	`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 10)
}

func TestNestedFunctionCalls(t *testing.T) {
	input := `
	let add = function(a, b) { a + b; };
	let applyFunc = function(a, b, f) { f(a, b); };
	applyFunc(2, 3, add);
	`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 5)
}

func TestClosureFunction(t *testing.T) {
	input := `
	let newAdder = function(x) {
		return function(y) { x + y; };
	};
	
	let addTwo = newAdder(2);
	addTwo(3);
	`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 5)
}

func TestRecursiveFunction(t *testing.T) {
	input := `
	let factorial = function(n) {
		if (n == 0) {
			return 1;
		} else {
			return n * factorial(n - 1);
		}
	};
	
	factorial(5);
	`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 120)
}

func TestFunctionWithReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`
			let noReturn = function(x) { x; };
			noReturn(10);
			`,
			10,
		},
		{
			`
			let earlyReturn = function(x) {
				return x;
				x + 1; // This should not be evaluated
			};
			earlyReturn(10);
			`,
			10,
		},
		{
			`
			let conditionalReturn = function(x) {
				if (x > 5) {
					return 10;
				}
				return 5;
			};
			conditionalReturn(7);
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			"let add = function(a, b) { a + b; }; add(1);",
			"Wrong number of arguments",
		},
		{
			"let add = function(a, b) { a + b; }; add(1, 2, 3);",
			"Wrong number of arguments",
		},
		{
			"1(1);",
			"Not a function",
		},
	}

	for _, tt := range tests {
		interpreter := New(tt.input)
		err := interpreter.Run()

		if err == nil {
			t.Errorf("expected error but got none for input: %s", tt.input)
			continue
		}

		if !strings.Contains(err.Error(), tt.expectedError) {
			t.Errorf("wrong error message. expected=%q to contain %q for input: %s",
				err.Error(), tt.expectedError, tt.input)
		}
	}
} 