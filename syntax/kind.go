package syntax

type Kind uint32

const (
	EOF Kind = iota

	// {
	LBRACE

	// }
	RBRACE

	// [
	LBRACKET

	// ]
	RBRACKET

	// (
	LPAREN

	// )
	RPAREN

	// =
	EQUAL

	// :
	COLON

	// ,
	COMMA

	COMMENT

	NEWLINE
	QUOTE
	AT
	DOLLAR
	DOT
	HASH
	STRING_SINGLE
	STRING_DOUBLE
	NUMBER
	IDENTIFIER

	// Whitespace
	WS

	// Tout autres caractère spécial
	SYMBOL
)
