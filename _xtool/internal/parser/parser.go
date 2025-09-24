package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/goplus/lib/c"
	clangutils "github.com/goplus/llcppg/_xtool/internal/clang"
	clang "github.com/goplus/llcppg/_xtool/internal/libclang"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/token"
)

type dbgFlags = int

var debugParse bool

const (
	DbgParse   dbgFlags = 1 << iota
	DbgFlagAll          = DbgParse
)

func SetDebug(dbgFlags dbgFlags) {
	debugParse = (dbgFlags & DbgParse) != 0
}

type Converter struct {
	file   *ast.File
	index  *clang.Index
	unit   *clang.TranslationUnit
	indent int // for verbose debug
}

var tagMap = map[string]ast.Tag{
	"struct": ast.Struct,
	"union":  ast.Union,
	"enum":   ast.Enum,
	"class":  ast.Class,
}

type ConverterConfig struct {
	File  string
	Args  []string
	IsCpp bool
}

func Do(config *ConverterConfig) (*ast.File, error) {
	converter, err := NewConverter(config)
	if err != nil {
		return nil, err
	}
	return converter.Convert()
}

func NewConverter(config *ConverterConfig) (*Converter, error) {
	if debugParse {
		fmt.Fprintln(os.Stderr, "NewConverter: config")
		fmt.Fprintln(os.Stderr, "config.File", config.File)
	}

	index, unit, err := clangutils.CreateTranslationUnit(&clangutils.Config{
		File:    config.File,
		Temp:    false,
		Args:    config.Args,
		IsCpp:   config.IsCpp,
		Options: clang.DetailedPreprocessingRecord,
	})
	if err != nil {
		return nil, err
	}

	return &Converter{
		index: index,
		unit:  unit,
		file:  &ast.File{},
	}, nil
}

func (ct *Converter) Dispose() {
	ct.logln("Dispose")
	ct.index.Dispose()
	ct.unit.Dispose()
}

func (ct *Converter) GetTokens(cursor clang.Cursor) []*ast.Token {
	ran := cursor.Extent()
	var numTokens c.Uint
	var tokens *clang.Token
	ct.unit.Tokenize(ran, &tokens, &numTokens)
	defer ct.unit.DisposeTokens(tokens, numTokens)

	tokensSlice := unsafe.Slice(tokens, int(numTokens))

	result := make([]*ast.Token, 0, int(numTokens))
	for _, tok := range tokensSlice {
		tokStr := ct.unit.Token(tok)
		result = append(result, &ast.Token{
			Token: toToken(tok),
			Lit:   c.GoString(tokStr.CStr()),
		})
		tokStr.Dispose()
	}
	return result
}

func (ct *Converter) logBase() string {
	return strings.Repeat(" ", ct.indent)
}

func (ct *Converter) incIndent() {
	ct.indent++
}

func (ct *Converter) decIndent() {
	if ct.indent > 0 {
		ct.indent--
	}
}

func (ct *Converter) logf(format string, args ...interface{}) {
	if debugParse {
		fmt.Fprintf(os.Stderr, ct.logBase()+format, args...)
	}
}
func (ct *Converter) logln(args ...interface{}) {
	if debugParse {
		if len(args) > 0 {
			firstArg := fmt.Sprintf("%s%v", ct.logBase(), args[0])
			fmt.Fprintln(os.Stderr, append([]interface{}{firstArg}, args[1:]...)...)
		} else {
			fmt.Fprintln(os.Stderr, ct.logBase())
		}
	}
}

func (ct *Converter) InFile(cursor clang.Cursor) bool {
	loc := cursor.Location()
	filePath, _, _ := clangutils.GetPresumedLocation(loc)
	ct.logf("GetCurFile: PresumedLocation %s cursor.Location() %s\n", filePath, clang.GoString(loc.File().FileName()))
	if filePath == "<built-in>" || filePath == "<command line>" {
		//todo(zzy): For some built-in macros, there is no file.
		ct.logln("GetCurFile: NO FILE")
		return false
	}
	return true
}

func (ct *Converter) CreateObject(cursor clang.Cursor, name *ast.Ident) ast.Object {
	base := ast.Object{
		Loc:    createLoc(cursor),
		Parent: ct.BuildScopingExpr(cursor.SemanticParent()),
		Name:   name,
	}
	commentGroup, isDoc := ct.ParseCommentGroup(cursor)
	if isDoc {
		base.Doc = commentGroup
	}
	return base
}

func createLoc(cursor clang.Cursor) *ast.Location {
	filename, _, _ := clangutils.GetPresumedLocation(cursor.Location())
	return &ast.Location{
		File: filename,
	}
}

// extracts and parses comments associated with a given Clang cursor,
// distinguishing between documentation comments and line comments.
//
// The function determines whether a comment is a documentation comment or a line comment by
// comparing the range of the comment node with the range of the declaration node in the AST.
//
// Note: In cases where both documentation comments and line comments conceptually exist,
// only the line comment will be preserved.
func (ct *Converter) ParseCommentGroup(cursor clang.Cursor) (comentGroup *ast.CommentGroup, isDoc bool) {
	rawComment := toStr(cursor.RawCommentText())
	commentGroup := &ast.CommentGroup{}
	if rawComment != "" {
		commentRange := cursor.CommentRange()
		cursorRange := cursor.Extent()
		isDoc := getOffset(commentRange.RangeStart()) < getOffset(cursorRange.RangeStart())
		commentGroup = ct.ParseComment(rawComment)
		if len(commentGroup.List) > 0 {
			return commentGroup, isDoc
		}
	}
	return nil, false
}

