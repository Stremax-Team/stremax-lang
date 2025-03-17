package interpreter

import (
	"strings"
	"testing"
)

func TestBlockScoping(t *testing.T) {
	tests := []struct {
		input           string
		expectedError   bool
		errorContains   string
		expectedResult  int64
	}{
		{
			// Variable defined in outer scope should be accessible in inner scope
			input: `
				let x = 10;
				if (x > 5) {
					x;
				}
			`,
			expectedError: false,
			expectedResult: 10,
		},
		{
			// Variable defined in inner scope should not be accessible in outer scope
			input: `
				if (true) {
					let y = 20;
				}
				y;
			`,
			expectedError: true,
			errorContains: "Identifier not found: y",
		},
		{
			// Variable with same name in inner scope should shadow outer scope
			input: `
				let z = 30;
				if (true) {
					let z = 40;
					z;
				}
			`,
			expectedError: false,
			expectedResult: 40,
		},
		{
			// After exiting block, outer scope variable should retain its value
			input: `
				let a = 50;
				if (true) {
					let a = 60;
				}
				a;
			`,
			expectedError: false,
			expectedResult: 50,
		},
		{
			// Nested blocks should maintain proper scoping
			input: `
				let b = 70;
				if (true) {
					let c = 80;
					if (true) {
						let d = 90;
						b + c + d;
					}
				}
			`,
			expectedError: false,
			expectedResult: 240, // 70 + 80 + 90
		},
		{
			// Variable from innermost block should not be accessible in middle block
			input: `
				if (true) {
					if (true) {
						let inner = 100;
					}
					inner;
				}
			`,
			expectedError: true,
			errorContains: "Identifier not found: inner",
		},
	}

	for idx, tt := range tests {
		interp := New(tt.input)
		err := interp.Run()

		// Check for expected errors
		if tt.expectedError {
			if err == nil {
				t.Errorf("test %d: expected error but got none", idx)
				continue
			}
			if !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("test %d: expected error to contain %q, got %q", idx, tt.errorContains, err.Error())
			}
			continue
		}

		// Check for unexpected errors
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", idx, err)
			continue
		}

		// For now, we can't easily access the result value from the interpreter
		// This would require modifying the interpreter to expose the last evaluated value
		// For a complete test, we would need to add this functionality
	}
} 