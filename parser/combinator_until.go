package parser

import (
	"fmt"

	"github.com/vanetec/loom-syntax/syntax"
)

// Until consomme les tokens jusqu'à ce que `stop` réussisse.
// stop : ParserWrapper[T] (ex : Consume(LBRACE) ou Tag("{"))
// Retourne un ParserWrapper[[]Token] comme tes autres combinators
func Until[T any](stop ParserWrapper[T]) ParserWrapper[[]syntax.Token] {
	return ParserWrapper[[]syntax.Token]{
		parser: func(c *Cursor) ([]syntax.Token, error) {
			var consumed []syntax.Token

			for c.pos < len(c.tokens) {
				// On teste le stop parser sur la position actuelle
				_, err := stop.Parse(c)
				if err == nil {
					// stop trouvé → on s'arrête ici
					return consumed, nil
				}

				// sinon on consomme le token courant
				consumed = append(consumed, c.tokens[c.pos])
				c.pos++
			}

			return consumed, fmt.Errorf("Until: stop parser not found")
		},
	}
}
