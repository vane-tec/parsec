package syntax

type Lexer struct {
	input  []byte
	length int
	tokens []Token
	pos    int
	line   int
	column int
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{
		input:  input,
		length: len(input),
		pos:    0,
		line:   1,
		column: 1,
	}
}

func (this *Lexer) Input() []byte {
	return this.input
}
func (this *Lexer) addToken(k Kind, start int) {
	token := Token{
		kind:   k,
		start:  start,
		end:    this.pos,
		line:   this.line,
		column: this.column,
	}
	this.tokens = append(this.tokens, token)
}

func (this *Lexer) peek() byte {
	if this.pos >= this.length {
		return 0
	}
	return this.input[this.pos]
}

func (this *Lexer) peekNext() byte {
	if this.pos+1 >= this.length {
		return 0
	}

	return this.input[this.pos+1]
}

func (this *Lexer) advance() byte {
	c := this.peek()

	switch c {
	case '\n':
		this.line++
		this.column = 0
	case '\r':
		// ignore
	default:
		this.column++
	}

	if this.pos >= this.length {
		return 0
	}

	this.pos++

	return c
}

func isIdentifier(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isIdentPart(c byte) bool {
	return isIdentifier(c) || isDigit(c)
}

func isDigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

func (this *Lexer) DoTokenize() {

	for this.pos < this.length {
		start := this.pos
		ch := this.peek()

		this.advance()
		var kind Kind

		switch ch {
		case '{':
			kind = LBRACE
		case '}':
			kind = RBRACE
		case '[':
			kind = LBRACKET
		case ']':
			kind = RBRACKET
		case '(':
			kind = LPAREN
		case ')':
			kind = RPAREN
		case '=':
			kind = EQUAL
		case ':':
			kind = COLON
		case ',':
			kind = COMMA
		case '$':
			kind = DOLLAR
		case '#':
			kind = HASH
		case '@':
			kind = AT
		case '.':
			kind = DOT
		case '"', '\'', '`':
			kind = QUOTE
		case ' ', '\t', '\r':
			kind = WS
		default:
			if isIdentifier(ch) {
				for this.pos < this.length {
					c := this.peek()
					if !isIdentPart(c) {
						break
					}
					this.advance()
				}

				kind = IDENTIFIER
			}
		}

		this.addToken(kind, start)
	}

}

func (this *Lexer) GetTokens() []Token {
	return this.tokens
}
