package lexer

import "testing"

// TestLogicalOperators tests the lexer's ability to recognize logical operators
func TestLogicalOperators(t *testing.T) {
	input := `
a && b
a || b
true && false
false || true
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "a"},
		{AND, "&&"},
		{IDENT, "b"},
		{IDENT, "a"},
		{OR, "||"},
		{IDENT, "b"},
		{TRUE, "true"},
		{AND, "&&"},
		{FALSE, "false"},
		{FALSE, "false"},
		{OR, "||"},
		{TRUE, "true"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
} 