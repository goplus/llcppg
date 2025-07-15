package symg

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/goplus/lib/c/clang"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	"github.com/goplus/llcppg/internal/name"

	llcppg "github.com/goplus/llcppg/config"
)

type HeaderSymbols map[string]*SymbolInfo

func (h HeaderSymbols) ToSymbolTable() (table []*llcppg.SymbolInfo) {
	for name, info := range h {
		table = append(table, &llcppg.SymbolInfo{
			Go:     info.GoName,
			CPP:    info.ProtoName,
			Mangle: name,
		})
	}
	return
}

type SymbolInfo struct {
	GoName    string
	ProtoName string
}

type collect struct {
	symName    string             // symbol name
	getSymInfo func() *SymbolInfo // get symbol info
}

type SymbolProcessor struct {
	curPkgFiles map[string]struct{}
	prefixes    []string
	symbolMap   map[string]*SymbolInfo
	nameCounts  map[string]int
	// custom symbol map like:
	// "sqlite3_finalize":".Close" -> method
	// "sqlite3_open":"Open" -> function
	customSymMap map[string]string
	// register queue
	collectQueue []*collect
}

func NewSymbolProcessor(curPkgFiles []string, prefixes []string, symMap map[string]string) *SymbolProcessor {
	p := &SymbolProcessor{
		prefixes:     prefixes,
		customSymMap: symMap,
		symbolMap:    make(map[string]*SymbolInfo),
		nameCounts:   make(map[string]int),
		curPkgFiles:  make(map[string]struct{}),
	}
	for _, file := range curPkgFiles {
		p.curPkgFiles[file] = struct{}{}
	}
	return p
}

func (p *SymbolProcessor) isSelfFile(filename string) bool {
	_, ok := p.curPkgFiles[filename]
	if !ok && dbgSymbol {
		fmt.Println("not in file: ", filename)
	}
	return ok
}

// Is the cursor in the current package
func (p *SymbolProcessor) inCurPkg(cursor clang.Cursor) bool {
	return p.isSelfFile(cursorFileName(cursor))
}

func (p *SymbolProcessor) GenMethodName(class, name string, isDestructor bool, isPointer bool) string {
	prefix := class + "."
	if isPointer {
		prefix = "(*" + class + ")."
	}
	if isDestructor {
		return prefix + "Dispose"
	}
	if class == name {
		return prefix + "Init"
	}
	return prefix + name
}

func (p *SymbolProcessor) pointerLevel(typ clang.Type) int {
	canonicalTypeGoString := clang.GoString(typ.CanonicalType().String())
	return strings.Count(canonicalTypeGoString, "*")
}

// check cursor can be a receiver
func (p *SymbolProcessor) beRecv(cur clang.Cursor) (ok bool, isPtr bool, typeName string) {
	typ := cur.Type()
	// Check if the type has more than one level of pointers (e.g., **int, ***char)
	// Multi-level pointers are not considered for method generation as they are
	// typically used for complex data structures rather than object instances
	if p.pointerLevel(typ) > 1 {
		return false, false, ""
	}

	// check the arg's type's location
	isInCurPkg := p.inCurPkg(underCursor(cur))

	if typ.Kind == clang.TypePointer {
		underTypeName := clang.GoString(typ.PointeeType().NamedType().String())
		return isInCurPkg, true, name.GoName(underTypeName, p.prefixes, isInCurPkg)
	}

	// Check if the type is an elaborated type (e.g., struct/class with full qualification)
	// or a typedef (type alias). These types may wrap underlying pointer types
	// that should be considered for method generation
	if typ.Kind == clang.TypeElaborated || typ.Kind == clang.TypeTypedef {
		underTyp := typ.CanonicalType()
		if underTyp.Kind == clang.TypePointer {
			return p.beRecv(underTyp.TypeDeclaration())
		}
	}
	namedType := clang.GoString(typ.NamedType().String())
	goName := name.GoName(namedType, p.prefixes, isInCurPkg)
	return isInCurPkg, false, goName
}

// sqlite3_finalize -> .Close -> method
// sqlite3_open -> Open -> function
func (p *SymbolProcessor) customGoName(mangled string) (goName string, isMethod bool, ok bool) {
	if customName, ok := p.customSymMap[mangled]; ok {
		name, found := strings.CutPrefix(customName, ".")
		return name, found, true
	}
	return "", false, false
}

func (p *SymbolProcessor) genGoName(cursor clang.Cursor, symbolName string) string {
	originName := clang.GoString(cursor.String())
	isDestructor := cursor.Kind == clang.CursorDestructor
	var convertedName string
	if isDestructor {
		convertedName = name.GoName(originName[1:], p.prefixes, p.inCurPkg(cursor))
	} else {
		convertedName = name.GoName(originName, p.prefixes, p.inCurPkg(cursor))
	}

	customGoName, toMethod, isCustom := p.customGoName(symbolName)

	// 1. for class method,gen method name
	if parent := cursor.SemanticParent(); parent.Kind == clang.CursorClassDecl {
		class := name.GoName(clang.GoString(parent.String()), p.prefixes, p.inCurPkg(cursor))
		// concat method name
		if isCustom {
			convertedName = customGoName
		}
		return p.AddSuffix(p.GenMethodName(class, convertedName, isDestructor, true))
	}

	// 2. check if can gen method name
	numArgs := cursor.NumArguments()
	// 3. Don't attempt to convert a variadic function to a method
	isValist := cursor.Type().IsFunctionTypeVariadic() > 0
	// also config to gen method name,if can't gen method,use the origin function type
	if numArgs > 0 && !isValist {
		// also can gen method name,but not want to be method,output func not method
		if isCustom && !toMethod {
			return p.AddSuffix(customGoName)
		}
		if ok, isPtr, typeName := p.beRecv(cursor.Argument(0)); ok {
			if isCustom {
				convertedName = customGoName
			}
			return p.AddSuffix(p.GenMethodName(typeName, convertedName, isDestructor, isPtr))
		}
	}

	// 3. normal function name
	if isCustom {
		return p.AddSuffix(customGoName)
	}
	return p.AddSuffix(convertedName)
}

