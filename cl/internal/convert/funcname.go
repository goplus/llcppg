package convert

import (
	"strings"

	"github.com/goplus/llcppg/ast"
)

// GoFuncSpec parses and stores components of a Go function name
// Input examples from llcppg.symb.json's "go" field:
// 1. Simple function: "AddPatchToArray"
// 2. Method with pointer receiver: "(*Sqlite3Stmt).Sqlite3BindParameterIndex"
// 3. Method with value receiver: "CJSONBool.CJSONCreateBool"
type GoFuncSpec struct {
	GoSymbName string // original full name from input
	FnName     string // function name without receiver
	IsMethod   bool   // if the function can be a method
	RecvName   string // receiver name
	PtrRecv    bool   // if the receiver is a pointer
}

// - "AddPatchToArray" -> {goSymbolName: "AddPatchToArray", funcName: "AddPatchToArray"}
// - "(*Sqlite3Stmt).Sqlite3BindParameterIndex" -> {goSymbolName: "...", recvName: "Sqlite3Stmt", funcName: "Sqlite3BindParameterIndex", ptrRecv: true}
// - "CJSONBool.CJSONCreateBool" -> {goSymbolName: "...", recvName: "CJSONBool", funcName: "CJSONCreateBool", ptrRecv: false}
func NewGoFuncSpec(name string, field []*ast.Field) *GoFuncSpec {
	l := strings.Split(name, ".")
	if len(l) < 2 {
		return &GoFuncSpec{GoSymbName: name, FnName: name, IsMethod: false}
	}
	recvName := l[0]
	ptrRecv := false
	if strings.HasPrefix(recvName, "(*") {
		ptrRecv = true
		recvName = strings.TrimPrefix(recvName, "(*")
		recvName = strings.TrimSuffix(recvName, ")")
	}

	// if the function can't be a method, we use the function name as the Go symbol name
	// not use the receiver style to link
	fnName := l[1]
	goSymbName := name
	beMethod := canBeMethod(field)
	if !beMethod {
		goSymbName = fnName
	}

	return &GoFuncSpec{
		GoSymbName: goSymbName,
		RecvName:   recvName,
		FnName:     fnName,
		PtrRecv:    ptrRecv,
		IsMethod:   beMethod,
	}
}

func (g *GoFuncSpec) IsIgnore() bool {
	return g.GoSymbName == "-"
}

func canBeMethod(fieldList []*ast.Field) bool {
	if len(fieldList) == 0 {
		return false
	}
	lastField := fieldList[len(fieldList)-1]
	_, isVar := lastField.Type.(*ast.Variadic)
	return !isVar
}