func (ct *Converter) ParseComment(rawComment string) *ast.CommentGroup {
	lines := strings.Split(rawComment, "\n")
	commentGroup := &ast.CommentGroup{}
	for _, line := range lines {
		commentGroup.List = append(commentGroup.List, &ast.Comment{Text: line + "\n"})
	}
	return commentGroup
}

// visit top decls (struct, class, function, enum & macro, include)
func (ct *Converter) visitTop(cursor, parent clang.Cursor) clang.ChildVisitResult {
	ct.incIndent()
	defer ct.decIndent()

	inFile := ct.InFile(cursor)

	name := toStr(cursor.String())
	ct.logf("visitTop: Cursor: %s\n", name)

	if !inFile {
		return clang.ChildVisit_Continue
	}

	switch cursor.Kind {
	case clang.CursorInclusionDirective:
		include, err := ct.ProcessInclude(cursor)
		if err != nil {
			ct.logln(err)
			return clang.ChildVisit_Continue
		}
		ct.file.Includes = append(ct.file.Includes, include)
		ct.logln("visitTop: ProcessInclude END ", include.Path)
	case clang.CursorMacroDefinition:
		macro := ct.ProcessMacro(cursor)
		if cursor.IsMacroBuiltin() == 0 {
			ct.file.Macros = append(ct.file.Macros, macro)
		}
		ct.logln("visitTop: ProcessMacro END ", macro.Name, "Tokens Length:", len(macro.Tokens))
	case clang.CursorEnumDecl:
		enum := ct.ProcessEnumDecl(cursor)
		ct.file.Decls = append(ct.file.Decls, enum)
		ct.logf("visitTop: ProcessEnumDecl END")
		if enum.Name != nil {
			ct.logln(enum.Name.Name)
		} else {
			ct.logln("ANONY")
		}

	case clang.CursorClassDecl:
		classDecl := ct.ProcessClassDecl(cursor)
		// todo(zzy):class need consider nested struct situation
		ct.file.Decls = append(ct.file.Decls, classDecl)
		// class havent anonymous situation
		ct.logln("visitTop: ProcessClassDecl END", classDecl.Name.Name)
	case clang.CursorStructDecl:
		decls := ct.ProcessStructDecl(cursor)
		ct.file.Decls = append(ct.file.Decls, decls...)
		ct.logf("visitTop: ProcessStructDecl END")
	case clang.CursorUnionDecl:
		decls := ct.ProcessUnionDecl(cursor)
		ct.file.Decls = append(ct.file.Decls, decls...)
		ct.logf("visitTop: ProcessUnionDecl END")
	case clang.CursorFunctionDecl, clang.CursorCXXMethod, clang.CursorConstructor, clang.CursorDestructor:
		// Handle functions and class methods (including out-of-class method)
		// Example: void MyClass::myMethod() { ... } out-of-class method
		funcDecl := ct.ProcessFuncDecl(cursor)
		ct.file.Decls = append(ct.file.Decls, funcDecl)
		ct.logln("visitTop: ProcessFuncDecl END", funcDecl.Name.Name, funcDecl.MangledName, "isStatic:", funcDecl.IsStatic, "isInline:", funcDecl.IsInline)
	case clang.CursorTypedefDecl:
		typedefDecl := ct.ProcessTypeDefDecl(cursor)
		if typedefDecl == nil {
			return clang.ChildVisit_Continue
		}
		ct.file.Decls = append(ct.file.Decls, typedefDecl)
		ct.logln("visitTop: ProcessTypeDefDecl END", typedefDecl.Name.Name)
	case clang.CursorNamespace:
		clangutils.VisitChildren(cursor, ct.visitTop)
	}
	return clang.ChildVisit_Continue
}

// for flatten ast, keep type order
// input is clang -E 's result
func (ct *Converter) Convert() (*ast.File, error) {
	cursor := ct.unit.Cursor()
	clangutils.VisitChildren(cursor, ct.visitTop)
	return ct.file, nil
}

