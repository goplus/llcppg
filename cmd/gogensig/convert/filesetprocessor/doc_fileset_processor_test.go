package filesetprocessor_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/convert/filesetprocessor"
	"github.com/goplus/llcppg/cmd/gogensig/visitor"
	"github.com/goplus/llcppg/llcppg"
)

func TestProcessValidSigfetchContent(t *testing.T) {
	content := []map[string]interface{}{
		{
			"_Type": "FileEntry",
			"path":  "temp.h",
			"doc": map[string]interface{}{
				"_Type": "File",
				"decls": []map[string]interface{}{
					{
						"_Type":  "FuncDecl",
						"Loc":    map[string]interface{}{"_Type": "Location", "File": "temp.h"},
						"Doc":    nil,
						"Parent": nil,
						"Name":   map[string]interface{}{"_Type": "Ident", "Name": "go_func_name"},
						"Type": map[string]interface{}{
							"_Type":  "FuncType",
							"Params": map[string]interface{}{"_Type": "FieldList", "List": []interface{}{}},
							"Ret":    map[string]interface{}{"_Type": "BuiltinType", "Kind": 6, "Flags": 0},
						},
						"IsInline":      false,
						"IsStatic":      false,
						"IsConst":       false,
						"IsExplicit":    false,
						"IsConstructor": false,
						"IsDestructor":  false,
						"IsVirtual":     false,
						"IsOverride":    false,
					},
				},
			},
		},
	}

	tempFileName, err := config.CreateTmpJSONFile("llcppg.sigfetch-test.json", content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFileName)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir, err := os.MkdirTemp(dir, "gogensig-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	p, _, err := filesetprocessor.New(&convert.Config{
		PkgName:   "files",
		SymbFile:  "",
		CfgFile:   "",
		OutputDir: tempDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = p.ProcessFileSetFromPath(tempFileName)
	if err != nil {
		t.Error(err)
	}

	t.Run("process", func(t *testing.T) {
		err = filesetprocessor.Process(&convert.Config{
			PkgName:      "files",
			SymbFile:     "",
			CfgFile:      "",
			OutputDir:    tempDir,
			SigfetchFile: tempFileName,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

}

func TestProcessFileNotExist(t *testing.T) {
	astConvert, err := convert.NewAstConvert(&convert.Config{
		PkgName:  "error",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal(err)
	}
	docVisitors := []visitor.DocVisitor{astConvert}
	manager := visitor.NewDocVisitorList(docVisitors)
	p := filesetprocessor.NewDocFileSetProcessor(&filesetprocessor.ProcesserConfig{
		Exec: func(file *llcppg.FileEntry) error {
			manager.Visit(file.Doc, file.Path, file.IncPath, file.IsSys, file.FileType)
			return nil
		},
		DepIncs: []string{},
	})
	err = p.ProcessFileSetFromPath("notexist.json")
	if !os.IsNotExist(err) {
		t.Error("expect no such file or directory error")
	}
}

func TestProcessInvalidSigfetchContent(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("%s", "expect panic")
		}
	}()

	invalidContent := "invalid json content"
	tempFileName, err := config.CreateTmpJSONFile("llcppg.sigfetch-panic.json", invalidContent)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFileName)

	astConvert, err := convert.NewAstConvert(&convert.Config{
		PkgName:  "panic",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal(err)
	}
	docVisitors := []visitor.DocVisitor{astConvert}
	manager := visitor.NewDocVisitorList(docVisitors)
	p := filesetprocessor.NewDocFileSetProcessor(&filesetprocessor.ProcesserConfig{
		Exec: func(file *llcppg.FileEntry) error {
			manager.Visit(file.Doc, file.Path, file.IncPath, file.IsSys, file.FileType)
			return nil
		},
		DepIncs: []string{},
	})
	err = p.ProcessFileSetFromPath(tempFileName)
	if err != nil {
		panic(err)
	}
}

var errCustomExec = errors.New("custom exec error")

func TestCustomExec(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("%s", "expect panic")
		}
	}()
	file := []*llcppg.FileEntry{
		{
			Path:  "/path/to/foo.h",
			IsSys: false,
			Doc:   &ast.File{},
		},
	}
	p := filesetprocessor.NewDocFileSetProcessor(&filesetprocessor.ProcesserConfig{
		Exec: func(file *llcppg.FileEntry) error {
			return errCustomExec
		},
	})
	err := p.ProcessFileSet(file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecOrder(t *testing.T) {
	depIncs := []string{"/path/to/int16_t.h"}
	fileSet := []*llcppg.FileEntry{
		{
			Path:    "/path/to/foo.h",
			IncPath: "foo.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{
					{Path: "/path/to/cdef.h"},
					{Path: "/path/to/stdint.h"},
				},
			},
		},
		{
			Path:    "/path/to/cdef.h",
			IncPath: "cdef.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{
					{Path: "/path/to/int8_t.h"},
					{Path: "/path/to/int16_t.h"},
				},
			},
		},
		{
			Path:    "/path/to/stdint.h",
			IncPath: "stdint.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{
					{Path: "/path/to/int8_t.h"},
					{Path: "/path/to/int16_t.h"},
				},
			},
		},
		{
			Path:    "/path/to/int8_t.h",
			IncPath: "int8_t.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{},
			},
		},
		{
			Path:    "/path/to/int16_t.h",
			IncPath: "int16_t.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{},
			},
		},
		{
			Path:    "/path/to/bar.h",
			IncPath: "bar.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{
					{Path: "/path/to/stdint.h"},
					{Path: "/path/to/a.h"},
				},
			},
		},
		// circular dependency
		{
			Path:    "/path/to/a.h",
			IncPath: "a.h",
			IsSys:   false,
			Doc: &ast.File{
				Includes: []*ast.Include{
					{Path: "/path/to/bar.h"},
					// will not appear in normal
					{Path: "/path/to/noexist.h"},
				},
			},
		},
	}
	var processFiles []string
	expectedOrder := []string{
		"/path/to/int8_t.h",
		"/path/to/cdef.h",
		"/path/to/stdint.h",
		"/path/to/foo.h",
		"/path/to/a.h",
		"/path/to/bar.h",
	}
	p := filesetprocessor.NewDocFileSetProcessor(&filesetprocessor.ProcesserConfig{
		Exec: func(file *llcppg.FileEntry) error {
			processFiles = append(processFiles, file.Path)
			return nil
		},
		DepIncs: depIncs,
	})
	p.ProcessFileSet(fileSet)
	if !reflect.DeepEqual(processFiles, expectedOrder) {
		t.Errorf("expect %v, got %v", expectedOrder, processFiles)
	}
}
