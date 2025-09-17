package parser

import (
	"fmt"

	"github.com/usememos/gomark/parser/tokenizer"
)

// ParseError represents different types of parsing errors.
type ParseError struct {
	Type     ErrorType
	Message  string
	Position int
	Token    *tokenizer.Token
	Context  string
}

// ErrorType defines the category of parsing error.
type ErrorType string

const (
	// SyntaxError indicates malformed markdown syntax.
	SyntaxError ErrorType = "syntax_error"

	// UnexpectedToken indicates an unexpected token was encountered.
	UnexpectedToken ErrorType = "unexpected_token"

	// UnmatchedDelimiter indicates unmatched opening/closing delimiters.
	UnmatchedDelimiter ErrorType = "unmatched_delimiter"

	// InvalidNesting indicates invalid nesting of markdown elements.
	InvalidNesting ErrorType = "invalid_nesting"

	// UnknownElement indicates an unknown or unsupported markdown element.
	UnknownElement ErrorType = "unknown_element"

	// IncompleteElement indicates an incomplete markdown element.
	IncompleteElement ErrorType = "incomplete_element"
)

func (e *ParseError) Error() string {
	if e.Token != nil {
		return fmt.Sprintf("%s at position %d: %s (token: %s)", e.Type, e.Position, e.Message, e.Token.Value)
	}
	return fmt.Sprintf("%s at position %d: %s", e.Type, e.Position, e.Message)
}

// NewSyntaxError creates a syntax error.
func NewSyntaxError(message string, position int, token *tokenizer.Token) *ParseError {
	return &ParseError{
		Type:     SyntaxError,
		Message:  message,
		Position: position,
		Token:    token,
	}
}

// NewUnexpectedTokenError creates an unexpected token error.
func NewUnexpectedTokenError(expected string, actual *tokenizer.Token, position int) *ParseError {
	message := fmt.Sprintf("expected %s, got %s", expected, actual.Type)
	return &ParseError{
		Type:     UnexpectedToken,
		Message:  message,
		Position: position,
		Token:    actual,
	}
}

// NewUnmatchedDelimiterError creates an unmatched delimiter error.
func NewUnmatchedDelimiterError(delimiter string, position int) *ParseError {
	return &ParseError{
		Type:     UnmatchedDelimiter,
		Message:  fmt.Sprintf("unmatched %s", delimiter),
		Position: position,
	}
}

// NewInvalidNestingError creates an invalid nesting error.
func NewInvalidNestingError(parent, child string, position int) *ParseError {
	return &ParseError{
		Type:     InvalidNesting,
		Message:  fmt.Sprintf("cannot nest %s inside %s", child, parent),
		Position: position,
	}
}

// NewUnknownElementError creates an unknown element error.
func NewUnknownElementError(element string, position int, token *tokenizer.Token) *ParseError {
	return &ParseError{
		Type:     UnknownElement,
		Message:  fmt.Sprintf("unknown markdown element: %s", element),
		Position: position,
		Token:    token,
	}
}

// NewIncompleteElementError creates an incomplete element error.
func NewIncompleteElementError(element string, position int) *ParseError {
	return &ParseError{
		Type:     IncompleteElement,
		Message:  fmt.Sprintf("incomplete %s element", element),
		Position: position,
	}
}

// ParseResult wraps parsing results with error information.
type ParseResult struct {
	Success  bool
	Errors   []*ParseError
	Warnings []*ParseError
}

// AddError adds an error to the result.
func (r *ParseResult) AddError(err *ParseError) {
	r.Errors = append(r.Errors, err)
	r.Success = false
}

// AddWarning adds a warning to the result.
func (r *ParseResult) AddWarning(warning *ParseError) {
	r.Warnings = append(r.Warnings, warning)
}

// HasErrors returns true if there are any errors.
func (r *ParseResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// HasWarnings returns true if there are any warnings.
func (r *ParseResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// FirstError returns the first error, or nil if none.
func (r *ParseResult) FirstError() error {
	if len(r.Errors) > 0 {
		return r.Errors[0]
	}
	return nil
}

// NewParseResult creates a new ParseResult.
func NewParseResult() *ParseResult {
	return &ParseResult{
		Success:  true,
		Errors:   []*ParseError{},
		Warnings: []*ParseError{},
	}
}
