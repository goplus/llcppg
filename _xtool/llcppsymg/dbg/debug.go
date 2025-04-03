package dbg

type dbgFlags = int

var flags dbgFlags

const (
	DbgSymbol        dbgFlags = 1 << iota
	DbgParseIsMethod          //print parse.go isMethod debug log info
	DbgEditSymMap             //print user edit sym map info
	DbgFlagAll       = DbgSymbol | DbgParseIsMethod
)

func SetDebugSymbol() {
	flags |= DbgSymbol
}

func GetDebugSymbol() bool {
	return flags&DbgSymbol != 0
}

func SetDebugParseIsMethod() {
	flags |= DbgParseIsMethod
}

func GetDebugParseIsMethod() bool {
	return flags&DbgParseIsMethod != 0
}

func SetDebugEditSymMap() {
	flags |= DbgEditSymMap
}

func GetDebugEditSymMap() bool {
	return flags&DbgEditSymMap != 0
}
