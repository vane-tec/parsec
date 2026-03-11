package parser

import (
	"fmt"
	"testing"

	"github.com/vanetec/loom-syntax/syntax"
)

type Structure struct {
	Name string
}

func TestParseStructure(t *testing.T) {

	parseStructure := func(c *Cursor) (Structure, error) {

		// structure keyword
		tok, err := Consume(syntax.IDENTIFIER).Parse(c)
		if err != nil {
			return Structure{}, err
		}

		if string(tok.Value(c.source)) != "structure" {
			t.Fatalf("expected keyword 'structure'")
		}

		// structure name
		Tag(" ").parser(c)
		nameTok, err := Consume(syntax.IDENTIFIER).Parse(c)
		if err != nil {
			return Structure{}, err
		}

		name := string(nameTok.Value(c.source))

		// {
		Tag(" ").parser(c)
		_, err = Consume(syntax.LBRACE).Parse(c)
		if err != nil {
			return Structure{}, err
		}

		// }
		_, err = Consume(syntax.RBRACE).Parse(c)
		if err != nil {
			return Structure{}, err
		}

		return Structure{
			Name: name,
		}, nil
	}

	source := []byte(`structure Test {}`)
	lexer := syntax.NewLexer(source)
	lexer.DoTokenize()

	cursor := NewCursor(lexer.Input(), lexer.GetTokens())
	// tokens := cursor.tokens
	// fmt.Printf("\n%+v \n", tokens)

	node, err := parseStructure(cursor)

	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if node.Name != "Test" {
		t.Fatalf("expected Test, got %s", node.Name)
	}

	fmt.Printf("%+v\n", node)

}
