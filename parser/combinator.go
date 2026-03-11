package parser

import (
	"bytes"

	"github.com/vanetec/loom-syntax/syntax"
)

// Parse exécute le parser sur un cursor
func (pw ParserWrapper[T]) Parse(c *Cursor) (T, error) {
	return pw.parser(c)
}

// FromParser convertit un Parser[T] classique en ParserWrapper[T]
func FromParser[T any](p Parser[T]) ParserWrapper[T] {
	return ParserWrapper[T]{parser: p}
}

// Expect ajoute un message d'erreur personnalisé
func (pw ParserWrapper[T]) Expect(message string) ParserWrapper[T] {
	wrapped := func(c *Cursor) (T, error) {
		val, err := pw.parser(c)
		if err != nil {
			tok, ok := c.Current()
			if ok {
				return val, errorFromToken(tok, message)
			}
			return val, ParseError{message, 0, 0}
		}
		return val, nil
	}
	return ParserWrapper[T]{parser: wrapped}
}

// Map transforme la valeur retournée par le parser
func Map[A any, B any](pw ParserWrapper[A], fn func(A) B) ParserWrapper[B] {
	return ParserWrapper[B]{
		parser: func(c *Cursor) (B, error) {
			val, err := pw.parser(c)
			if err != nil {
				var zero B
				return zero, err
			}
			return fn(val), nil
		},
	}
}

// Many parse une liste répétée d'éléments
func Many[T any](pw ParserWrapper[T]) ParserWrapper[[]T] {
	return ParserWrapper[[]T]{
		parser: func(c *Cursor) ([]T, error) {
			var result []T
			for {
				pos := c.Peek()
				val, err := pw.parser(c)
				if err != nil {
					c.Restore(pos)
					break
				}
				result = append(result, val)
			}
			return result, nil
		},
	}
}

// Opt rend un parser optionnel
func Opt[T any](pw ParserWrapper[T]) ParserWrapper[*T] {
	return ParserWrapper[*T]{
		parser: func(c *Cursor) (*T, error) {
			pos := c.Peek()
			val, err := pw.parser(c)
			if err != nil {
				c.Restore(pos)
				return nil, nil
			}
			return &val, nil
		},
	}
}

//
// Consume : vérifie le type de token
//

func Consume(kind syntax.Kind) ParserWrapper[syntax.Token] {
	p := func(c *Cursor) (syntax.Token, error) {
		tok, ok := c.Current()
		if !ok {
			return syntax.Token{}, ParseError{"unexpected end of input", 0, 0}
		}
		if tok.Kind() != kind {
			return syntax.Token{}, errorFromToken(tok, "unexpected token")
		}
		c.Advance()
		return tok, nil
	}
	return ParserWrapper[syntax.Token]{parser: p}
}

// ===========

//
// Tag : vérifie le texte exact du token
//

func Tag(text string) ParserWrapper[syntax.Token] {
	p := func(c *Cursor) (syntax.Token, error) {
		tok, ok := c.Current()
		if !ok {
			return syntax.Token{}, ParseError{"unexpected EOF", 0, 0}
		}
		value := tok.Value(c.source)
		if !bytes.Equal(value, []byte(text)) {
			return syntax.Token{}, errorFromToken(tok, "unexpected token value")
		}
		c.Advance()
		return tok, nil
	}
	return ParserWrapper[syntax.Token]{parser: p}
}

//
// Delimited : parse open + inner + close
//

func Delimited[A any, B any, C any](
	open ParserWrapper[A],
	inner ParserWrapper[B],
	close ParserWrapper[C],
) ParserWrapper[B] {
	return ParserWrapper[B]{
		parser: func(c *Cursor) (B, error) {
			if _, err := open.Parse(c); err != nil {
				var zero B
				return zero, err
			}

			val, err := inner.Parse(c)
			if err != nil {
				var zero B
				return zero, err
			}

			if _, err := close.Parse(c); err != nil {
				var zero B
				return zero, err
			}

			return val, nil
		},
	}
}

//
// Terminated : parse p + term
//

func Terminated[A any, B any](
	p ParserWrapper[A],
	term ParserWrapper[B],
) ParserWrapper[A] {
	return ParserWrapper[A]{
		parser: func(c *Cursor) (A, error) {
			val, err := p.Parse(c)
			if err != nil {
				var zero A
				return zero, err
			}
			if _, err := term.Parse(c); err != nil {
				var zero A
				return zero, err
			}
			return val, nil
		},
	}
}

//
// SeparatedList : liste d'éléments séparés par un séparateur
//

func SeparatedList[T any, S any](
	item ParserWrapper[T],
	sep ParserWrapper[S],
) ParserWrapper[[]T] {
	return ParserWrapper[[]T]{
		parser: func(c *Cursor) ([]T, error) {
			var result []T
			val, err := item.Parse(c)
			if err != nil {
				return nil, err
			}
			result = append(result, val)

			for {
				pos := c.Peek()
				if _, err := sep.Parse(c); err != nil {
					c.Restore(pos)
					break
				}

				val, err := item.Parse(c)
				if err != nil {
					c.Restore(pos)
					break
				}
				result = append(result, val)
			}

			return result, nil
		},
	}
}

//
// Choice : essaie plusieurs parsers jusqu'au succès
//

func Choice[T any](parsers ...ParserWrapper[T]) ParserWrapper[T] {
	return ParserWrapper[T]{
		parser: func(c *Cursor) (T, error) {
			var lastErr error
			for _, p := range parsers {
				pos := c.Peek()
				val, err := p.Parse(c)
				if err == nil {
					return val, nil
				}
				lastErr = err
				c.Restore(pos)
			}
			var zero T
			return zero, lastErr
		},
	}
}

//
// Peek : regarde un parser sans avancer le curseur
//

func Peek[T any](pw ParserWrapper[T]) ParserWrapper[T] {
	return ParserWrapper[T]{
		parser: func(c *Cursor) (T, error) {
			pos := c.Peek()
			val, err := pw.Parse(c)
			c.Restore(pos)
			return val, err
		},
	}
}