func (p *SymbolProcessor) genProtoName(cursor clang.Cursor) string {
	scopingParts := clangutils.BuildScopingParts(cursor.SemanticParent())

	var builder strings.Builder
	for _, part := range scopingParts {
		builder.WriteString(part)
		builder.WriteString("::")
	}

	builder.WriteString(clang.GoString(cursor.DisplayName()))
	return builder.String()
}

func (p *SymbolProcessor) AddSuffix(name string) string {
	p.nameCounts[name]++
	if count := p.nameCounts[name]; count > 1 {
		return name + "__" + strconv.Itoa(count-1)
	}
	return name
}

func (p *SymbolProcessor) collectFuncInfo(cursor clang.Cursor) {
	// On Linux, C++ symbols typically have one leading underscore
	// On macOS, C++ symbols may have two leading underscores
	// For consistency, we remove the first leading underscore on macOS
	if dbgSymbol {
		fmt.Printf("collectFuncInfo: %s %s\n", clang.GoString(cursor.Mangling()), clang.GoString(cursor.String()))
	}
	symbolName := clang.GoString(cursor.Mangling())
	if runtime.GOOS == "darwin" {
		symbolName = strings.TrimPrefix(symbolName, "_")
	}

	// In C, multiple declarations of the same function are allowed.
	// Functions with identical signatures will have the same mangled name.
	// We treat them as the same function rather than overloads, so we only
	// process the first occurrence and skip subsequent declarations.
	if _, exists := p.symbolMap[symbolName]; exists {
		return
	}
	p.symbolMap[symbolName] = &SymbolInfo{}
	p.collectQueue = append(p.collectQueue, &collect{
		symName: symbolName,
		getSymInfo: func() *SymbolInfo {
			return &SymbolInfo{
				GoName:    p.genGoName(cursor, symbolName),
				ProtoName: p.genProtoName(cursor),
			}
		},
	})
}

func (p *SymbolProcessor) visitTop(cursor, parent clang.Cursor) clang.ChildVisitResult {
	filename := cursorFileName(cursor)
	switch cursor.Kind {
	case clang.CursorNamespace, clang.CursorClassDecl:
		clangutils.VisitChildren(cursor, p.visitTop)
	case clang.CursorCXXMethod, clang.CursorFunctionDecl, clang.CursorConstructor, clang.CursorDestructor:
		isPublicFunc := cursor.Kind == clang.CursorFunctionDecl &&
			cursor.StorageClass() != clang.SCStatic

		isPublicMethod := cursor.CXXAccessSpecifier() == clang.CXXPublic &&
			cursor.Kind == clang.CursorCXXMethod ||
			cursor.Kind == clang.CursorConstructor ||
			cursor.Kind == clang.CursorDestructor

		if p.isSelfFile(filename) && (isPublicFunc || isPublicMethod) {
			p.collectFuncInfo(cursor)
		}
	}
	return clang.ChildVisit_Continue
}

// processCollect processes the symbol collection queue and prioritizes custom go names.
// Custom symbols (defined in llcppg.cfg/symMap) are processed before regular symbols
// to ensure user-defined mappings take precedence.
func (p *SymbolProcessor) processCollect() {
	sort.SliceStable(p.collectQueue, func(i, j int) bool {
		_, customI := p.customSymMap[p.collectQueue[i].symName]
		_, customJ := p.customSymMap[p.collectQueue[j].symName]
		return customI && !customJ
	})
	for _, collect := range p.collectQueue {
		p.symbolMap[collect.symName] = collect.getSymInfo()
	}
}

// Get the underlying cursor of the cursor
// if cur is a pointer,return the underlying cursor
func underCursor(arg clang.Cursor) clang.Cursor {
	typ := arg.Type()
	if typ.Kind == clang.TypePointer {
		typ = typ.PointeeType()
	}
	return typ.TypeDeclaration()
}

// https://releases.llvm.org/19.1.0/tools/clang/docs/ReleaseNotes.html#libclang
// cursor.Location() in llvm@19 cannot get the fileinfo for a macro expansion,so we dirrect use PresumedLocation
func cursorFileName(cursor clang.Cursor) string {
	loc := cursor.Location()
	filePath, _, _ := clangutils.GetPresumedLocation(loc)
	return filePath
}

func ParseHeaderFile(combileFile string, curPkgFiles []string, prefixes []string, cflags []string, symMap map[string]string, isCpp bool) (HeaderSymbols, error) {
	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
		File:    combileFile,
		IsCpp:   isCpp,
		Index:   clang.CreateIndex(0, 0),
		Args:    cflags,
		Options: clang.DetailedPreprocessingRecord,
	})
	if err != nil {
		return nil, err
	}
	defer unit.Dispose()
	defer index.Dispose()
	cursor := unit.Cursor()
	processer := NewSymbolProcessor(curPkgFiles, prefixes, symMap)
	clangutils.VisitChildren(cursor, processer.visitTop)
	processer.processCollect()
	return HeaderSymbols(processer.symbolMap), nil
}