func (ct *Converter) ProcessType(t clang.Type) ast.Expr {
	ct.incIndent()
	defer ct.decIndent()

	typeName, typeKind := getTypeDesc(t)
	ct.logln("ProcessType: TypeName:", typeName, "TypeKind:", typeKind)

	if t.Kind == clang.TypeUnexposed {
		// https://github.com/goplus/llcppg/issues/497
		return ct.ProcessType(t.CanonicalType())
	}

	if t.Kind >= clang.TypeFirstBuiltin && t.Kind <= clang.TypeLastBuiltin {
		return ct.ProcessBuiltinType(t)
	}

	if t.Kind == clang.TypeElaborated {
		return ct.ProcessElaboratedType(t)
	}

	if t.Kind == clang.TypeTypedef {
		return ct.ProcessTypeDefType(t)
	}

	var expr ast.Expr
	switch t.Kind {
	case clang.TypePointer:
		name, kind := getTypeDesc(t.PointeeType())
		ct.logln("ProcessType: PointerType  Pointee TypeName:", name, "TypeKind:", kind)
		expr = &ast.PointerType{X: ct.ProcessType(t.PointeeType())}
	case clang.TypeBlockPointer:
		name, kind := getTypeDesc(t)
		ct.logln("ProcessType: BlockPointerType  TypeName:", name, "TypeKind:", kind)
		typ := ct.ProcessType(t.PointeeType())
		fnType, ok := typ.(*ast.FuncType)
		if !ok {
			panic("BlockPointerType: not FuncType")
		}
		expr = &ast.BlockPointerType{X: fnType}
	case clang.TypeLValueReference:
		name, kind := getTypeDesc(t.NonReferenceType())
		ct.logln("ProcessType: LvalueRefType  NonReference TypeName:", name, "TypeKind:", kind)
		expr = &ast.LvalueRefType{X: ct.ProcessType(t.NonReferenceType())}
	case clang.TypeRValueReference:
		name, kind := getTypeDesc(t.NonReferenceType())
		ct.logln("ProcessType: RvalueRefType  NonReference TypeName:", name, "TypeKind:", kind)
		expr = &ast.RvalueRefType{X: ct.ProcessType(t.NonReferenceType())}
	case clang.TypeFunctionProto, clang.TypeFunctionNoProto:
		// treating TypeFunctionNoProto as a general function without parameters
		// function type will only collect return type, params will be collected in ProcessFuncDecl
		name, kind := getTypeDesc(t)
		ct.logln("ProcessType: FunctionType  TypeName:", name, "TypeKind:", kind)
		expr = ct.ProcessFunctionType(t)
	case clang.TypeConstantArray, clang.TypeIncompleteArray, clang.TypeVariableArray, clang.TypeDependentSizedArray:
		if t.Kind == clang.TypeConstantArray {
			len := (*c.Char)(c.Malloc(unsafe.Sizeof(c.Char(0)) * 20))
			c.Sprintf(len, c.Str("%lld"), t.ArraySize())
			defer c.Free(unsafe.Pointer(len))
			expr = &ast.ArrayType{
				Elt: ct.ProcessType(t.ArrayElementType()),
				Len: &ast.BasicLit{Kind: ast.IntLit, Value: c.GoString(len)},
			}
		} else if t.Kind == clang.TypeIncompleteArray {
			// incomplete array havent len expr
			expr = &ast.ArrayType{
				Elt: ct.ProcessType(t.ArrayElementType()),
			}
		}
	default:
		name, kind := getTypeDesc(t)
		ct.logln("ProcessType: Unknown Type TypeName:", name, "TypeKind:", kind)
	}
	return expr
}

// For function types, we can only obtain the parameter types, but not the parameter names.
// This is because we cannot reverse-lookup the corresponding declaration node from a function type.
// Note: For function declarations, parameter names are collected in the ProcessFuncDecl method.
func (ct *Converter) ProcessFunctionType(t clang.Type) *ast.FuncType {
	ct.incIndent()
	defer ct.decIndent()
	typeName, typeKind := getTypeDesc(t)
	ct.logln("ProcessFunctionType: TypeName:", typeName, "TypeKind:", typeKind)
	// Note: Attempting to get the type declaration for a function type will result in CursorNoDeclFound
	// cursor := t.TypeDeclaration()
	// This would return CursorNoDeclFound
	resType := t.ResultType()

	name, kind := getTypeDesc(resType)
	ct.logln("ProcessFunctionType: ResultType TypeName:", name, "TypeKind:", kind)

	ret := ct.ProcessType(resType)
	params := &ast.FieldList{}
	numArgs := t.NumArgTypes()
	for i := 0; i < int(numArgs); i++ {
		argType := t.ArgType(c.Uint(i))
		params.List = append(params.List, &ast.Field{
			Type: ct.ProcessType(argType),
		})
	}
	if t.IsFunctionTypeVariadic() != 0 {
		params.List = append(params.List, &ast.Field{
			Type: &ast.Variadic{},
		})
	}

	return &ast.FuncType{
		Ret:    ret,
		Params: params,
	}
}

func (ct *Converter) ProcessTypeDefDecl(cursor clang.Cursor) *ast.TypedefDecl {
	ct.incIndent()
	defer ct.decIndent()
	name, kind := getCursorDesc(cursor)
	ct.logln("ProcessTypeDefDecl: CursorName:", name, "CursorKind:", kind, "CursorTypeKind:", toStr(cursor.Type().Kind.String()))

	typ := ct.ProcessUnderlyingType(cursor)
	// For cases like: typedef struct { int x; } Name;
	// libclang incorrectly reports the anonymous structure as a named structure
	// with the same name as the typedef. Since the anonymous structure definition
	// has already been collected when processing its declaration cursor,
	// we skip this redundant typedef declaration by returning nil.
	if typ == nil {
		return nil
	}

	decl := &ast.TypedefDecl{
		Object: ct.CreateObject(cursor, &ast.Ident{Name: name}),
		Type:   typ,
	}
	return decl
}

