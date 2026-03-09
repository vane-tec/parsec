package main

import (
	"os"
	"testing"
	"time"

	"github.com/vanetec/loom-syntax/syntax"
)

func BenchmarkIdlTokenize(b *testing.B) {
	data, err := os.ReadFile("/home/admi/Bureau/projects/loom-spec/loom/loom-syntax/test/prelude_test.loom")

	if err != nil {
		b.Fatal(err)
	}

	// data := []byte(`
	// // COmment
	// structure Test {}
	// `)

	lexer := syntax.NewLexer(data)
	b.ResetTimer()
	startTime := time.Now()

	for i := 0; i < b.N; i++ {
		lexer.DoTokenize()
	}

	duration := time.Since(startTime)
	b.ReportMetric(float64(duration.Milliseconds()), "ms_total")
}
