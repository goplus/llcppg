package cl

import "testing"

func TestCstyleToGoForceCamelCase(t *testing.T) {
	tests := []struct {
		name   string
		force  bool
		rename map[string]string
		input  string
		want   string
	}{
		{name: "default", input: "Cursor_StructDecl", want: "Cursor_StructDecl"},
		{name: "enabled", force: true, input: "Cursor_StructDecl", want: "CursorStructDecl"},
		{name: "uppercase suffix", force: true, input: "Cursor_ID", want: "Cursor_ID"},
		{name: "macro", force: true, input: "FLAG_READ_ONLY", want: "FLAG_READ_ONLY"},
		{name: "snake case", force: true, input: "set_name", want: "SetName"},
		{name: "part rename", force: true, rename: map[string]string{"StructDecl": "Decl"}, input: "Cursor_StructDecl", want: "CursorDecl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pkgCtx{forceCamelCase: tt.force, rename: tt.rename}
			if got := p.cstyleToGo(tt.input, true); got != tt.want {
				t.Errorf("cstyleToGo(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFuncNameRenamePrecedesForceCamelCase(t *testing.T) {
	p := pkgCtx{
		forceCamelCase: true,
		rename:         map[string]string{"Cursor_StructDecl": "SpecialName"},
	}
	if got := p.funcName("Cursor_StructDecl", -1, "", "", true, false); got != "SpecialName" {
		t.Errorf("funcName() = %q, want %q", got, "SpecialName")
	}
}