func (ct *Converter) ProcessUnderlyingType(cursor clang.Cursor) ast.Expr {
	underlyingTyp := cursor.TypedefDeclUnderlyingType()

	if underlyingTyp.Kind != clang.TypeElaborated {
		ct.logln("ProcessUnderlyingType: not elaborated")
		return ct.ProcessType(underlyingTyp)
	}

	defName := toStr(cursor.String())
	// Using getActualTypeCursor to recursively find the actual declaration of the underlying type,
	// handles cases with multi-level typedef chains
	underName := toStr(ct.getActualTypeCursor(underlyingTyp.TypeDeclaration()).String())
	ct.logln("ProcessUnderlyingType: defName:", defName, "underName:", underName)

	// For a typedef like "typedef struct xxx xxx;", the underlying type declaration
	// can appear in two locations:
	// 1. Inside the typedef itself when the struct is defined inline
	// 2. At the implementation location when there's a separate struct xxx definition
	// in the source file
	// Therefore, we shouldn't use declaration location to determine whether to remove
	// extra typedef nodes
	//
	// Note: This handles both direct self-references (e.g., typedef struct Foo Foo;) and
	// multi-level typedef chains that refer back to the original declaration (e.g., typedef enum algorithm {...} algorithm_t; typedef algorithm_t algorithm;)
	if defName == underName {
		ct.logln("ProcessUnderlyingType: is self reference")
		return nil
	}

	return ct.ProcessElaboratedType(underlyingTyp)
}

// getActualType gets the actual type by handling only the outer Elaborated and Typedef types.
// Note: We don't use CanonicalType() because it recursively resolves all types, including parameter types
// in function signatures, which would cause typedef types within function signatures to be desugared.
// For example:
//
//	typedef struct OSSL_CORE_HANDLE OSSL_CORE_HANDLE;
//	typedef struct OSSL_DISPATCH OSSL_DISPATCH;
//	typedef int (OSSL_provider_init_fn)(const OSSL_CORE_HANDLE *handle,
//	                            const OSSL_DISPATCH *in,
//	                            const OSSL_DISPATCH **out,
//	                            void **provctx);
//	OSSL_provider_init_fn OSSL_provider_init;
//
// Using CanonicalType() would desugar OSSL_CORE_HANDLE and OSSL_DISPATCH in the function signature
// to their underlying struct types, which is not what we want.
func (ct *Converter) getActualType(t clang.Type) clang.Type {
	ct.incIndent()
	defer ct.decIndent()
	typName, typKind := getTypeDesc(t)
	ct.logln("getActualType: TypeName:", typName, "TypeKind:", typKind)
	switch t.Kind {
	case clang.TypeElaborated:
		ct.logln("getActualType: TypeElaborated")
		return ct.getActualType(t.NamedType())
	case clang.TypeTypedef:
		ct.logln("getActualType: TypeTypedef")
		return ct.getActualType(t.TypeDeclaration().TypedefDeclUnderlyingType())
	default:
		return t
	}
}

// getActualTypeCursor recursively gets the actual underlying cursor of a type declaration.
// For multi-level nested typedef chains, it continues recursion until finding the original non-typedef type.
// Example:
// - typedef enum Foo {...} Foo_t;
// - typedef Foo_t Foo;
// When processing Foo, it recursively finds the declaration cursor of enum Foo
func (ct *Converter) getActualTypeCursor(cursor clang.Cursor) clang.Cursor {
	ct.incIndent()
	defer ct.decIndent()
	typName, typKind := getCursorDesc(cursor)
	ct.logln("getActualTypeCursor: TypeName:", typName, "TypeKind:", typKind)
	switch cursor.Kind {
	case clang.CursorTypedefDecl:
		return ct.getActualTypeCursor(cursor.TypedefDeclUnderlyingType().TypeDeclaration())
	default:
		return cursor
	}
}

// converts functions, methods, constructors, destructors (including out-of-class decl) to ast.FuncDecl nodes.
func (ct *Converter) ProcessFuncDecl(cursor clang.Cursor) *ast.FuncDecl {
	ct.incIndent()
	defer ct.decIndent()
	name, kind := getCursorDesc(cursor)
	mangledName := toStr(cursor.Mangling())
	ct.logln("ProcessFuncDecl: CursorName:", name, "CursorKind:", kind, "mangledName:", mangledName)

	// function type will only collect return type
	// ProcessType can't get the field names, will collect in follows
	fnType := cursor.Type()
	typName, typKind := getTypeDesc(fnType)
	ct.logln("ProcessFuncDecl: TypeName:", typName, "TypeKind:", typKind)

	typeToProcess := fnType
	if fnType.Kind == clang.TypeElaborated {
		typeToProcess = ct.getActualType(fnType)
		actualTypeName, actualTypeKind := getTypeDesc(typeToProcess)
		ct.logln("ProcessFuncDecl: ActualType TypeName:", actualTypeName, "TypeKind:", actualTypeKind)
	}
	funcType, ok := ct.ProcessType(typeToProcess).(*ast.FuncType)
	if !ok {
		ct.logln("ProcessFuncDecl: failed to process function type")
		return nil
	}
	ct.logln("ProcessFuncDecl: ProcessFieldList")

	// For function type references (e.g. `typedef void (fntype)(); fntype foo;`),
	// params are already processed in ProcessType via CanonicalType
	if fnType.Kind != clang.TypeElaborated {
		numArgs := cursor.NumArguments()
		numFields := c.Int(len(funcType.Params.List))
		for i := c.Int(0); i < numArgs; i++ {
			arg := cursor.Argument(c.Uint(i))
			name := clang.GoString(arg.DisplayName())
			if len(name) > 0 && i < numFields {
				field := funcType.Params.List[i]
				field.Names = []*ast.Ident{&ast.Ident{Name: name}}
			}
		}
	}

	// Linux has one less leading underscore than macOS, so remove one leading underscore on macOS
	if runtime.GOOS == "darwin" {
		mangledName = strings.TrimPrefix(mangledName, "_")
	}

	funcDecl := &ast.FuncDecl{
		Object:      ct.CreateObject(cursor, &ast.Ident{Name: name}),
		Type:        funcType,
		MangledName: mangledName,
	}

	if cursor.IsFunctionInlined() != 0 {
		funcDecl.IsInline = true
	}

	if isMethod(cursor) {
		ct.logln("ProcessFuncDecl: is method, ProcessMethodAttributes")
		ct.ProcessMethodAttributes(cursor, funcDecl)
	} else {
		if cursor.StorageClass() == clang.SCStatic {
			funcDecl.IsStatic = true
		}
	}

	return funcDecl
}

