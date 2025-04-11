package dbg

type dbgFlags = int

var flags dbgFlags

const (
	DbgSymbolNotFound dbgFlags = 1 << iota
	DbgError                   // print when error ocur
	DbgLog                     // print log info
	DbgSetCurFile
	DbgNew
	DbgWrite
	DbgUnmarshalling
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

func SetDebugNew() {
	flags |= DbgNew
}

func GetDebugNew() bool {
	return flags&DbgNew != 0
}

func SetDebugWrite() {
	flags |= DbgWrite
}

func GetDebugWrite() bool {
	return flags&DbgWrite != 0
}

func SetDebugUnmarshalling() {
	flags |= DbgUnmarshalling
}

func GetDebugUnmarshalling() bool {
	return flags&DbgUnmarshalling != 0
}
