package parse

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llcppg/_xtool/llcppsymg/dbg"
	"github.com/goplus/llcppg/_xtool/llcppsymg/names"
	"github.com/goplus/llgo/c/clang"
)

type SymbolInfo struct {
	GoName    string
	ProtoName string
}

type SymbolProcessor struct {
	Files      []string
	Prefixes   []string
	SymbolMap  map[string]*SymbolInfo
	NameCounts map[string]int
	// for independent files,signal that the file has been processed
	// will clean in a translation unit process end
	processingFiles map[string]struct{}
	processedFiles  map[string]struct{}
}

func NewSymbolProcessor(Files []string, Prefixes []string) *SymbolProcessor {
	return &SymbolProcessor{
		Files:           Files,
		Prefixes:        Prefixes,
		SymbolMap:       make(map[string]*SymbolInfo),
		NameCounts:      make(map[string]int),
		processedFiles:  make(map[string]struct{}),
		processingFiles: make(map[string]struct{}),
	}
}

func (p *SymbolProcessor) isSelfFile(filename string) bool {
	for _, file := range p.Files {
		if file == filename {
			return true
		}
	}
	return false
}

func (p *SymbolProcessor) typeCursor(arg clang.Cursor) clang.Cursor {
	typ := arg.Type()
	if typ.Kind == clang.TypePointer {
		typ = typ.PointeeType()
	}
	return typ.TypeDeclaration()
}

func (p *SymbolProcessor) cursorFileName(cur clang.Cursor, isArg bool) (ret string) {
	if isArg {
		typCursor := p.typeCursor(cur)
		filename := ""
		if len(clang.GoString(typCursor.String())) > 0 {
			filename = clang.GoString(typCursor.Location().File().FileName())
		}
		return filename
	}
	return clang.GoString(cur.Location().File().FileName())
}

