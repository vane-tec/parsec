package syntax

type Token struct {
	kind   Kind
	start  int
	end    int
	line   int
	column int
}

func NewToken(kind Kind, start int, end int, line int, column int, length int) *Token {
	return &Token{
		kind:   kind,
		start:  start,
		end:    end,
		line:   line,
		column: column,
	}
}

func (this *Token) Kind() Kind {
	return this.kind
}

func (this *Token) Start() int {
	return this.start
}

func (this *Token) End() int {
	return this.end
}

func (this *Token) Line() int {
	return this.line
}

func (this *Token) Column() int {
	return this.column
}

// Retourne la portion de texte correspondant au token.
// Aucune copie n'est faite : On slice directement le buffer.
func (this *Token) Value(src []byte) []byte {
	return src[this.start:this.end]
}