// get Methods Attributes
func (ct *Converter) ProcessMethodAttributes(cursor clang.Cursor, fn *ast.FuncDecl) {
	if parent := cursor.SemanticParent(); parent.Equal(cursor.LexicalParent()) != 1 {
		fn.Parent = ct.BuildScopingExpr(cursor.SemanticParent())
	}

	switch cursor.Kind {
	case clang.CursorDestructor:
		fn.IsDestructor = true
	case clang.CursorConstructor:
		fn.IsConstructor = true
		if cursor.IsExplicit() != 0 {
			fn.IsExplicit = true
		}
	}

	if cursor.IsStatic() != 0 {
		fn.IsStatic = true
	}
	if cursor.IsVirtual() != 0 || cursor.IsPureVirtual() != 0 {
		fn.IsVirtual = true
	}
	if cursor.IsConst() != 0 {
		fn.IsConst = true
	}

	var numOverridden c.Uint
	var overridden *clang.Cursor
	cursor.OverriddenCursors(&overridden, &numOverridden)
	if numOverridden > 0 {
		fn.IsOverride = true
	}
	overridden.DisposeOverriddenCursors()
}

func (ct *Converter) ProcessEnumType(cursor clang.Cursor) *ast.EnumType {
	items := make([]*ast.EnumItem, 0)

	clangutils.VisitChildren(cursor, func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Kind == clang.CursorEnumConstantDecl {
			name := cursor.String()
			defer name.Dispose()

			val := (*c.Char)(c.Malloc(unsafe.Sizeof(c.Char(0)) * 20))
			c.Sprintf(val, c.Str("%lld"), cursor.EnumConstantDeclValue())
			defer c.Free(unsafe.Pointer(val))

			enum := &ast.EnumItem{
				Name: &ast.Ident{Name: c.GoString(name.CStr())},
				Value: &ast.BasicLit{
					Kind:  ast.IntLit,
					Value: c.GoString(val),
				},
			}
			items = append(items, enum)
		}
		return clang.ChildVisit_Continue
	})

	return &ast.EnumType{
		Items: items,
	}
}

func (ct *Converter) ProcessEnumDecl(cursor clang.Cursor) *ast.EnumTypeDecl {
	cursorName, cursorKind := getCursorDesc(cursor)
	ct.logln("ProcessEnumDecl: CursorName:", cursorName, "CursorKind:", cursorKind)

	decl := &ast.EnumTypeDecl{
		Object: ct.CreateObject(cursor, nil),
		Type:   ct.ProcessEnumType(cursor),
	}

	anony := cursor.IsAnonymous()
	if anony == 0 {
		decl.Name = &ast.Ident{Name: cursorName}
		ct.logln("ProcessEnumDecl: has name", cursorName)
	} else {
		ct.logln("ProcessEnumDecl: is anonymous")
	}

	return decl
}

// current only collect macro which defined in file
func (ct *Converter) ProcessMacro(cursor clang.Cursor) *ast.Macro {
	macro := &ast.Macro{
		Loc:    createLoc(cursor),
		Name:   clang.GoString(cursor.String()),
		Tokens: ct.GetTokens(cursor),
	}
	return macro
}

func (ct *Converter) ProcessInclude(cursor clang.Cursor) (*ast.Include, error) {
	name := toStr(cursor.String())
	includedFile := cursor.IncludedFile()
	includedPath := toStr(includedFile.FileName())
	if includedPath == "" {
		return nil, fmt.Errorf("%s: failed to get included file", name)
	}
	return &ast.Include{Path: filepath.Clean(includedPath)}, nil
}

