package parser

import (
	"fmt"

	"github.com/vanetec/loom-syntax/syntax"
)

// UntilAny consomme les tokens jusqu'à ce que l'un des stop parsers réussisse.
// stops : liste de ParserWrapper[T] représentant les conditions d'arrêt
func UntilAny[T any](stops ...ParserWrapper[T]) ParserWrapper[[]syntax.Token] {
	return ParserWrapper[[]syntax.Token]{
		parser: func(c *Cursor) ([]syntax.Token, error) {
			var consumed []syntax.Token

			for c.pos < len(c.tokens) {
				// on teste tous les stop parsers
				found := false
				for _, stop := range stops {
					_, err := stop.Parse(c)
					if err == nil {
						// un stop trouvé → on s'arrête
						found = true
						break
					}
				}

				if found {
					return consumed, nil
				}

				// sinon consomme le token courant
				consumed = append(consumed, c.tokens[c.pos])
				c.pos++
			}

			return consumed, fmt.Errorf("UntilAny: aucun stop parser trouvé")
		},
	}
}
