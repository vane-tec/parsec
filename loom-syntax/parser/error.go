package parser

import (
	"fmt"

	"github.com/vanetec/loom-syntax/syntax"
)

// ParseError représente une erreur de parsing avec localisation
type ParseError struct {
	message string
	line    int
	column  int
}

func (this ParseError) Error() string {
	return fmt.Sprintf("%s at %d:%d", this.message, this.line, this.column)
}

// Helper pour générer une erreur depuis un token.
func errorFromToken(t syntax.Token, m string) error {
	return ParseError{
		message: m,
		line:    t.Line(),
		column:  t.Column(),
	}
}
