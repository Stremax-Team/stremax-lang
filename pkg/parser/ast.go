package parser

import (
	"bytes"
	"strings"
)

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node in the AST
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node in the AST
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal of the token associated with the node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a string representation of the program
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ContractStatement represents a contract declaration
type ContractStatement struct {
	Token      Token // the 'contract' token
	Name       *Identifier
	StateBlock *StateBlockStatement
	Body       *BlockStatement
}

func (cs *ContractStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (cs *ContractStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// String returns a string representation of the contract statement
func (cs *ContractStatement) String() string {
	var out bytes.Buffer

	out.WriteString("contract ")
	out.WriteString(cs.Name.String())
	out.WriteString(" ")
	out.WriteString(cs.Body.String())

	return out.String()
}

// StateBlockStatement represents a state block in a contract
type StateBlockStatement struct {
	Token Token // the 'state' token
	Body  *BlockStatement
}

func (sb *StateBlockStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (sb *StateBlockStatement) TokenLiteral() string {
	return sb.Token.Literal
}

// String returns a string representation of the state block
func (sb *StateBlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("state ")
	out.WriteString(sb.Body.String())

	return out.String()
}

// FunctionStatement represents a function declaration
type FunctionStatement struct {
	Token      Token // the 'function' token
	Name       *Identifier
	Parameters []*ParameterStatement
	ReturnType *TypeExpression
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (fs *FunctionStatement) TokenLiteral() string {
	return fs.Token.Literal
}

// String returns a string representation of the function statement
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	if fs.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(fs.ReturnType.String())
		out.WriteString(" ")
	}

	out.WriteString(fs.Body.String())

	return out.String()
}

// ConstructorStatement represents a constructor declaration
type ConstructorStatement struct {
	Token      Token // the 'constructor' token
	Parameters []*ParameterStatement
	Body       *BlockStatement
}

func (cs *ConstructorStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (cs *ConstructorStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// String returns a string representation of the constructor statement
func (cs *ConstructorStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range cs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(cs.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(cs.Body.String())

	return out.String()
}

// EventStatement represents an event declaration
type EventStatement struct {
	Token      Token // the 'event' token
	Name       *Identifier
	Parameters []*ParameterStatement
}

func (es *EventStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (es *EventStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String returns a string representation of the event statement
func (es *EventStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range es.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(es.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(es.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	return out.String()
}

// ParameterStatement represents a parameter in a function or constructor
type ParameterStatement struct {
	Token Token // the parameter name token
	Name  *Identifier
	Type  *TypeExpression
}

func (ps *ParameterStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ps *ParameterStatement) TokenLiteral() string {
	return ps.Token.Literal
}

// String returns a string representation of the parameter statement
func (ps *ParameterStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ps.Name.String())
	out.WriteString(": ")
	out.WriteString(ps.Type.String())

	return out.String()
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Token      Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String returns a string representation of the block statement
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")

	return out.String()
}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	Token      Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String returns a string representation of the expression statement
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// LetStatement represents a variable declaration
type LetStatement struct {
	Token Token // the 'let' token
	Name  *Identifier
	Type  *TypeExpression
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String returns a string representation of the let statement
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())

	if ls.Type != nil {
		out.WriteString(": ")
		out.WriteString(ls.Type.String())
	}

	if ls.Value != nil {
		out.WriteString(" = ")
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Token       Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String returns a string representation of the return statement
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// RequireStatement represents a require statement
type RequireStatement struct {
	Token     Token      // the 'require' token
	Condition Expression // the condition to check
	Message   Expression // the error message (optional)
}

func (rs *RequireStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (rs *RequireStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String returns a string representation of the require statement
func (rs *RequireStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString("(")
	out.WriteString(rs.Condition.String())

	if rs.Message != nil {
		out.WriteString(", ")
		out.WriteString(rs.Message.String())
	}

	out.WriteString(")")

	return out.String()
}

// EmitStatement represents an emit statement
type EmitStatement struct {
	Token     Token        // the 'emit' token
	EventName *Identifier  // the event name
	Arguments []Expression // the event arguments
}

func (es *EmitStatement) statementNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (es *EmitStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String returns a string representation of the emit statement
func (es *EmitStatement) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range es.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(es.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(es.EventName.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// Identifier represents an identifier
type Identifier struct {
	Token Token // the identifier token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of the identifier
func (i *Identifier) String() string {
	return i.Value
}

// TypeExpression represents a type expression
type TypeExpression struct {
	Token     Token           // the type token
	Type      string          // the type name (e.g., Int, String, Address, Map)
	KeyType   *TypeExpression // for Map types, the key type
	ValueType *TypeExpression // for Map types, the value type
}

func (te *TypeExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (te *TypeExpression) TokenLiteral() string {
	return te.Token.Literal
}

// String returns a string representation of the type expression
func (te *TypeExpression) String() string {
	var out bytes.Buffer

	out.WriteString(te.Type)

	// If it's a Map type, include the key and value types
	if te.Type == "Map" && te.KeyType != nil && te.ValueType != nil {
		out.WriteString("<")
		out.WriteString(te.KeyType.String())
		out.WriteString(", ")
		out.WriteString(te.ValueType.String())
		out.WriteString(">")
	}

	return out.String()
}

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Token Token // the integer token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String returns a string representation of the integer literal
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// StringLiteral represents a string literal
type StringLiteral struct {
	Token Token // the string token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String returns a string representation of the string literal
func (sl *StringLiteral) String() string {
	return "\"" + sl.Value + "\""
}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Token Token // the boolean token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

// String returns a string representation of the boolean literal
func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

// PrefixExpression represents a prefix expression
type PrefixExpression struct {
	Token    Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String returns a string representation of the prefix expression
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an infix expression
type InfixExpression struct {
	Token    Token // the operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String returns a string representation of the infix expression
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// AssignExpression represents an assignment expression
type AssignExpression struct {
	Token    Token // the = token
	Left     Expression
	Operator string
	Right    Expression
}

func (ae *AssignExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ae *AssignExpression) TokenLiteral() string {
	return ae.Token.Literal
}

// String returns a string representation of the assignment expression
func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Left.String())
	out.WriteString(" " + ae.Operator + " ")
	out.WriteString(ae.Right.String())

	return out.String()
}

// IndexExpression represents an index expression (array[index])
type IndexExpression struct {
	Token Token // the [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String returns a string representation of the index expression
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// CallExpression represents a function call expression
type CallExpression struct {
	Token     Token // the ( token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String returns a string representation of the call expression
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// DotExpression represents a dot expression (e.g., obj.property)
type DotExpression struct {
	Token Token      // the '.' token
	Left  Expression // the expression on the left of the dot
	Right Expression // the identifier on the right of the dot
}

func (de *DotExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (de *DotExpression) TokenLiteral() string {
	return de.Token.Literal
}

// String returns a string representation of the dot expression
func (de *DotExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(de.Left.String())
	out.WriteString(".")
	out.WriteString(de.Right.String())
	out.WriteString(")")

	return out.String()
}

// IfExpression represents an if expression
type IfExpression struct {
	Token       Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral returns the literal of the token associated with the node
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String returns a string representation of the if expression
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
