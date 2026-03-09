package parser

/*
*
* Chaque combinator :
* lit les tokens
* avance pos
* retourne un résultat
 */
type Parser[T any] func(*Cursor) (T, error)

type ParserWrapper[T any] struct {
	parser Parser[T]
}
