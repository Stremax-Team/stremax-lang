package interpreter

import (
	"strings"
	"testing"
)

func TestEnhancedStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// String + String
		{`"Hello" + " " + "World";`, "Hello World"},
		{`"" + "empty";`, "empty"},
		{`"empty" + "";`, "empty"},
		{`"" + "";`, ""},
		
		// String + Integer
		{`"The answer is: " + 42;`, "The answer is: 42"},
		{`"Value: " + 0;`, "Value: 0"},
		{`"Negative: " + -10;`, "Negative: -10"},
		
		// Integer + String
		{`100 + " points";`, "100 points"},
		{`0 + " zero";`, "0 zero"},
		
		// String + Boolean
		{`"Is it true? " + true;`, "Is it true? true"},
		{`"Is it false? " + false;`, "Is it false? false"},
		
		// Boolean + String
		{`true + " is correct";`, "true is correct"},
		{`false + " is incorrect";`, "false is incorrect"},
		
		// Complex concatenation
		{`"Count: " + 1 + 2 + 3;`, "Count: 123"},
		{`"Valid: " + (10 > 5);`, "Valid: true"},
		{`"Result: " + (10 + 20);`, "Result: 30"},
		
		// Empty string with various types
		{`"" + 123;`, "123"},
		{`"" + true;`, "true"},
		{`123 + "";`, "123"},
		{`true + "";`, "true"},
		
		// Long string concatenation
		{
			`"This is a very long string. " + "It contains multiple parts. " + "We want to ensure " + "it works correctly " + "with many concatenations.";`,
			"This is a very long string. It contains multiple parts. We want to ensure it works correctly with many concatenations.",
		},
		
		// String concatenation in expressions
		{`let a = "Hello"; let b = " World"; a + b;`, "Hello World"},
		{`let x = 42; "The value of x is: " + x;`, "The value of x is: 42"},
		{`let isTrue = true; "The statement is " + isTrue;`, "The statement is true"},
		
		// Concatenation with the result of a conditional expression
		{`"The result is: " + if (10 > 5) { "positive" } else { "negative" };`, "The result is: positive"},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		
		strObj, ok := evaluated.(*String)
		if !ok {
			t.Errorf("test %d: object is not String. got=%T", i, evaluated)
			continue
		}
		
		if strObj.Value != tt.expected {
			t.Errorf("test %d: wrong string value. expected=%q, got=%q", i, tt.expected, strObj.Value)
		}
	}
}

func TestStringConcatenationErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedContains string
	}{
		// Using operators other than + for strings
		{`"Hello" - "World";`, "Unknown operator: -"},
		{`"Hello" * "World";`, "Unknown operator: *"},
		{`"Hello" / "World";`, "Unknown operator: /"},
	}

	for i, tt := range tests {
		interpreter := New(tt.input)
		err := interpreter.Run()
		
		if err == nil {
			t.Errorf("test %d: expected error but got none", i)
			continue
		}
		
		if !strings.Contains(err.Error(), tt.expectedContains) {
			t.Errorf("test %d: expected error to contain %q, got %q", 
				i, tt.expectedContains, err.Error())
		}
	}
} 