func (p *SymbolProcessor) inCurPkg(cur clang.Cursor, isArg bool) bool {
	if false {
		typ := p.typeCursor(cur)
		fmt.Println(
			clang.GoString(cur.DisplayName()),
			p.cursorFileName(cur, isArg),
			"type:", clang.GoString(typ.String()),
			p.isSelfFile(p.cursorFileName(cur, isArg)),
		)
	}
	return p.isSelfFile(p.cursorFileName(cur, isArg))
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

func (p *SymbolProcessor) printTypeInfo(typ clang.Type, isArg bool, prefix string) {
	if dbg.GetDebugParseIsMethod() {
		definitionType := clang.GoString(typ.TypeDeclaration().Definition().Type().String())
		canonicalType := clang.GoString(typ.CanonicalType().String())
		fmt.Println("**********", prefix, "**********")
		fmt.Println(
			"typ.String():", clang.GoString(typ.String()),
			"typ.NamedType().String():", clang.GoString(typ.NamedType().String()),
			"isTypePointer:", typ.Kind == clang.TypePointer,
			"isTypeElaborated:", typ.Kind == clang.TypeElaborated,
			"isTypeTypedef:", typ.Kind == clang.TypeTypedef,
			"isInCurPkg:", p.inCurPkg(typ.TypeDeclaration(), isArg),
			"definitionType:", definitionType,
			"canonicalType:", canonicalType,
			"pointLevel:", p.pointerLevel(typ),
		)
	}
}

func printResult(isInCurPkg, isPointer bool, goName, prefix string) {
	if dbg.GetDebugParseIsMethod() {
		fmt.Println("===========", prefix, "===========")
		fmt.Println("isInCurPkg:", isInCurPkg, "isPointer:", isPointer, "goName", goName)
	}
}

func (p *SymbolProcessor) pointerLevel(typ clang.Type) int {
	canonicalTypeGoString := clang.GoString(typ.CanonicalType().String())
	return strings.Count(canonicalTypeGoString, "*")
}

func (p *SymbolProcessor) isMethod(cur clang.Cursor, isArg bool) (bool, bool, string) {
	typ := cur.Type()
	if p.pointerLevel(typ) > 1 {
		return false, false, names.GoName(clang.GoString(typ.String()), p.Prefixes, false)
	}
	isInCurPkg := p.inCurPkg(cur, isArg)
	p.printTypeInfo(typ, isArg, "typ")
	if typ.Kind == clang.TypePointer {
		pointeeType := typ.PointeeType()
		p.printTypeInfo(pointeeType, isArg, "typ.PointeeType()")
		pointeeTypeNamedType := pointeeType.NamedType()
		namedTypeGoString := clang.GoString(pointeeTypeNamedType.String())
		p.printTypeInfo(pointeeTypeNamedType, isArg, "typ.PointeeType().NamedType()")
		if len(namedTypeGoString) > 0 {
			goName := names.GoName(namedTypeGoString, p.Prefixes, isInCurPkg)
			printResult(isInCurPkg, true, goName, "typ.pointeeType().NamedType()")
			return isInCurPkg, true, goName
		}
		return p.isMethod(pointeeType.TypeDeclaration(), isArg)
	} else if typ.Kind == clang.TypeElaborated ||
		typ.Kind == clang.TypeTypedef {
		canonicalType := typ.CanonicalType()
		p.printTypeInfo(canonicalType, isArg, "typ.CanonicalType()")
		if canonicalType.Kind == clang.TypePointer {
			return p.isMethod(canonicalType.TypeDeclaration(), isArg)
		}
	}
	namedType := typ.NamedType()
	namedTypeGoString := clang.GoString(namedType.String())
	if len(namedTypeGoString) > 0 {
		goName := names.GoName(namedTypeGoString, p.Prefixes, isInCurPkg)
		printResult(isInCurPkg, false, goName, "typ.NamedType()")
		return isInCurPkg, false, goName
	}
	typeGoString := clang.GoString(typ.String())
	goName := names.GoName(typeGoString, p.Prefixes, isInCurPkg)
	printResult(isInCurPkg, false, goName, "typ")
	return isInCurPkg, false, goName
}

func (p *SymbolProcessor) genGoName(cursor clang.Cursor) string {
	originName := clang.GoString(cursor.String())
	isDestructor := cursor.Kind == clang.CursorDestructor
	var convertedName string
	if isDestructor {
		convertedName = names.GoName(originName[1:], p.Prefixes, p.inCurPkg(cursor, false))
	} else {
		convertedName = names.GoName(originName, p.Prefixes, p.inCurPkg(cursor, false))
	}

	if parent := cursor.SemanticParent(); parent.Kind == clang.CursorClassDecl {
		class := names.GoName(clang.GoString(parent.String()), p.Prefixes, p.inCurPkg(cursor, false))
		return p.AddSuffix(p.GenMethodName(class, convertedName, isDestructor, true))
	} else if cursor.Kind == clang.CursorFunctionDecl {
		numArgs := cursor.NumArguments()
		if numArgs > 0 {
			if ok, isPtr, typeName := p.isMethod(cursor.Argument(0), true); ok {
				return p.AddSuffix(p.GenMethodName(typeName, convertedName, isDestructor, isPtr))
			}
		}
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
	p.NameCounts[name]++
	if count := p.NameCounts[name]; count > 1 {
		return name + "__" + strconv.Itoa(count-1)
	}
	return name
}

func (p *SymbolProcessor) collectFuncInfo(cursor clang.Cursor) {
	// On Linux, C++ symbols typically have one leading underscore
	// On macOS, C++ symbols may have two leading underscores
	// For consistency, we remove the first leading underscore on macOS
	if dbg.GetDebugSymbol() {
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
	if _, exists := p.SymbolMap[symbolName]; exists {
		return
	}
	p.SymbolMap[symbolName] = &SymbolInfo{
		GoName:    p.genGoName(cursor),
		ProtoName: p.genProtoName(cursor),
	}
}

func (p *SymbolProcessor) visitTop(cursor, parent clang.Cursor) clang.ChildVisitResult {
	filename := clang.GoString(cursor.Location().File().FileName())
	if _, ok := p.processedFiles[filename]; ok {
		if dbg.GetDebugSymbol() {
			fmt.Printf("visitTop: %s has been processed: \n", filename)
		}
		return clang.ChildVisit_Continue
	}
	if filename == "" {
		return clang.ChildVisit_Continue
	}
	p.processingFiles[filename] = struct{}{}
	if dbg.GetDebugSymbol() && filename != "" {
		fmt.Printf("visitTop: %s\n", filename)
	}
	switch cursor.Kind {
	case clang.CursorNamespace, clang.CursorClassDecl:
		clangutils.VisitChildren(cursor, p.visitTop)
	case clang.CursorCXXMethod, clang.CursorFunctionDecl, clang.CursorConstructor, clang.CursorDestructor:
		isPublicMethod := (cursor.CXXAccessSpecifier() == clang.CXXPublic) && cursor.Kind == clang.CursorCXXMethod || cursor.Kind == clang.CursorConstructor || cursor.Kind == clang.CursorDestructor
		if p.isSelfFile(filename) && (cursor.Kind == clang.CursorFunctionDecl || isPublicMethod) {
			p.collectFuncInfo(cursor)
		}
	}
	return clang.ChildVisit_Continue
}

func (p *SymbolProcessor) collect(cfg *clangutils.Config) error {
	filename := cfg.File
	if cfg.Temp {
		filename = clangutils.TEMP_FILE
	}
	if _, ok := p.processedFiles[filename]; ok {
		if dbg.GetDebugSymbol() {
			fmt.Printf("%s has been processed: \n", filename)
		}
		return nil
	}
	if dbg.GetDebugSymbol() {
		fmt.Printf("create translation unit: \nfile:%s\nIsCpp:%v\nTemp:%v\nArgs:%v\n", filename, cfg.IsCpp, cfg.Temp, cfg.Args)
	}
	_, unit, err := clangutils.CreateTranslationUnit(cfg)
	if err != nil {
		return errors.New("Unable to parse translation unit for file " + filename)
	}
	cursor := unit.Cursor()
	if dbg.GetDebugSymbol() {
		fmt.Printf("%s start collect \n", filename)
	}
	clangutils.VisitChildren(cursor, p.visitTop)
	for filename := range p.processingFiles {
		p.processedFiles[filename] = struct{}{}
	}
	p.processingFiles = make(map[string]struct{})
	unit.Dispose()
	return nil
}

func ParseHeaderFile(files []string, prefixes []string, cflags []string, isCpp bool, isTemp bool) (map[string]*SymbolInfo, error) {
	index := clang.CreateIndex(0, 0)
	if isTemp {
		files = append(files, clangutils.TEMP_FILE)
	}
	processer := NewSymbolProcessor(files, prefixes)
	for _, file := range files {
		processer.collect(&clangutils.Config{
			File:  file,
			Temp:  isTemp,
			IsCpp: isCpp,
			Index: index,
			Args:  cflags,
		})
	}
	index.Dispose()
	return processer.SymbolMap, nil
}
