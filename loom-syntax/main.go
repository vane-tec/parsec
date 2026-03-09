package main

import (
	"fmt"
	"os"

	"github.com/vanetec/loom-syntax/syntax"
)

func main() {
	TestA()
	// parser.TestC()
}

func TestA() {
	// { 15
	// data := []byte(`
	// // COmment
	// structure Test {}
	// `)

	data, err := os.ReadFile("/home/admi/Bureau/projects/loom-spec/loom/loom-syntax/test/prelude_test.loom")

	if err != nil {
		panic(err)
	}

	lexer := syntax.NewLexer(data)
	lexer.DoTokenize()

	// dev
	tokens := lexer.GetTokens()
	for _, t := range tokens {
		e := t.End()
		s := t.Start()
		l := t.Line()
		c := t.Column()
		k := t.Kind()
		fmt.Printf("[l : %v, c : %v, s : %v, e : %v , k: %v] %s \n", l, c, s, e, k, data[s:e])
	}
	// fmt.Printf("%+v\n", tokens)
	// println(string(data[15:16]))
	// println(len(tokens))
}
