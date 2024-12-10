package processor

import (
	"log"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	"github.com/goplus/llcppg/cmd/gogensig/visitor"
)

type DocVisitorManager struct {
	VisitorList []visitor.DocVisitor
}

func NewDocVisitorManager(visitorList []visitor.DocVisitor) *DocVisitorManager {
	return &DocVisitorManager{VisitorList: visitorList}
}

func (p *DocVisitorManager) Visit(node ast.Node, path string, incPath string, isSys bool) bool {
	for _, v := range p.VisitorList {
		v.VisitStart(path, incPath, isSys)
		v.Visit(node)
		v.VisitDone(path)
	}
	return true
}

type DocFileSetProcessor struct {
	visitedFile map[string]struct{}
	processing  map[string]struct{}
	exec        Exec     // execute a single file
	done        func()   // done callback
	depIncs     []string // abs path
}

type Exec func(*unmarshal.FileEntry) error

type ProcesserConfig struct {
	Exec    Exec
	Done    func()
	DepIncs []string // abs path
}

// allDepIncs is the absolute path of all dependent include files
// such as /path/to/foo.h, etc. skip these files,because they are already processed
func NewDocFileSetProcessor(cfg *ProcesserConfig) *DocFileSetProcessor {
	p := &DocFileSetProcessor{
		processing:  make(map[string]struct{}),
		visitedFile: make(map[string]struct{}),
		done:        cfg.Done,
		depIncs:     cfg.DepIncs,
	}
	if cfg.Exec != nil {
		p.exec = cfg.Exec
	}
	return p
}

func (p *DocFileSetProcessor) visitFile(path string, files unmarshal.FileSet) {
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
		err := p.exec(&findFile)
		if err != nil {
			log.Panic("visit file error: ", err, " file: ", findFile.Path)
		}
	}
	p.visitedFile[findFile.Path] = struct{}{}
	delete(p.processing, findFile.Path)
}

func (p *DocFileSetProcessor) ProcessFileSet(files unmarshal.FileSet) error {
	for _, inc := range p.depIncs {
		idx := FindEntry(files, inc)
		if idx < 0 {
			continue
		}
		p.visitedFile[files[idx].Path] = struct{}{}
	}
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
func FindEntry(files unmarshal.FileSet, path string) int {
	for i, e := range files {
		if e.Path == path {
			return i
		}
	}
	return -1
}
