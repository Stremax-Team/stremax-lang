package interpreter

import (
	"bytes"
	"fmt"
	"github.com/Stremax-Team/stremax-lang/pkg/blockchain"
	"github.com/Stremax-Team/stremax-lang/pkg/errors"
	"github.com/Stremax-Team/stremax-lang/pkg/lexer"
	"github.com/Stremax-Team/stremax-lang/pkg/parser"
	"strings"
	"hash/fnv"
)

// Object represents a runtime value in the Stremax-Lang interpreter.
// All values in Stremax-Lang implement this interface, which provides
// methods for type identification and string representation.
type Object interface {
	// Type returns the type of the object as a string.
	Type() string

	// Inspect returns a string representation of the object.
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

// Function represents a function definition
type Function struct {
	Parameters []*parser.ParameterStatement
	Body       *parser.BlockStatement
	ReturnType *parser.TypeExpression
	Env        *Environment
	Name       string // Optional, for named functions
}

// Type returns the type of the Function object
func (f *Function) Type() string { return "FUNCTION" }

// Inspect returns a string representation of the Function object
func (f *Function) Inspect() string {
	return fmt.Sprintf("function %s", f.Name)
}

// ReturnValue represents a value returned from a function
type ReturnValue struct {
	Value Object
}

// Type returns the type of the ReturnValue object
func (rv *ReturnValue) Type() string { return "RETURN_VALUE" }

// Inspect returns a string representation of the ReturnValue object
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Array represents an array object
type Array struct {
	Elements []Object
}

// Type returns the type of the Array object
func (a *Array) Type() string { return "ARRAY" }

// Inspect returns a string representation of the Array object
func (a *Array) Inspect() string {
	var out bytes.Buffer
	
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	
	return out.String()
}

// HashKey represents a key in a hash map
type HashKey struct {
	Type  string
	Value uint64
}

// Hashable is an interface for objects that can be used as hash keys
type Hashable interface {
	HashKey() HashKey
}

// HashPair represents a key-value pair in a hash map
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents a hash map object
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the type of the Hash object
func (h *Hash) Type() string { return "HASH" }

// Inspect returns a string representation of the Hash object
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	
	return out.String()
}

// StringHashKey makes String hashable
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// BooleanHashKey makes Boolean hashable
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	
	return HashKey{Type: b.Type(), Value: value}
}

// IntegerHashKey makes Integer hashable
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

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

// Interpreter represents an interpreter for Stremax-Lang.
// It handles lexing, parsing, and evaluating Stremax-Lang code,
// maintaining the execution environment and blockchain state.
type Interpreter struct {
	source string
	lexer  *lexer.Lexer
	parser *parser.Parser
	env    *Environment
	bc     *blockchain.Blockchain
}

// New creates a new Stremax-Lang interpreter with the given source code.
// It initializes the lexer, parser, environment, and blockchain components
// needed for execution.
//
// Parameters:
//   - source: The Stremax-Lang source code to interpret
//
// Returns:
//   - A new Interpreter instance ready to execute the provided code
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

// Run executes the Stremax-Lang source code provided to the interpreter.
// It parses the program, evaluates it, and returns any errors encountered
// during execution.
//
// Returns:
//   - An error if parsing or evaluation fails, nil otherwise
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
	value, err := i.evalExpression(stmt.ReturnValue)
	if err != nil {
		return nil, err
	}
	
	// Wrap the return value
	return &ReturnValue{Value: value}, nil
}

// evalBlockStatement evaluates a block statement
func (i *Interpreter) evalBlockStatement(block *parser.BlockStatement) (Object, error) {
	var result Object
	var err error

	// Create a new enclosed environment for this block
	previousEnv := i.env
	blockEnv := NewEnclosedEnvironment(i.env)
	i.env = blockEnv

	// Evaluate statements in the block environment
	for _, stmt := range block.Statements {
		result, err = i.evalStatement(stmt)
		if err != nil {
			i.env = previousEnv // Restore the previous environment in case of error
			return nil, err
		}
		
		// Check if it's a return value, if so, return early
		if result != nil && result.Type() == "RETURN_VALUE" {
			i.env = previousEnv // Restore the previous environment
			return result, nil
		}
	}

	// Restore the previous environment
	i.env = previousEnv

	return result, nil
}

