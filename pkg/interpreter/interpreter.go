package interpreter

import (
	"fmt"
	"github.com/Stremax-Team/stremax-lang/pkg/blockchain"
	"github.com/Stremax-Team/stremax-lang/pkg/errors"
	"github.com/Stremax-Team/stremax-lang/pkg/lexer"
	"github.com/Stremax-Team/stremax-lang/pkg/parser"
)

// Object represents a runtime value
type Object interface {
	Type() string
	Inspect() string
}

// Integer represents an integer value
type Integer struct {
	Value int64
}

// Type returns the type of the Integer object
func (i *Integer) Type() string { return "INTEGER" }

// Inspect returns a string representation of the Integer object
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// String represents a string value
type String struct {
	Value string
}

// Type returns the type of the String object
func (s *String) Type() string { return "STRING" }

// Inspect returns a string representation of the String object
func (s *String) Inspect() string { return s.Value }

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

// Type returns the type of the Boolean object
func (b *Boolean) Type() string { return "BOOLEAN" }

// Inspect returns a string representation of the Boolean object
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Address represents a blockchain address
type Address struct {
	Value blockchain.Address
}

// Type returns the type of the Address object
func (a *Address) Type() string { return "ADDRESS" }

// Inspect returns a string representation of the Address object
func (a *Address) Inspect() string { return string(a.Value) }

// Environment represents a variable environment
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

// NewEnclosedEnvironment creates a new environment with an outer environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a variable from the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a variable in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Interpreter represents an interpreter for Stremax-Lang
type Interpreter struct {
	source string
	lexer  *lexer.Lexer
	parser *parser.Parser
	env    *Environment
	bc     *blockchain.Blockchain
}

// New creates a new interpreter
func New(source string) *Interpreter {
	l := lexer.New(source)
	p := parser.New(l)
	
	return &Interpreter{
		source: source,
		lexer:  l,
		parser: p,
		env:    NewEnvironment(),
		bc:     blockchain.New(),
	}
}

// Run executes the source code
func (i *Interpreter) Run() error {
	// Parse the program
	program := i.parser.ParseProgram()
	if len(i.parser.Errors()) != 0 {
		for _, msg := range i.parser.Errors() {
			fmt.Printf("Parser error: %s\n", msg)
		}
		return errors.NewSyntaxError("Failed to parse program", 0, 0, "")
	}
	
	// Evaluate the program
	result, err := i.evalProgram(program)
	if err != nil {
		return err
	}
	
	if result != nil {
		fmt.Printf("Result: %s\n", result.Inspect())
	}
	
	return nil
}

