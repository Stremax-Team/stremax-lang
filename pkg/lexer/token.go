package lexer

// TokenType represents the type of a token
type TokenType string

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Token types
const (
	// Special tokens
	ILLEGAL = "ILLEGAL" // Token we don't know about
	EOF     = "EOF"     // End of file
	
	// Identifiers and literals
	IDENT  = "IDENT"  // add, x, y, etc.
	INT    = "INT"    // 123456
	STRING = "STRING" // "hello"
	
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="
	
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."
	
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	
	// Keywords
	FUNCTION   = "FUNCTION"
	CONTRACT   = "CONTRACT"
	STATE      = "STATE"
	LET        = "LET"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	IF         = "IF"
	ELSE       = "ELSE"
	RETURN     = "RETURN"
	REQUIRE    = "REQUIRE"
	EMIT       = "EMIT"
	EVENT      = "EVENT"
	ADDRESS    = "ADDRESS"
	MAP        = "MAP"
	CONSTRUCTOR = "CONSTRUCTOR"
)

// Keywords maps string literals to their token types
var Keywords = map[string]TokenType{
	"function":    FUNCTION,
	"contract":    CONTRACT,
	"state":       STATE,
	"let":         LET,
	"true":        TRUE,
	"false":       FALSE,
	"if":          IF,
	"else":        ELSE,
	"return":      RETURN,
	"require":     REQUIRE,
	"emit":        EMIT,
	"event":       EVENT,
	"Address":     ADDRESS,
	"Map":         MAP,
	"constructor": CONSTRUCTOR,
}

// LookupIdent checks if the given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
} 