// NULL represents a null value
var NULL = &Null{}

// Null represents a null value
type Null struct{}

// Type returns the type of the Null object
func (n *Null) Type() string { return "NULL" }

// Inspect returns a string representation of the Null object
func (n *Null) Inspect() string { return "null" }

// evalExpression evaluates an expression and returns the result
func (i *Interpreter) evalExpression(expr parser.Expression) (Object, error) {
	// Wrap in a type switch to handle different expression types
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return &Integer{Value: e.Value}, nil
	case *parser.StringLiteral:
		return &String{Value: e.Value}, nil
	case *parser.BooleanLiteral:
		return &Boolean{Value: e.Value}, nil
	case *parser.PrefixExpression:
		return i.evalPrefixExpression(e)
	case *parser.InfixExpression:
		return i.evalInfixExpression(e)
	case *parser.IfExpression:
		return i.evalIfExpression(e)
	case *parser.Identifier:
		return i.evalIdentifier(e)
	case *parser.CallExpression:
		return i.evalCallExpression(e)
	case *parser.FunctionLiteral:
		return i.evalFunctionLiteral(e)
	case *parser.ArrayLiteral:
		return i.evalArrayLiteral(e)
	case *parser.IndexExpression:
		return i.evalIndexExpression(e)
	case *parser.HashLiteral:
		return i.evalHashLiteral(e)
	default:
		fmt.Printf("DEBUG: Unknown expression type: %T\n", e)
		return nil, errors.NewRuntimeError(fmt.Sprintf("Unknown expression type: %T", e), 0, 0, "")
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
		// For "+" operator, support both addition and concatenation
		if expr.Operator == "+" {
			// If either integer is very large or has a special format, treat as concatenation
			leftVal := left.(*Integer).Value
			rightVal := right.(*Integer).Value
			leftStr := fmt.Sprintf("%d", leftVal)
			rightStr := fmt.Sprintf("%d", rightVal)
			
			// Check if this should be treated as concatenation
			if strings.HasPrefix(leftStr, "0") || strings.HasPrefix(rightStr, "0") {
				// Handle as string concatenation when numbers have leading zeros
				return &String{Value: leftStr + rightStr}, nil
			} else {
				// Regular integer addition
				return i.evalIntegerInfixExpression(expr.Operator, left, right)
			}
		} else {
			// Other operators proceed with normal integer operations
			return i.evalIntegerInfixExpression(expr.Operator, left, right)
		}
	case left.Type() == "STRING" && right.Type() == "STRING":
		return i.evalStringInfixExpression(expr.Operator, left, right)
	// Support string concatenation with other types
	case left.Type() == "STRING" && expr.Operator == "+":
		return i.evalMixedStringConcatExpression(left, right, true)
	case right.Type() == "STRING" && expr.Operator == "+":
		return i.evalMixedStringConcatExpression(right, left, false)
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

// evalLogicalExpression evaluates a logical expression with short-circuit evaluation
func (i *Interpreter) evalLogicalExpression(expr *parser.InfixExpression) (Object, error) {
	// Evaluate the left operand
	left, err := i.evalExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	// Check if left operand is a boolean
	if left.Type() != "BOOLEAN" {
		return nil, errors.NewTypeError(
			fmt.Sprintf("Left operand of %s must be a boolean, got %s", expr.Operator, left.Type()),
			expr.Token.Line,
			expr.Token.Column,
			"",
		)
	}

	leftBool := left.(*Boolean).Value

	// Short-circuit evaluation
	if expr.Operator == "&&" {
		// If left is false, return false without evaluating right
		if !leftBool {
			return &Boolean{Value: false}, nil
		}
	} else if expr.Operator == "||" {
		// If left is true, return true without evaluating right
		if leftBool {
			return &Boolean{Value: true}, nil
		}
	}

	// Evaluate the right operand
	right, err := i.evalExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	// Check if right operand is a boolean
	if right.Type() != "BOOLEAN" {
		return nil, errors.NewTypeError(
			fmt.Sprintf("Right operand of %s must be a boolean, got %s", expr.Operator, right.Type()),
			expr.Token.Line,
			expr.Token.Column,
			"",
		)
	}

	rightBool := right.(*Boolean).Value

	// Determine the result based on the operator
	if expr.Operator == "&&" {
		return &Boolean{Value: leftBool && rightBool}, nil
	} else { // expr.Operator == "||"
		return &Boolean{Value: leftBool || rightBool}, nil
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
	case "<=":
		return &Boolean{Value: leftVal <= rightVal}, nil
	case ">=":
		return &Boolean{Value: leftVal >= rightVal}, nil
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

// evalStringConcatExpression evaluates string concatenation with any other type
func (i *Interpreter) evalStringConcatExpression(strObj Object, otherObj Object) (Object, error) {
	strVal := strObj.(*String).Value
	
	// Convert the other object to a string
	var otherVal string
	switch other := otherObj.(type) {
	case *Integer:
		otherVal = fmt.Sprintf("%d", other.Value)
	case *Boolean:
		otherVal = fmt.Sprintf("%t", other.Value)
	case *String:
		otherVal = other.Value
	case *Address:
		otherVal = string(other.Value)
	default:
		// For any other type, use the Inspect method
		otherVal = other.Inspect()
	}
	
	return &String{Value: strVal + otherVal}, nil
}

// evalMixedStringConcatExpression handles both string+other and other+string cases
func (i *Interpreter) evalMixedStringConcatExpression(strObj Object, otherObj Object, stringIsLeft bool) (Object, error) {
	strVal := strObj.(*String).Value
	
	// Convert the other object to a string
	var otherVal string
	switch other := otherObj.(type) {
	case *Integer:
		otherVal = fmt.Sprintf("%d", other.Value)
	case *Boolean:
		otherVal = fmt.Sprintf("%t", other.Value)
	case *String:
		otherVal = other.Value
	case *Address:
		otherVal = string(other.Value)
	default:
		// For any other type, use the Inspect method
		otherVal = other.Inspect()
	}
	
	if stringIsLeft {
		return &String{Value: strVal + otherVal}, nil
	} else {
		return &String{Value: otherVal + strVal}, nil
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
	// Evaluate the function expression to get the function object
	function, err := i.evalExpression(expr.Function)
	if err != nil {
		return nil, err
	}
	
	// Check if it's actually a function
	fn, ok := function.(*Function)
	if !ok {
		return nil, errors.NewTypeError(
			fmt.Sprintf("Not a function: %s", function.Type()),
			expr.Token.Line,
			expr.Token.Column,
			"",
		)
	}
	
	// Evaluate the arguments
	args, err := i.evalExpressions(expr.Arguments)
	if err != nil {
		return nil, err
	}
	
	// Check if the number of arguments matches the number of parameters
	if len(args) != len(fn.Parameters) {
		return nil, errors.NewTypeError(
			fmt.Sprintf("Wrong number of arguments: expected %d, got %d",
				len(fn.Parameters), len(args)),
			expr.Token.Line,
			expr.Token.Column,
			"",
		)
	}
	
	// Create a new environment for the function call
	extendedEnv := NewEnclosedEnvironment(fn.Env)
	
	// Bind the arguments to the parameters
	for i, param := range fn.Parameters {
		extendedEnv.Set(param.Name.Value, args[i])
	}
	
	// Save the current environment and set the function's environment
	previousEnv := i.env
	i.env = extendedEnv
	
	// Evaluate the function body
	result, err := i.evalBlockStatement(fn.Body)
	
	// Restore the previous environment
	i.env = previousEnv
	
	if err != nil {
		return nil, err
	}
	
	// Unwrap the return value if it's a return value
	if returnValue, ok := result.(*ReturnValue); ok {
		return returnValue.Value, nil
	}
	
	return result, nil
}

// evalExpressions evaluates a list of expressions
func (i *Interpreter) evalExpressions(exps []parser.Expression) ([]Object, error) {
	var result []Object
	
	for _, exp := range exps {
		evaluated, err := i.evalExpression(exp)
		if err != nil {
			return nil, err
		}
		result = append(result, evaluated)
	}
	
	return result, nil
}

// evalIndexExpression evaluates an index expression
func (i *Interpreter) evalIndexExpression(expr *parser.IndexExpression) (Object, error) {
	left, err := i.evalExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	index, err := i.evalExpression(expr.Index)
	if err != nil {
		return nil, err
	}

	return i.evalElementAccess(left, index, expr.Token)
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
	// Create a name for anonymous functions if necessary
	name := ""
	if stmt.Name != nil {
		name = stmt.Name.Value
	}
	
	function := &Function{
		Parameters: stmt.Parameters,
		Body:       stmt.Body,
		ReturnType: stmt.ReturnType,
		Env:        i.env,
		Name:       name,
	}
	
	// Store the function in the current environment if it has a name
	if name != "" {
		i.env.Set(name, function)
	}
	
	return function, nil
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

// evalFunctionLiteral evaluates a function literal expression
func (i *Interpreter) evalFunctionLiteral(fl *parser.FunctionLiteral) (Object, error) {
	function := &Function{
		Parameters: fl.Parameters,
		Body:       fl.Body,
		ReturnType: fl.ReturnType,
		Env:        i.env, // Capture the current environment for closures
	}
	
	return function, nil
}

// evalArrayLiteral evaluates an array literal expression
func (i *Interpreter) evalArrayLiteral(node *parser.ArrayLiteral) (Object, error) {
	elements := []Object{}

	for _, el := range node.Elements {
		evaluated, err := i.evalExpression(el)
		if err != nil {
			return nil, err
		}
		elements = append(elements, evaluated)
	}

	return &Array{Elements: elements}, nil
}

// evalElementAccess handles accessing elements from arrays
func (i *Interpreter) evalElementAccess(left, index Object, token parser.Token) (Object, error) {
	switch {
	case left.Type() == "ARRAY" && index.Type() == "INTEGER":
		return i.evalArrayIndexExpression(left, index, token)
	case left.Type() == "HASH":
		return i.evalHashIndexExpression(left, index, token)
	default:
		return nil, errors.NewRuntimeError(
			fmt.Sprintf("index operator not supported: %s", left.Type()),
			token.Line, token.Column, "")
	}
}

// evalArrayIndexExpression implements array indexing
func (i *Interpreter) evalArrayIndexExpression(array, index Object, token parser.Token) (Object, error) {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL, nil
	}

	return arrayObject.Elements[idx], nil
}

// evalHashLiteral evaluates a hash literal expression
func (i *Interpreter) evalHashLiteral(node *parser.HashLiteral) (Object, error) {
	pairs := make(map[HashKey]HashPair)

	for keyNode, valueNode := range node.Pairs {
		key, err := i.evalExpression(keyNode)
		if err != nil {
			return nil, err
		}
		
		hashKey, ok := key.(Hashable)
		if !ok {
			return nil, errors.NewRuntimeError(
				fmt.Sprintf("unusable as hash key: %s", key.Type()),
				node.Token.Line, node.Token.Column, "")
		}

		value, err := i.evalExpression(valueNode)
		if err != nil {
			return nil, err
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key: key, Value: value}
	}

	return &Hash{Pairs: pairs}, nil
}

// evalHashIndexExpression handles hash element access with [key]
func (i *Interpreter) evalHashIndexExpression(hash, index Object, token parser.Token) (Object, error) {
	hashObject := hash.(*Hash)
	
	key, ok := index.(Hashable)
	if !ok {
		return nil, errors.NewRuntimeError(
			fmt.Sprintf("unusable as hash key: %s", index.Type()),
			token.Line, token.Column, "")
	}
	
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL, nil
	}
	
	return pair.Value, nil
}
