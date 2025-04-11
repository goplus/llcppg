package dbg

type dbgFlags = int

var flags dbgFlags

const (
	DbgSymbol        dbgFlags = 1 << iota
	DbgParseIsMethod          //print parse.go isMethod debug log info
	DbgEditSymMap             //print user edit sym map info
	DbgVisitTop               //print visitTop
	DbgCollectFuncInfo
	DbgNewSymbol
	DbgFileType
	DbgFlagAll = 0
)

func SetDebugAll() {
	flags = DbgFlagAll
}

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

func SetDebugVisitTop() {
	flags |= DbgVisitTop
}

func GetDebugVisitTop() bool {
	return flags&DbgVisitTop != 0
}

func SetDebugCollectFuncInfo() {
	flags |= DbgCollectFuncInfo
}

func GetDebugCollectFuncInfo() bool {
	return flags&DbgCollectFuncInfo != 0
}

func SetDebugNewSymbol() {
	flags |= DbgNewSymbol
}

func GetDebugNewSymbol() bool {
	return flags&DbgNewSymbol != 0
}

func SetDebugFileType() {
	flags |= DbgFileType
}

func GetDebugFileType() bool {
	return flags&DbgFileType != 0
}
