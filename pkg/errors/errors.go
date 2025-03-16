package errors

import (
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// Error types
	SyntaxError    ErrorType = "SyntaxError"
	TypeError      ErrorType = "TypeError"
	ReferenceError ErrorType = "ReferenceError"
	RuntimeError   ErrorType = "RuntimeError"
	BlockchainError ErrorType = "BlockchainError"
	ContractError  ErrorType = "ContractError"
)

// Error represents a Stremax-Lang error
type Error struct {
	Type    ErrorType
	Message string
	Line    int
	Column  int
	File    string
}

// Error returns a string representation of the error
func (e *Error) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("%s: %s at %s:%d:%d", e.Type, e.Message, e.File, e.Line, e.Column)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewSyntaxError creates a new syntax error
func NewSyntaxError(message string, line, column int, file string) *Error {
	return &Error{
		Type:    SyntaxError,
		Message: message,
		Line:    line,
		Column:  column,
		File:    file,
	}
}

// NewTypeError creates a new type error
func NewTypeError(message string, line, column int, file string) *Error {
	return &Error{
		Type:    TypeError,
		Message: message,
		Line:    line,
		Column:  column,
		File:    file,
	}
}

// NewReferenceError creates a new reference error
func NewReferenceError(message string, line, column int, file string) *Error {
	return &Error{
		Type:    ReferenceError,
		Message: message,
		Line:    line,
		Column:  column,
		File:    file,
	}
}

// NewRuntimeError creates a new runtime error
func NewRuntimeError(message string, line, column int, file string) *Error {
	return &Error{
		Type:    RuntimeError,
		Message: message,
		Line:    line,
		Column:  column,
		File:    file,
	}
}

// NewBlockchainError creates a new blockchain error
func NewBlockchainError(message string) *Error {
	return &Error{
		Type:    BlockchainError,
		Message: message,
	}
}

// NewContractError creates a new contract error
func NewContractError(message string, contractAddress string) *Error {
	return &Error{
		Type:    ContractError,
		Message: fmt.Sprintf("%s (contract: %s)", message, contractAddress),
	}
}

// FormatErrorWithSource formats an error with the source code
func FormatErrorWithSource(err *Error, source string) string {
	if err.Line <= 0 || err.Column <= 0 {
		return err.Error()
	}

	// Split the source into lines
	lines := splitLines(source)
	if err.Line > len(lines) {
		return err.Error()
	}

	// Get the line with the error
	line := lines[err.Line-1]

	// Format the error message
	result := fmt.Sprintf("%s\n\n", err.Error())
	
	// Add the line with the error
	result += fmt.Sprintf("%4d | %s\n", err.Line, line)
	
	// Add a pointer to the column
	result += fmt.Sprintf("     | %s^\n", spaces(err.Column-1))
	
	return result
}

// splitLines splits a string into lines
func splitLines(s string) []string {
	var lines []string
	var line string
	
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(r)
		}
	}
	
	if line != "" {
		lines = append(lines, line)
	}
	
	return lines
}

// spaces returns a string with n spaces
func spaces(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
} 