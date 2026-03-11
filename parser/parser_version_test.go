package parser

import (
	"fmt"
	"testing"

	"github.com/vanetec/loom-syntax/syntax"
)

type VersionDirective struct {
	Value string
}

func TestParseVersionDirective_StrictQuote(t *testing.T) {

	parseVersionDirective := func(c *Cursor) (VersionDirective, error) {
		Consume(syntax.DOLLAR).Parse(c)

		// 1️⃣ Keyword $version
		tok, err := Consume(syntax.IDENTIFIER).Parse(c)
		if err != nil {
			return VersionDirective{}, err
		}

		if string(tok.Value(c.source)) != "version" {
			return VersionDirective{}, fmt.Errorf("expected '$version', got %s", tok.Value(c.source))
		}

		// 2️⃣ Optionnel : espaces multiples
		Many(Tag(" ")).Parse(c)

		// 3️⃣ Séparateur :
		_, err = Consume(syntax.COLON).Parse(c)
		if err != nil {
			return VersionDirective{}, fmt.Errorf("expected ':' after $version")
		}

		// 4️⃣ Optionnel : espaces après :
		Many(Tag(" ")).Parse(c)

		// 5️⃣ Ouverture quote
		_, err = Consume(syntax.QUOTE).Parse(c)
		if err != nil {
			return VersionDirective{}, fmt.Errorf("expected opening quote")
		}

		// 6️⃣ Consomme tous les tokens jusqu'à quote fermante ou fin de ligne
		var valueTokens []syntax.Token
		for c.pos < c.length {
			tok := c.tokens[c.pos]

			if tok.Kind() == syntax.QUOTE {
				// quote fermante trouvée → stop
				break
			}

			if tok.Kind() == syntax.NEWLINE {
				// fin de ligne atteinte avant quote fermante → erreur
				return VersionDirective{}, fmt.Errorf("expected closing quote before end of line")
			}

			valueTokens = append(valueTokens, tok)
			c.pos++
		}

		// 7️⃣ Consomme quote fermante
		_, err = Consume(syntax.QUOTE).Parse(c)
		if err != nil {
			return VersionDirective{}, fmt.Errorf("expected closing quote")
		}

		value := string(TokenSliceToBytes(valueTokens, c.source))

		return VersionDirective{Value: value}, nil
	}

	source := []byte(`$version: "1.0"`)
	lexer := syntax.NewLexer(source)
	lexer.DoTokenize()

	cursor := NewCursor(lexer.Input(), lexer.GetTokens())

	node, err := parseVersionDirective(cursor)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if node.Value != "1.0" {
		t.Fatalf("expected 1.0, got %s", node.Value)
	}

	fmt.Printf("%+v\n", node)
}

// helper : transforme slice de Token → []byte
func TokenSliceToBytes(tokens []syntax.Token, source []byte) []byte {
	if len(tokens) == 0 {
		return nil
	}
	start := tokens[0].Start()
	end := tokens[len(tokens)-1].End()
	return source[start:end]
}
