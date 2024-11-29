package convert

import "strings"

type GoFuncName struct {
	goSymbolName string
	recvName     string
	funcName     string
}

func NewGoFuncName(name string) *GoFuncName {
	l := strings.Split(name, ".")
	if len(l) < 2 {
		return &GoFuncName{goSymbolName: name, funcName: name}
	}
	return &GoFuncName{goSymbolName: name, recvName: l[0], funcName: l[1]}
}

func (p *GoFuncName) HasReceiver() bool {
	return len(p.recvName) > 0
}

func (p *GoFuncName) OriginGoSymbolName() string {
	return p.goSymbolName
}
