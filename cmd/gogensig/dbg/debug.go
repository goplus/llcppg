package dbg

type dbgFlags = int

var flags dbgFlags

const (
	DbgSymbolNotFound dbgFlags = 1 << iota
	DbgError                   // print when error ocur
	DbgLog                     // print log info
	DbgSetCurFile
	DbgNewFuncDecl
	DbgFlagAll = 0 | DbgError | DbgLog
)

func SetDebugSymbolNotFound() {
	flags |= DbgSymbolNotFound
}

func GetDebugSymbolNotFound() bool {
	return flags&DbgSymbolNotFound != 0
}

func SetDebugError() {
	flags |= DbgError
}

func GetDebugError() bool {
	return flags&DbgError != 0
}

func SetDebugLog() {
	flags |= DbgLog
}

func GetDebugLog() bool {
	return flags&DbgLog != 0
}

func SetDebugAll() {
	flags = DbgFlagAll
}

func SetDebugSetCurFile() {
	flags |= DbgSetCurFile
}

func GetDebugSetCurFile() bool {
	return flags&DbgSetCurFile != 0
}

func SetDebugNewFuncDecl() {
	flags |= DbgNewFuncDecl
}

func GetDebugNewFuncDecl() bool {
	return flags&DbgNewFuncDecl != 0
}