// evalProgram evaluates a program
func (i *Interpreter) evalProgram(program *parser.Program) (Object, error) {
	var result Object
	var err error
	
	for _, stmt := range program.Statements {
		result, err = i.evalStatement(stmt)
		if err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

// evalStatement evaluates a statement
func (i *Interpreter) evalStatement(stmt parser.Statement) (Object, error) {
	switch s := stmt.(type) {
	case *parser.LetStatement:
		return i.evalLetStatement(s)
	case *parser.ReturnStatement:
		return i.evalReturnStatement(s)
	case *parser.ExpressionStatement:
		return i.evalExpression(s.Expression)
	case *parser.BlockStatement:
		return i.evalBlockStatement(s)
	case *parser.ContractStatement:
		return i.evalContractStatement(s)
	case *parser.FunctionStatement:
		return i.evalFunctionStatement(s)
	case *parser.RequireStatement:
		return i.evalRequireStatement(s)
	case *parser.EmitStatement:
		return i.evalEmitStatement(s)
	default:
		return nil, errors.NewRuntimeError("Unknown statement type", 0, 0, "")
	}
}

// evalLetStatement evaluates a let statement
func (i *Interpreter) evalLetStatement(stmt *parser.LetStatement) (Object, error) {
	val, err := i.evalExpression(stmt.Value)
	if err != nil {
		return nil, err
	}
	
	i.env.Set(stmt.Name.Value, val)
	return val, nil
}

// evalReturnStatement evaluates a return statement
func (i *Interpreter) evalReturnStatement(stmt *parser.ReturnStatement) (Object, error) {
	return i.evalExpression(stmt.ReturnValue)
}

// evalBlockStatement evaluates a block statement
func (i *Interpreter) evalBlockStatement(block *parser.BlockStatement) (Object, error) {
	var result Object
	var err error
	
	for _, stmt := range block.Statements {
		result, err = i.evalStatement(stmt)
		if err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

// evalExpression evaluates an expression
func (i *Interpreter) evalExpression(expr parser.Expression) (Object, error) {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return &Integer{Value: e.Value}, nil
	case *parser.StringLiteral:
		return &String{Value: e.Value}, nil
	case *parser.BooleanLiteral:
		return &Boolean{Value: e.Value}, nil
	case *parser.Identifier:
		return i.evalIdentifier(e)
	case *parser.PrefixExpression:
		return i.evalPrefixExpression(e)
	case *parser.InfixExpression:
		return i.evalInfixExpression(e)
	case *parser.IfExpression:
		return i.evalIfExpression(e)
	case *parser.CallExpression:
		return i.evalCallExpression(e)
	case *parser.IndexExpression:
		return i.evalIndexExpression(e)
	case *parser.DotExpression:
		return i.evalDotExpression(e)
	default:
		return nil, errors.NewRuntimeError("Unknown expression type", 0, 0, "")
	}
}

// evalIdentifier evaluates an identifier
func (i *Interpreter) evalIdentifier(ident *parser.Identifier) (Object, error) {
	val, ok := i.env.Get(ident.Value)
	if !ok {
		return nil, errors.NewReferenceError(fmt.Sprintf("Identifier not found: %s", ident.Value), ident.Token.Line, ident.Token.Column, "")
	}
	return val, nil
}

// evalPrefixExpression evaluates a prefix expression
func (i *Interpreter) evalPrefixExpression(expr *parser.PrefixExpression) (Object, error) {
	right, err := i.evalExpression(expr.Right)
	if err != nil {
		return nil, err
	}
	
	switch expr.Operator {
	case "!":
		return i.evalBangOperatorExpression(right)
	case "-":
		return i.evalMinusPrefixOperatorExpression(right)
	default:
		return nil, errors.NewRuntimeError(fmt.Sprintf("Unknown operator: %s", expr.Operator), expr.Token.Line, expr.Token.Column, "")
	}
}

// evalBangOperatorExpression evaluates a bang operator expression
func (i *Interpreter) evalBangOperatorExpression(right Object) (Object, error) {
	switch right := right.(type) {
	case *Boolean:
		return &Boolean{Value: !right.Value}, nil
	default:
		return &Boolean{Value: false}, nil
	}
}

// evalMinusPrefixOperatorExpression evaluates a minus prefix operator expression
func (i *Interpreter) evalMinusPrefixOperatorExpression(right Object) (Object, error) {
	if right.Type() != "INTEGER" {
		return nil, errors.NewTypeError("Cannot negate non-integer", 0, 0, "")
	}
	
	value := right.(*Integer).Value
	return &Integer{Value: -value}, nil
}

// evalInfixExpression evaluates an infix expression
func (i *Interpreter) evalInfixExpression(expr *parser.InfixExpression) (Object, error) {
	left, err := i.evalExpression(expr.Left)
	if err != nil {
		return nil, err
	}
	
	right, err := i.evalExpression(expr.Right)
	if err != nil {
		return nil, err
	}
	
	switch {
	case left.Type() == "INTEGER" && right.Type() == "INTEGER":
		return i.evalIntegerInfixExpression(expr.Operator, left, right)
	case left.Type() == "STRING" && right.Type() == "STRING":
		return i.evalStringInfixExpression(expr.Operator, left, right)
	case expr.Operator == "==":
		return &Boolean{Value: left == right}, nil
	case expr.Operator == "!=":
		return &Boolean{Value: left != right}, nil
	default:
		return nil, errors.NewTypeError(
			fmt.Sprintf("Type mismatch: %s %s %s", left.Type(), expr.Operator, right.Type()),
			expr.Token.Line,
			expr.Token.Column,
			"",
		)
	}
}

// evalIntegerInfixExpression evaluates an integer infix expression
func (i *Interpreter) evalIntegerInfixExpression(operator string, left, right Object) (Object, error) {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value
	
	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}, nil
	case "-":
		return &Integer{Value: leftVal - rightVal}, nil
	case "*":
		return &Integer{Value: leftVal * rightVal}, nil
	case "/":
		if rightVal == 0 {
			return nil, errors.NewRuntimeError("Division by zero", 0, 0, "")
		}
		return &Integer{Value: leftVal / rightVal}, nil
	case "<":
		return &Boolean{Value: leftVal < rightVal}, nil
	case ">":
		return &Boolean{Value: leftVal > rightVal}, nil
	case "==":
		return &Boolean{Value: leftVal == rightVal}, nil
	case "!=":
		return &Boolean{Value: leftVal != rightVal}, nil
	default:
		return nil, errors.NewRuntimeError(fmt.Sprintf("Unknown operator: %s", operator), 0, 0, "")
	}
}

// evalStringInfixExpression evaluates a string infix expression
func (i *Interpreter) evalStringInfixExpression(operator string, left, right Object) (Object, error) {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	
	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}, nil
	case "==":
		return &Boolean{Value: leftVal == rightVal}, nil
	case "!=":
		return &Boolean{Value: leftVal != rightVal}, nil
	default:
		return nil, errors.NewRuntimeError(fmt.Sprintf("Unknown operator: %s", operator), 0, 0, "")
	}
}

