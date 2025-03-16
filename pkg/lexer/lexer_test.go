package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		let five = 5;
		let ten = 10;
		
		function add(x: Int, y: Int): Int {
			return x + y;
		}
		
		let result = add(five, ten);
		
		if (5 < 10) {
			return true;
		} else {
			return false;
		}
		
		10 == 10;
		10 != 9;
		
		// This is a line comment
		
		/*
			This is a block comment
		*/
		
		contract TokenContract {
			state {
				owner: Address
				totalSupply: Int
				balances: Map<Address, Int>
			}
			
			constructor(initialSupply: Int) {
				owner = msg.sender;
				totalSupply = initialSupply;
				balances[owner] = initialSupply;
			}
			
			function transfer(to: Address, amount: Int) {
				require(balances[msg.sender] >= amount, "Insufficient balance");
				
				balances[msg.sender] -= amount;
				balances[to] += amount;
				
				emit Transfer(msg.sender, to, amount);
			}
			
			event Transfer(from: Address, to: Address, amount: Int)
		}
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},

		{FUNCTION, "function"},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COLON, ":"},
		{IDENT, "Int"},
		{COMMA, ","},
		{IDENT, "y"},
		{COLON, ":"},
		{IDENT, "Int"},
		{RPAREN, ")"},
		{COLON, ":"},
		{IDENT, "Int"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},

		{LET, "let"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},

		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},

		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},

		{INT, "10"},
		{NotEq, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},

		{CONTRACT, "contract"},
		{IDENT, "TokenContract"},
		{LBRACE, "{"},

		{STATE, "state"},
		{LBRACE, "{"},
		{IDENT, "owner"},
		{COLON, ":"},
		{ADDRESS, "Address"},
		{IDENT, "totalSupply"},
		{COLON, ":"},
		{IDENT, "Int"},
		{IDENT, "balances"},
		{COLON, ":"},
		{MAP, "Map"},
		{LT, "<"},
		{ADDRESS, "Address"},
		{COMMA, ","},
		{IDENT, "Int"},
		{GT, ">"},
		{RBRACE, "}"},

		{CONSTRUCTOR, "constructor"},
		{LPAREN, "("},
		{IDENT, "initialSupply"},
		{COLON, ":"},
		{IDENT, "Int"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "owner"},
		{ASSIGN, "="},
		{IDENT, "msg"},
		{DOT, "."},
		{IDENT, "sender"},
		{SEMICOLON, ";"},

		{IDENT, "totalSupply"},
		{ASSIGN, "="},
		{IDENT, "initialSupply"},
		{SEMICOLON, ";"},

		{IDENT, "balances"},
		{LBRACKET, "["},
		{IDENT, "owner"},
		{RBRACKET, "]"},
		{ASSIGN, "="},
		{IDENT, "initialSupply"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},

		{FUNCTION, "function"},
		{IDENT, "transfer"},
		{LPAREN, "("},
		{IDENT, "to"},
		{COLON, ":"},
		{ADDRESS, "Address"},
		{COMMA, ","},
		{IDENT, "amount"},
		{COLON, ":"},
		{IDENT, "Int"},
		{RPAREN, ")"},
		{LBRACE, "{"},

		{REQUIRE, "require"},
		{LPAREN, "("},
		{IDENT, "balances"},
		{LBRACKET, "["},
		{IDENT, "msg"},
		{DOT, "."},
		{IDENT, "sender"},
		{RBRACKET, "]"},
		{GT, ">"},
		{ASSIGN, "="},
		{IDENT, "amount"},
		{COMMA, ","},
		{STRING, "Insufficient balance"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},

		{IDENT, "balances"},
		{LBRACKET, "["},
		{IDENT, "msg"},
		{DOT, "."},
		{IDENT, "sender"},
		{RBRACKET, "]"},
		{MINUS, "-"},
		{ASSIGN, "="},
		{IDENT, "amount"},
		{SEMICOLON, ";"},

		{IDENT, "balances"},
		{LBRACKET, "["},
		{IDENT, "to"},
		{RBRACKET, "]"},
		{PLUS, "+"},
		{ASSIGN, "="},
		{IDENT, "amount"},
		{SEMICOLON, ";"},

		{EMIT, "emit"},
		{IDENT, "Transfer"},
		{LPAREN, "("},
		{IDENT, "msg"},
		{DOT, "."},
		{IDENT, "sender"},
		{COMMA, ","},
		{IDENT, "to"},
		{COMMA, ","},
		{IDENT, "amount"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},

		{EVENT, "event"},
		{IDENT, "Transfer"},
		{LPAREN, "("},
		{IDENT, "from"},
		{COLON, ":"},
		{ADDRESS, "Address"},
		{COMMA, ","},
		{IDENT, "to"},
		{COLON, ":"},
		{ADDRESS, "Address"},
		{COMMA, ","},
		{IDENT, "amount"},
		{COLON, ":"},
		{IDENT, "Int"},
		{RPAREN, ")"},
		{RBRACE, "}"},

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
