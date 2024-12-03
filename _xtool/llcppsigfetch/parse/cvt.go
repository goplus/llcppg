package parse

import (
	"fmt"
	"os"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llgo/c"
	"github.com/goplus/llgo/c/clang"
)

type FileEntry struct {
	Path    string
	IncPath string
	IsSys   bool
	Doc     *ast.File
}

type Converter struct {
	Files     []*FileEntry
	FileOrder []string // todo(zzy): more efficient struct
	curLoc    ast.Location
	index     *clang.Index
	unit      *clang.TranslationUnit

	indent int // for verbose debug
}

type Config struct {
	File  string
	Temp  bool
	Args  []string
	IsCpp bool
}

func NewConverter(config *clangutils.Config) (*Converter, error) {
	index, unit, err := clangutils.CreateTranslationUnit(config)
	if err != nil {
		return nil, err
	}

	files := initFileEntries(unit)

	return &Converter{
		Files: files,
		index: index,
		unit:  unit,
	}, nil

}

func (ct *Converter) Dispose() {
	ct.logln("Dispose")
	ct.index.Dispose()
	ct.unit.Dispose()
}

func initFileEntries(unit *clang.TranslationUnit) []*FileEntry {
	files := make([]*FileEntry, 0)
	clangutils.GetInclusions(unit, func(inced clang.File, incins []clang.SourceLocation) {
		loc := unit.GetLocation(inced, 1, 1)
		incedFile := toStr(inced.FileName())
		var incPath string
		if len(incins) > 0 {
			cur := unit.GetCursor(&incins[0])
			incPath = toStr(cur.String())
		}
		files = append(files, &FileEntry{
			Path:    incedFile,
			IncPath: incPath,
			IsSys:   loc.IsInSystemHeader() != 0,
			Doc:     &ast.File{},
		})
	})
	return files
}

func (ct *Converter) logBase() string {
	return strings.Repeat(" ", ct.indent)
}

func (ct *Converter) logln(args ...interface{}) {
	if debugParse {
		if len(args) > 0 {
			firstArg := fmt.Sprintf("%s%v", ct.logBase(), args[0])
			fmt.Fprintln(os.Stderr, append([]interface{}{firstArg}, args[1:]...)...)
		}
	}
}

func (ct *Converter) GetCurFile(cursor clang.Cursor) *ast.File {
	loc := cursor.Location()
	var file clang.File
	loc.SpellingLocation(&file, nil, nil, nil)

	if file.FileName().CStr() == nil {
		ct.curLoc = ast.Location{File: ""}
		ct.logln("GetCurFile: NO FILE") // 这个注释掉就可以正常执行，不注释就会崩
		return nil
	}

	filePath := toStr(file.FileName())

	ct.curLoc = ast.Location{File: filePath}

	// todo(zzy): more efficient
	for i, entry := range ct.Files {
		if entry.Path == filePath {
			ct.logln("GetCurFile: found", filePath)
			return ct.Files[i].Doc
		}
	}
	// ct.logln("GetCurFile: Create New ast.File", filePath)
	// entry := &FileEntry{Path: filePath, Doc: &ast.File{}, IsSys: false}
	// if loc.IsInSystemHeader() != 0 {
	// 	entry.IsSys = true
	// }
	// ct.Files = append(ct.Files, entry)
	return nil
}

// visit top decls (struct,class,function,enum & macro,include)
func (ct *Converter) visitTop(cursor, parent clang.Cursor) clang.ChildVisitResult {
	curFile := ct.GetCurFile(cursor)
	if curFile == nil {
		return clang.ChildVisit_Continue
	}
	switch cursor.Kind {
	case clang.CursorFunctionDecl:
		funcDecl := ct.ProcessFuncDecl(cursor)
		curFile.Decls = append(curFile.Decls, funcDecl)
		ct.logln("visitTop: ProcessFuncDecl END", funcDecl.Name.Name, funcDecl.MangledName, "isStatic:", funcDecl.IsStatic, "isInline:", funcDecl.IsInline)
	}
	return clang.ChildVisit_Continue
}

func (ct *Converter) Convert() ([]*FileEntry, error) {
	cursor := ct.unit.Cursor()
	// visit top decls (struct,class,function & macro,include)
	clangutils.VisitChildren(cursor, ct.visitTop)
	return ct.Files, nil
}

func (ct *Converter) ProcessFuncDecl(cursor clang.Cursor) *ast.FuncDecl {
	funcDecl := &ast.FuncDecl{DeclBase: ast.DeclBase{Loc: &ast.Location{File: ct.curLoc.File}, Doc: nil, Parent: nil}, Name: &ast.Ident{Name: "cjsonfree"},
		Type: &ast.FuncType{Ret: &ast.BuiltinType{Kind: ast.Void}, Params: &ast.FieldList{
			List: []*ast.Field{{Type: &ast.PointerType{X: &ast.BuiltinType{Kind: ast.Void}}, Names: []*ast.Ident{{Name: "object"}}}}},
		},
		MangledName: "cjsonfree",
	}
	return funcDecl
}

func toStr(clangStr clang.String) (str string) {
	defer clangStr.Dispose() // 无论这个是否存在，只要上方的 ct.logln("GetCurFile: NO FILE") 没有注释掉就会崩溃
	if clangStr.CStr() != nil {
		str = c.GoString(clangStr.CStr())
	}
	return
}