// evalIfExpression evaluates an if expression
func (i *Interpreter) evalIfExpression(expr *parser.IfExpression) (Object, error) {
	condition, err := i.evalExpression(expr.Condition)
	if err != nil {
		return nil, err
	}
	
	if isTruthy(condition) {
		return i.evalBlockStatement(expr.Consequence)
	} else if expr.Alternative != nil {
		return i.evalBlockStatement(expr.Alternative)
	} else {
		return nil, nil
	}
}

// isTruthy determines if an object is truthy
func isTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case *Boolean:
		return obj.Value
	case *Integer:
		return obj.Value != 0
	default:
		return true
	}
}

// evalCallExpression evaluates a call expression
func (i *Interpreter) evalCallExpression(expr *parser.CallExpression) (Object, error) {
	// For now, just return nil
	return nil, errors.NewRuntimeError("Function calls not implemented yet", expr.Token.Line, expr.Token.Column, "")
}

// evalIndexExpression evaluates an index expression
func (i *Interpreter) evalIndexExpression(expr *parser.IndexExpression) (Object, error) {
	// For now, just return nil
	return nil, errors.NewRuntimeError("Index expressions not implemented yet", expr.Token.Line, expr.Token.Column, "")
}

// evalDotExpression evaluates a dot expression
func (i *Interpreter) evalDotExpression(expr *parser.DotExpression) (Object, error) {
	// For now, just return nil
	return nil, errors.NewRuntimeError("Dot expressions not implemented yet", expr.Token.Line, expr.Token.Column, "")
}

// evalContractStatement evaluates a contract statement
func (i *Interpreter) evalContractStatement(stmt *parser.ContractStatement) (Object, error) {
	// For now, just return nil
	return nil, errors.NewRuntimeError("Contract statements not implemented yet", stmt.Token.Line, stmt.Token.Column, "")
}

// evalFunctionStatement evaluates a function statement
func (i *Interpreter) evalFunctionStatement(stmt *parser.FunctionStatement) (Object, error) {
	// For now, just return nil
	return nil, errors.NewRuntimeError("Function statements not implemented yet", stmt.Token.Line, stmt.Token.Column, "")
}

// evalRequireStatement evaluates a require statement
func (i *Interpreter) evalRequireStatement(stmt *parser.RequireStatement) (Object, error) {
	condition, err := i.evalExpression(stmt.Condition)
	if err != nil {
		return nil, err
	}
	
	if !isTruthy(condition) {
		message := "Requirement failed"
		if stmt.Message != nil {
			msgObj, err := i.evalExpression(stmt.Message)
			if err != nil {
				return nil, err
			}
			
			if msgObj.Type() == "STRING" {
				message = msgObj.(*String).Value
			}
		}
		return nil, errors.NewRuntimeError(message, stmt.Token.Line, stmt.Token.Column, "")
	}
	
	return nil, nil
}

// evalEmitStatement evaluates an emit statement
func (i *Interpreter) evalEmitStatement(stmt *parser.EmitStatement) (Object, error) {
	// For now, just print the event
	fmt.Printf("Event emitted: %s\n", stmt.EventName.Value)
	
	for _, arg := range stmt.Arguments {
		argObj, err := i.evalExpression(arg)
		if err != nil {
			return nil, err
		}
		fmt.Printf("  Argument: %s\n", argObj.Inspect())
	}
	
	return nil, nil
} 