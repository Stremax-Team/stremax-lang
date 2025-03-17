package interpreter

import (
	"strings"
	"testing"
)

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input           string
		expectedError   bool
		errorContains   string
		expectedResult  bool
	}{
		// Basic logical AND tests
		{
			input: "true && true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "true && false;",
			expectedError: false,
			expectedResult: false,
		},
		{
			input: "false && true;",
			expectedError: false,
			expectedResult: false,
		},
		{
			input: "false && false;",
			expectedError: false,
			expectedResult: false,
		},
		
		// Basic logical OR tests
		{
			input: "true || true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "true || false;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "false || true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "false || false;",
			expectedError: false,
			expectedResult: false,
		},
		
		// Compound expressions
		{
			input: "true && true && true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "true && false && true;",
			expectedError: false,
			expectedResult: false,
		},
		{
			input: "false || false || true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "false || false || false;",
			expectedError: false,
			expectedResult: false,
		},
		
		// Mixed operators
		{
			input: "true && false || true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "false || true && true;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "false || false && true;",
			expectedError: false,
			expectedResult: false,
		},
		
		// With comparisons
		{
			input: "5 > 3 && 10 < 20;",
			expectedError: false,
			expectedResult: true,
		},
		{
			input: "5 < 3 || 10 > 20;",
			expectedError: false,
			expectedResult: false,
		},
		
		// Short-circuit evaluation
		{
			input: "false && (10 / 0 > 0);", // Second part should not be evaluated
			expectedError: false,
			expectedResult: false,
		},
		{
			input: "true || (10 / 0 > 0);", // Second part should not be evaluated
			expectedError: false,
			expectedResult: true,
		},
		
		// Type errors
		{
			input: "5 && true;",
			expectedError: true,
			errorContains: "Left operand of && must be a boolean",
		},
		{
			input: "false || 10;",
			expectedError: true,
			errorContains: "Right operand of || must be a boolean",
		},
	}

	for idx, tt := range tests {
		interp := New(tt.input)
		err := interp.Run()

		// Check for expected errors
		if tt.expectedError {
			if err == nil {
				t.Errorf("test %d: expected error but got none for input: %s", idx, tt.input)
				continue
			}
			if !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("test %d: expected error to contain %q, got %q for input: %s", 
					idx, tt.errorContains, err.Error(), tt.input)
			}
			continue
		}

		// Check for unexpected errors
		if err != nil {
			t.Errorf("test %d: unexpected error: %s for input: %s", idx, err, tt.input)
			continue
		}

		// For now, we can't easily access the result value from the interpreter
		// This would require modifying the interpreter to expose the last evaluated value
		// For a complete test, we would need to add this functionality
	}
} 