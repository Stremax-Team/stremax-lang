package parser

import (
	"fmt"
	"github.com/Stremax-Team/stremax-lang/pkg/lexer"
	"strconv"
)

// Precedence levels for operators
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	DOT         // obj.property
)

// Operator precedence map
var precedences = map[lexer.TokenType]int{
	lexer.EQ:       EQUALS,
	lexer.NotEq:    EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.SLASH:    PRODUCT,
	lexer.ASTERISK: PRODUCT,
	lexer.LPAREN:   CALL,
	lexer.LBRACKET: INDEX,
	lexer.DOT:      DOT,
}

// Parser represents a parser for Stremax-Lang.
// It implements a recursive descent parser with Pratt parsing for expressions,
// which allows for handling operator precedence and various expression types.
// The parser builds an abstract syntax tree (AST) from the token stream
// provided by the lexer.
type Parser struct {
	l         *lexer.Lexer
	errors    []string
	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// New creates a new Parser with the given lexer.
// It initializes the parser state, reads the first two tokens,
// and registers all the parsing functions for different expression types.
//
// Parameters:
//   - l: The lexer that provides the token stream
//
// Returns:
//   - A new Parser instance ready to parse Stremax-Lang code
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	// Register prefix parse functions
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.IF, p.parseIfExpression)

	// Register infix parse functions
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NotEq, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parseDotExpression)

	return p
}

// Errors returns all errors encountered during parsing.
// This can be used to check if parsing was successful and
// to report any syntax errors to the user.
//
// Returns:
//   - A slice of error messages as strings
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses a complete Stremax-Lang program.
// It processes the token stream until it reaches the end of file,
// building an abstract syntax tree (AST) representation of the program.
// The function handles all statement types in the language and
// collects any parsing errors encountered.
//
// Returns:
//   - A Program struct containing the AST of the parsed program
//   - If parsing errors occur, they can be retrieved using the Errors() method
func (p *Parser) ParseProgram() *Program {
	program := &Program{
		Statements: []Statement{},
	}

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.CONTRACT:
		return p.parseContractStatement()
	case lexer.FUNCTION:
		return p.parseFunctionStatement()
	case lexer.CONSTRUCTOR:
		return p.parseConstructorStatement()
	case lexer.EVENT:
		return p.parseEventStatement()
	case lexer.REQUIRE:
		return p.parseRequireStatement()
	case lexer.EMIT:
		return p.parseEmitStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses a let statement
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for type annotation
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // Skip the colon
		p.nextToken() // Move to the type
		stmt.Type = p.parseTypeExpression()
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseContractStatement parses a contract statement
func (p *Parser) parseContractStatement() *ContractStatement {
	stmt := &ContractStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Parse the contract body
	body := &BlockStatement{Token: p.curToken}
	body.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		if p.curTokenIs(lexer.STATE) {
			// Parse state block
			stmt.StateBlock = p.parseStateBlockStatement()
		} else {
			// Parse other statements (functions, constructors, events)
			statement := p.parseStatement()
			if statement != nil {
				body.Statements = append(body.Statements, statement)
			}
		}
		p.nextToken()
	}

	stmt.Body = body

	return stmt
}

// parseStateBlockStatement parses a state block statement
func (p *Parser) parseStateBlockStatement() *StateBlockStatement {
	stmt := &StateBlockStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionStatement parses a function statement
func (p *Parser) parseFunctionStatement() *FunctionStatement {
	stmt := &FunctionStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseParameters()

	// Check for return type
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // Skip the colon
		p.nextToken() // Move to the type
		stmt.ReturnType = p.parseTypeExpression()
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseConstructorStatement parses a constructor statement
func (p *Parser) parseConstructorStatement() *ConstructorStatement {
	stmt := &ConstructorStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseParameters()

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseEventStatement parses an event statement
func (p *Parser) parseEventStatement() *EventStatement {
	stmt := &EventStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseParameters()

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseRequireStatement parses a require statement
func (p *Parser) parseRequireStatement() *RequireStatement {
	stmt := &RequireStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COMMA) {
		return nil
	}

	p.nextToken()
	stmt.Message = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseEmitStatement parses an emit statement
func (p *Parser) parseEmitStatement() *EmitStatement {
	stmt := &EmitStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.EventName = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Parse arguments
	args := []Expression{}

	// Check if there are any arguments
	if !p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))

		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken()
			p.nextToken()
			args = append(args, p.parseExpression(LOWEST))
		}
	}

	stmt.Arguments = args

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement parses a block statement
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpression parses an expression
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseIdentifier parses an identifier
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral parses an integer literal
func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseBooleanLiteral parses a boolean literal
func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

// parsePrefixExpression parses a prefix expression
func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression parses an infix expression
func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression parses a grouped expression
func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression parses an if expression
func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// parseCallExpression parses a call expression
func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

// parseIndexExpression parses an index expression
func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

// parseExpressionList parses a list of expressions
func (p *Parser) parseExpressionList(end lexer.TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseParameters parses function parameters
func (p *Parser) parseParameters() []*ParameterStatement {
	parameters := []*ParameterStatement{}

	// Check if there are any parameters
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	param := &ParameterStatement{
		Token: p.curToken,
		Name:  &Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	param.Type = p.parseTypeExpression()
	parameters = append(parameters, param)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()

		param := &ParameterStatement{
			Token: p.curToken,
			Name:  &Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()
		param.Type = p.parseTypeExpression()
		parameters = append(parameters, param)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return parameters
}

// parseTypeExpression parses a type expression
func (p *Parser) parseTypeExpression() *TypeExpression {
	expr := &TypeExpression{Token: p.curToken}

	// Basic type (e.g., Int, String, Address)
	if p.curTokenIs(lexer.IDENT) {
		expr.Type = p.curToken.Literal
	} else if p.curTokenIs(lexer.ADDRESS) {
		expr.Type = "Address"
	} else if p.curTokenIs(lexer.MAP) {
		// Map type (e.g., Map<Address, Int>)
		expr.Type = "Map"

		if !p.expectPeek(lexer.LT) {
			return nil
		}

		p.nextToken()
		expr.KeyType = p.parseTypeExpression()

		if !p.expectPeek(lexer.COMMA) {
			return nil
		}

		p.nextToken()
		expr.ValueType = p.parseTypeExpression()

		if !p.expectPeek(lexer.GT) {
			return nil
		}
	} else {
		msg := fmt.Sprintf("expected type expression, got %s instead", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	return expr
}

// curTokenIs checks if the current token is of the given type
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if the next token is of the given type
func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the next token is of the given type and advances if it is
func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekPrecedence returns the precedence of the next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence returns the precedence of the current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// peekError adds an error for an unexpected token
func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// noPrefixParseFnError adds an error for a token that doesn't have a prefix parse function
func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// registerPrefix registers a prefix parse function
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseDotExpression parses a dot expression (e.g., obj.property)
func (p *Parser) parseDotExpression(left Expression) Expression {
	expression := &DotExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	if !p.curTokenIs(lexer.IDENT) {
		msg := fmt.Sprintf("expected identifier after dot, got %s instead", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	expression.Right = p.parseIdentifier()

	return expression
}
