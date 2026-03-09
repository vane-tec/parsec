package parser

import (
	"testing"

	"github.com/vanetec/loom-syntax/syntax"
)

// type Structure struct {
// 	Name string
// }

func parseStructure(c *Cursor) (Structure, error) {

	// structure keyword
	tok, err := Consume(syntax.IDENTIFIER).Parse(c)
	if err != nil {
		return Structure{}, err
	}

	if string(tok.Value(c.source)) != "structure" {
		return Structure{}, err
	}

	// space
	Many(Tag(" ")).Parse(c)
	// Tag(" ").Parse(c)

	// name
	nameTok, err := Consume(syntax.IDENTIFIER).Parse(c)
	if err != nil {
		return Structure{}, err
	}

	name := string(nameTok.Value(c.source))

	// space
	Tag(" ").parser(c)

	// {
	_, err = Consume(syntax.LBRACE).Parse(c)
	if err != nil {
		return Structure{}, err
	}

	// }
	_, err = Consume(syntax.RBRACE).Parse(c)
	if err != nil {
		return Structure{}, err
	}

	return Structure{Name: name}, nil
}

func BenchmarkParseStructure(b *testing.B) {

	source := []byte(`structure   Test {}`)

	for i := 0; i < b.N; i++ {

		lexer := syntax.NewLexer(source)
		lexer.DoTokenize()

		cursor := NewCursor(lexer.Input(), lexer.GetTokens())

		_, err := parseStructure(cursor)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseStructure2(b *testing.B) {

	source := []byte(`structure Test {}`)

	lexer := syntax.NewLexer(source)
	lexer.DoTokenize()

	tokens := lexer.GetTokens()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		cursor := NewCursor(source, tokens)

		_, err := parseStructure(cursor)
		if err != nil {
			b.Fatal(err)
		}
	}
}