func (ct *Converter) createBaseField(cursor clang.Cursor) *ast.Field {
	ct.incIndent()
	defer ct.decIndent()

	fieldName := toStr(cursor.String())

	typ := cursor.Type()

	typeName, typeKind := getTypeDesc(typ)

	ct.logf("createBaseField: ProcessType %s TypeKind: %s", typeName, typeKind)

	field := &ast.Field{
		Type: ct.ProcessType(typ),
	}

	commentGroup, isDoc := ct.ParseCommentGroup(cursor)
	if commentGroup != nil {
		if isDoc {
			field.Doc = commentGroup
		} else {
			field.Comment = commentGroup
		}
	}
	// NOTE(MeteorsLiu): In non C++ mode, an anonymous field name may be `unname struct` instead of empty string
	// so check it via IsAnonymous()
	if cursor.IsAnonymous() == 0 {
		field.Names = []*ast.Ident{{Name: fieldName}}
	}
	return field
}

// For Record Type(struct, union ...)'s FieldList
func (ct *Converter) ProcessFieldList(cursor clang.Cursor) *ast.FieldList {
	ct.incIndent()
	defer ct.decIndent()
	flds := &ast.FieldList{}
	ct.logln("ProcessFieldList: VisitChildren")
	clangutils.VisitChildren(cursor, func(subcsr, parent clang.Cursor) clang.ChildVisitResult {
		switch subcsr.Kind {
		case clang.CursorFieldDecl:
			// In C language, parameter lists do not have similar parameter grouping in Go.
			// func foo(a, b int)

			// For follows struct, it will also parse to two FieldDecl
			// struct A {
			// 	int a, b;
			// };
			ct.logln("ProcessFieldList: CursorFieldDecl")
			field := ct.createBaseField(subcsr)
			field.Access = ast.AccessSpecifier(subcsr.CXXAccessSpecifier())
			flds.List = append(flds.List, field)
		case clang.CursorVarDecl:
			if subcsr.StorageClass() == clang.SCStatic {
				// static member variable
				field := ct.createBaseField(subcsr)
				field.Access = ast.AccessSpecifier(subcsr.CXXAccessSpecifier())
				field.IsStatic = true
				flds.List = append(flds.List, field)
			}
		}
		return clang.ChildVisit_Continue
	})
	return flds
}

// Note:Public Method is considered
func (ct *Converter) ProcessMethods(cursor clang.Cursor) []*ast.FuncDecl {
	methods := make([]*ast.FuncDecl, 0)
	clangutils.VisitChildren(cursor, func(subcsr, parent clang.Cursor) clang.ChildVisitResult {
		if isMethod(subcsr) && subcsr.CXXAccessSpecifier() == clang.CXXPublic {
			method := ct.ProcessFuncDecl(subcsr)
			if method != nil {
				methods = append(methods, method)
			}
		}
		return clang.ChildVisit_Continue
	})
	return methods
}

func (ct *Converter) ProcessRecordDecl(cursor clang.Cursor) []ast.Decl {
	var decls []ast.Decl
	ct.incIndent()
	defer ct.decIndent()
	cursorName, cursorKind := getCursorDesc(cursor)
	ct.logln("ProcessRecordDecl: CursorName:", cursorName, "CursorKind:", cursorKind)

	childs := PostOrderVisitChildren(cursor, func(child, parent clang.Cursor) bool {
		// if we found a nested enum, handle it like nested struct
		if child.Kind == clang.CursorEnumDecl {
			return true
		}
		return (child.Kind == clang.CursorStructDecl || child.Kind == clang.CursorUnionDecl) && child.IsAnonymous() == 0
	})

	for _, child := range childs {
		switch child.Kind {
		case clang.CursorStructDecl, clang.CursorUnionDecl:
			// note(zzy):use len(typ.Fields.List) to ensure it has fields not a forward declaration
			// but maybe make the forward decl in to AST is also good.
			childName := clang.GoString(child.String())
			ct.logln("ProcessRecordDecl: Found named nested struct:", childName)
			// Check if this is a named nested struct/union
			typ := ct.ProcessRecordType(child)
			// note(zzy):use len(typ.Fields.List) to ensure it has fields not a forward declaration
			// but maybe make the forward decl in to AST is also good.
			if child.IsAnonymous() == 0 && typ.Fields != nil {
				decls = append(decls, &ast.TypeDecl{
					Object: ct.CreateObject(child, &ast.Ident{Name: childName}),
					Type:   ct.ProcessRecordType(child),
				})
			}
		case clang.CursorEnumDecl:
			childName := clang.GoString(child.String())

			ct.logln("ProcessRecordDecl: Found named nested enum:", childName)

			ct.incIndent()
			decls = append(decls, ct.ProcessEnumDecl(child))
			ct.decIndent()
		}
	}
	ct.logln("ProcessRecordDecl: process record: ", cursorName)

	decl := &ast.TypeDecl{
		Object: ct.CreateObject(cursor, nil),
		Type:   ct.ProcessRecordType(cursor),
	}

	// NOTE(MeteorsLiu): IsAnonymousRecordDecl may return fake results when we're in non Cpp mode
	// to avoid that case, we have to check the IsAnonymous result
	isAnonymousRecord := cursor.IsAnonymousRecordDecl() > 0 || cursor.IsAnonymous() > 0

	if !isAnonymousRecord {
		decl.Name = &ast.Ident{Name: cursorName}
		ct.logln("ProcessRecordDecl: has name", cursorName)
	} else {
		ct.logln("ProcessRecordDecl: is anonymous")
	}

	decls = append(decls, decl)
	return decls
}

