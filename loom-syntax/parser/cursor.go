package parser

import "github.com/vanetec/loom-syntax/syntax"

// Représente la position actuelle dans la liste de tokens.
//
// Tous les combinators manipulent ce curseur.
// Le design évite toute copie de slice ou allocation inutile.
type Cursor struct {
	source []byte
	tokens []syntax.Token
	length int
	pos    int
}

func NewCursor(s []byte, t []syntax.Token) *Cursor {
	return &Cursor{source: s, tokens: t, length: len(s), pos: 0}
}

// Retourne le token actuel
func (this *Cursor) Current() (syntax.Token, bool) {
	if this.pos >= this.length {
		return syntax.Token{}, false
	}
	return this.tokens[this.pos], true
}

// Avance le curseur d'un token.
func (this *Cursor) Advance() {
	this.pos++
}

// Permet de sauvegarder la position actuelle.
// Utile pour les parsers backtraacking.
func (this *Cursor) Peek() int {
	return this.pos
}

// Remet la position précédente
func (this *Cursor) Restore(pos int) {
	this.pos = pos
}
