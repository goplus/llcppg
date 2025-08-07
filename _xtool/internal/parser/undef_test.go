package parser_test

import (
	"testing"

	"github.com/goplus/llcppg/_xtool/internal/parser"
)

func TestUndefType(t *testing.T) {
	// Test that undefined types are handled correctly
	ast, err := parser.Do(&parser.ConverterConfig{
		File:  "../testdata/undef_type/temp.h",
		IsCpp: false,
		Args:  []string{"-fparse-all-comments"},
	})
	if err != nil {
		t.Fatal("Do failed:", err)
	}
	
	// We expect no function declarations since the undefined type should be detected
	if len(ast.Decls) > 0 {
		t.Fatalf("Expected no declarations for undefined type, got %d declarations", len(ast.Decls))
	}
}