func (ct *Converter) ProcessStructDecl(cursor clang.Cursor) []ast.Decl {
	return ct.ProcessRecordDecl(cursor)
}

func (ct *Converter) ProcessUnionDecl(cursor clang.Cursor) []ast.Decl {
	return ct.ProcessRecordDecl(cursor)
}

func (ct *Converter) ProcessClassDecl(cursor clang.Cursor) *ast.TypeDecl {
	cursorName, cursorKind := getCursorDesc(cursor)
	ct.logln("ProcessClassDecl: CursorName:", cursorName, "CursorKind:", cursorKind)

	// Pushing class scope before processing its type and popping after
	base := ct.CreateObject(cursor, &ast.Ident{Name: cursorName})
	typ := ct.ProcessRecordType(cursor)

	decl := &ast.TypeDecl{
		Object: base,
		Type:   typ,
	}

	return decl
}

func (ct *Converter) ProcessRecordType(cursor clang.Cursor) *ast.RecordType {
	ct.incIndent()
	defer ct.decIndent()

	typ := &ast.RecordType{}

	cursorName, cursorKind := getCursorDesc(cursor)
	ct.logln("ProcessRecordType: CursorName:", cursorName, "CursorKind:", cursorKind)

	typ.Tag = toTag(cursor.Kind)
	ct.logln("ProcessRecordType: toTag", typ.Tag)

	if cursor.IsCursorDefinition() == 0 {
		ct.logln("ProcessRecordType: forward declaration, no definition")
		return typ
	}

	ct.logln("ProcessRecordType: ProcessFieldList")
	typ.Fields = ct.ProcessFieldList(cursor)

	ct.logln("ProcessRecordType: ProcessMethods")
	typ.Methods = ct.ProcessMethods(cursor)

	return typ
}

// process ElaboratedType Reference
//
// 1. Named elaborated type references:
// - Examples: struct MyStruct, union MyUnion, class MyClass, enum MyEnum
// - Handling: Constructed as TagExpr or ScopingExpr references
//
// 2. Anonymous elaborated type references:
// - Examples: struct { int x; int y; }, union { int a; float b; }
// - Handling: Retrieve their corresponding concrete types
func (ct *Converter) ProcessElaboratedType(t clang.Type) ast.Expr {
	ct.incIndent()
	defer ct.decIndent()
	typeName, typeKind := getTypeDesc(t)
	ct.logln("ProcessElaboratedType: TypeName:", typeName, "TypeKind:", typeKind)

	decl := t.TypeDeclaration()
	isAnonymousDecl := decl.IsAnonymous() > 0

	if isAnonymousDecl && decl.Kind != clang.CursorEnumDecl {
		return ct.ProcessRecordType(decl)
	}
	parts := clangutils.BuildScopingParts(decl)
	hasParent := clangutils.HasParent(decl)
	// NOTE(MeteorsLiu): nested enum behaves different from nested struct, for example, we can find its semantic parent
	// however, it will cause we misidentified it as a class method expr, so take it out
	if isAnonymousDecl && decl.Kind == clang.CursorEnumDecl {
		// case 1: anonymous enum, but not nested (anonymous enum decl variable case)
		if !hasParent {
			// this is not a nested enum, handle it normally
			return ct.ProcessEnumType(decl)
		}
		// case 2: anonymous enum, nested (normal nested struct reference)
		// by default, the type of an anonymous enum is int
		// NOTE(MeteorsLiu): see disscussion https://github.com/goplus/llcppg/pull/530
		return &ast.BuiltinType{Kind: ast.Int}

		// case 3: named enum, nested, fallback to process as a ElaboratedType (nornaml nested struct)
		// case 4: named enum, non-nested, fallback to process as a ElaboratedType normally. (typedef enum case)
	}

	// for elaborated type, it could have a tag description
	// like struct A, union B, class C, enum D
	typeParts := strings.SplitN(typeName, " ", 2)

	if len(typeParts) == 2 {
		if tagValue, ok := tagMap[typeParts[0]]; ok {
			return &ast.TagExpr{
				Tag:  tagValue,
				Name: buildScopingFromParts(parts),
			}
		}
	}

	return buildScopingFromParts(parts)
}

func (ct *Converter) ProcessTypeDefType(t clang.Type) ast.Expr {
	cursor := t.TypeDeclaration()
	ct.logln("ProcessTypeDefType: Typedef TypeDeclaration", toStr(cursor.String()), toStr(t.String()))
	if name := toStr(cursor.String()); name != "" {
		return &ast.Ident{Name: name}
	}
	ct.logln("ProcessTypeDefType: typedef type have no name")
	return nil
}

