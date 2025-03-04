package filesetprocessor

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/visitor"
	"github.com/goplus/llcppg/llcppg"
)

type DocFileSetProcessor struct {
	visitedFile map[string]struct{}
	processing  map[string]struct{}
	exec        Exec   // execute a single file
	done        func() // done callback
}

type Exec func(*llcppg.FileEntry) error

type ProcesserConfig struct {
	Exec Exec
	Done func()
}

func NewDocFileSetProcessor(cfg *ProcesserConfig) *DocFileSetProcessor {
	p := &DocFileSetProcessor{
		processing:  make(map[string]struct{}),
		visitedFile: make(map[string]struct{}),
		done:        cfg.Done,
	}
	if cfg.Exec != nil {
		p.exec = cfg.Exec
	}
	return p
}

func (p *DocFileSetProcessor) visitFile(path string, files []*llcppg.FileEntry) {
	if _, ok := p.visitedFile[path]; ok {
		return
	}
	if _, ok := p.processing[path]; ok {
		return
	}
	p.processing[path] = struct{}{}
	idx := FindEntry(files, path)
	if idx < 0 {
		return
	}
	findFile := files[idx]
	for _, include := range findFile.Doc.Includes {
		p.visitFile(include.Path, files)
	}
	if p.exec != nil {
		err := p.exec(findFile)
		if err != nil {
			log.Panic("visit file error: ", err, " file: ", findFile.Path)
		}
	}
	p.visitedFile[findFile.Path] = struct{}{}
	delete(p.processing, findFile.Path)
}

func (p *DocFileSetProcessor) ProcessFileSet(files []*llcppg.FileEntry) error {
	for _, file := range files {
		p.visitFile(file.Path, files)
	}
	if p.done != nil {
		p.done()
	}
	return nil
}

func (p *DocFileSetProcessor) ProcessFileSetFromByte(data []byte) error {
	fileSet, err := config.GetCppgSigfetchFromByte(data)
	if err != nil {
		return err
	}
	return p.ProcessFileSet(fileSet)
}

func (p *DocFileSetProcessor) ProcessFileSetFromPath(filePath string) error {
	data, err := config.ReadFile(filePath)
	if err != nil {
		return err
	}
	return p.ProcessFileSetFromByte(data)
}

// FindEntry finds the entry in FileSet. If useIncPath is true, it searches by IncPath, otherwise by Path
func FindEntry(files []*llcppg.FileEntry, path string) int {
	for i, e := range files {
		if e.Path == path {
			return i
		}
	}
	return -1
}

func ReadSigfetchFile(sigfetchFile string) ([]byte, error) {
	_, file := filepath.Split(sigfetchFile)
	var data []byte
	var err error
	if file == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(sigfetchFile)
	}
	return data, err
}

func New(cfg *convert.Config) (*DocFileSetProcessor, *convert.Package, error) {
	astConvert, err := convert.NewAstConvert(cfg)
	if err != nil {
		return nil, nil, err
	}

	docVisitors := []visitor.DocVisitor{astConvert}
	visitorList := visitor.NewDocVisitorList(docVisitors)

	return NewDocFileSetProcessor(&ProcesserConfig{
		Exec: func(file *llcppg.FileEntry) error {
			visitorList.Visit(file.Doc, file.Path, file.FileType)
			return nil
		},
		Done: func() {
			astConvert.WritePkgFiles()
			astConvert.WriteLinkFile()
			astConvert.WritePubFile()
		},
	}), astConvert.Pkg, nil
}

func Process(cfg *convert.Config) error {
	astConvert, err := convert.NewAstConvert(cfg)
	if err != nil {
		return err
	}

	docVisitors := []visitor.DocVisitor{astConvert}
	visitorList := visitor.NewDocVisitorList(docVisitors)

	p := NewDocFileSetProcessor(&ProcesserConfig{
		Exec: func(file *llcppg.FileEntry) error {
			visitorList.Visit(file.Doc, file.Path, file.FileType)
			return nil
		},
		Done: func() {
			astConvert.WritePkgFiles()
			astConvert.WriteLinkFile()
			astConvert.WritePubFile()
		},
	})

	sigfetchFileData, err := ReadSigfetchFile(cfg.SigfetchFile)
	if err != nil {
		return err
	}

	return p.ProcessFileSetFromByte(sigfetchFileData)
}
