package interpreter

import (
	"github.com/Stremax-Team/stremax-lang/pkg/lexer"
	"github.com/Stremax-Team/stremax-lang/pkg/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
		{"-5;", -5},
		{"-10;", -10},
		{"5 + 5 + 5 + 5 - 10;", 10},
		{"2 * 2 * 2 * 2 * 2;", 32},
		{"-50 + 100 + -50;", 0},
		{"5 * 2 + 10;", 20},
		{"5 + 2 * 10;", 25},
		{"20 + 2 * -10;", 0},
		{"50 / 2 * 2 + 10;", 60},
		{"2 * (5 + 10);", 30},
		{"3 * 3 * 3 + 10;", 37},
		{"3 * (3 * 3) + 10;", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1 < 2;", true},
		{"1 > 2;", false},
		{"1 < 1;", false},
		{"1 > 1;", false},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"1 == 2;", false},
		{"1 != 2;", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (1 < 2) { 10; }", 10},
		{"if (1 > 2) { 10; }", nil},
		{"if (1 > 2) { 10; } else { 20; }", 20},
		{"if (1 < 2) { 10; } else { 20; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	input := `"Hello World!";`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!";`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true;",
			"Cannot negate non-integer",
		},
		{
			"true + false;",
			"Type mismatch: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Type mismatch: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Type mismatch: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"Identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		i := New(tt.input)

		_, err := i.evalProgram(program)

		if err == nil {
			t.Errorf("no error returned. expected to contain=%q", tt.expectedMessage)
			continue
		}

		if !contains(err.Error(), tt.expectedMessage) {
			t.Errorf("wrong error message. expected to contain=%q, got=%q",
				tt.expectedMessage, err.Error())
		}
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func testEval(t *testing.T, input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	i := New(input)

	result, err := i.evalProgram(program)
	if err != nil {
		t.Fatalf("evalProgram error: %s", err)
	}

	return result
}

func testIntegerObject(t *testing.T, obj Object, expected int64) {
	result, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}
}

func testBooleanObject(t *testing.T, obj Object, expected bool) {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}
}

func testNullObject(t *testing.T, obj Object) {
	if obj != nil {
		t.Errorf("object is not nil. got=%T (%+v)", obj, obj)
	}
}