func (ct *Converter) ProcessBuiltinType(t clang.Type) *ast.BuiltinType {
	ct.incIndent()
	defer ct.decIndent()
	typeName, typeKind := getTypeDesc(t)
	ct.logln("ProcessBuiltinType: TypeName:", typeName, "TypeKind:", typeKind)

	kind := ast.Void
	var flags ast.TypeFlag

	switch t.Kind {
	case clang.TypeVoid:
		kind = ast.Void
	case clang.TypeBool:
		kind = ast.Bool
	case clang.TypeCharU, clang.TypeUChar, clang.TypeCharS, clang.TypeSChar:
		kind = ast.Char
	case clang.TypeChar16:
		kind = ast.Char16
	case clang.TypeChar32:
		kind = ast.Char32
	case clang.TypeWChar:
		kind = ast.WChar
	case clang.TypeShort, clang.TypeUShort:
		kind = ast.Int
		flags |= ast.Short
	case clang.TypeInt, clang.TypeUInt:
		kind = ast.Int
	case clang.TypeLong, clang.TypeULong:
		kind = ast.Int
		flags |= ast.Long
	case clang.TypeLongLong, clang.TypeULongLong:
		kind = ast.Int
		flags |= ast.LongLong
	case clang.TypeInt128, clang.TypeUInt128:
		kind = ast.Int128
	case clang.TypeFloat:
		kind = ast.Float
	case clang.TypeHalf, clang.TypeFloat16:
		kind = ast.Float16
	case clang.TypeDouble:
		kind = ast.Float
		flags |= ast.Double
	case clang.TypeLongDouble:
		kind = ast.Float
		flags |= ast.Long | ast.Double
	case clang.TypeFloat128:
		kind = ast.Float128
	case clang.TypeComplex:
		kind = ast.Complex
		complexKind := t.ElementType().Kind
		if complexKind == clang.TypeLongDouble {
			flags |= ast.Long | ast.Double
		} else if complexKind == clang.TypeDouble {
			flags |= ast.Double
		}
		// float complfex flag is not set
	default:
		// like IBM128,NullPtr,Accum
		kindStr := toStr(t.Kind.String())
		fmt.Fprintln(os.Stderr, "todo: unknown builtin type:", kindStr)
	}

	if IsExplicitSigned(t) {
		flags |= ast.Signed
	} else if IsExplicitUnsigned(t) {
		flags |= ast.Unsigned
	}

	return &ast.BuiltinType{
		Kind:  kind,
		Flags: flags,
	}
}

// Constructs a complete scoping expression by traversing the semantic parents, starting from the given clang.Cursor
// For anonymous decl of typedef references, use their anonymous name
func (ct *Converter) BuildScopingExpr(cursor clang.Cursor) ast.Expr {
	parts := clangutils.BuildScopingParts(cursor)
	return buildScopingFromParts(parts)
}

func PostOrderVisitChildren(cursor clang.Cursor, collect func(c, p clang.Cursor) bool) []clang.Cursor {
	var children []clang.Cursor
	clangutils.VisitChildren(cursor, func(child, parent clang.Cursor) clang.ChildVisitResult {
		if collect(child, parent) {
			childs := PostOrderVisitChildren(child, collect)
			children = append(children, childs[:]...)
			children = append(children, child)
		}
		return clang.ChildVisit_Continue
	})
	return children
}

func IsExplicitSigned(t clang.Type) bool {
	return t.Kind == clang.TypeCharS || t.Kind == clang.TypeSChar
}

func IsExplicitUnsigned(t clang.Type) bool {
	return t.Kind == clang.TypeCharU || t.Kind == clang.TypeUChar ||
		t.Kind == clang.TypeUShort || t.Kind == clang.TypeUInt ||
		t.Kind == clang.TypeULong || t.Kind == clang.TypeULongLong ||
		t.Kind == clang.TypeUInt128
}

func toTag(kind clang.CursorKind) ast.Tag {
	switch kind {
	case clang.CursorStructDecl:
		return ast.Struct
	case clang.CursorUnionDecl:
		return ast.Union
	case clang.CursorClassDecl:
		return ast.Class
	default:
		panic(fmt.Sprintf("Unexpected cursor kind in toTag: %v", kind))
	}
}

func toToken(tok clang.Token) token.Token {
	if tok.Kind() < clang.Punctuation || tok.Kind() > clang.Comment {
		return token.ILLEGAL
	} else {
		return token.Token(tok.Kind() + 1)
	}
}
func isMethod(cursor clang.Cursor) bool {
	return cursor.Kind == clang.CursorCXXMethod || cursor.Kind == clang.CursorConstructor || cursor.Kind == clang.CursorDestructor
}

func buildScopingFromParts(parts []string) ast.Expr {
	if len(parts) == 0 {
		return nil
	}
	var expr ast.Expr = &ast.Ident{Name: parts[0]}
	for _, part := range parts[1:] {
		expr = &ast.ScopingExpr{
			Parent: expr,
			X:      &ast.Ident{Name: part},
		}
	}
	return expr
}

func getOffset(location clang.SourceLocation) c.Uint {
	_, _, _, offset := clangutils.GetLocation(location)
	return offset
}

func toStr(clangStr clang.String) (str string) {
	defer clangStr.Dispose()
	if clangStr.CStr() != nil {
		str = c.GoString(clangStr.CStr())
	}
	return
}

func getTypeDesc(t clang.Type) (name string, kind string) {
	name = toStr(t.String())
	kind = toStr(t.Kind.String())
	return
}

func getCursorDesc(cursor clang.Cursor) (name string, kind string) {
	name = toStr(cursor.String())
	kind = toStr(cursor.Kind.String())
	return
}
