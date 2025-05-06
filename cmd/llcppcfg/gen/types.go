package gen

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ObjFile struct {
	OFile string
	HFile string
	Deps  []string
}

func NewObjFile(oFile, hFile string) *ObjFile {
	return &ObjFile{
		OFile: oFile,
		HFile: hFile,
		Deps:  make([]string, 0),
	}
}

func NewObjFileString(str string) *ObjFile {
	fields := strings.Split(str, ":")
	if len(fields) != 2 {
		return nil
	}
	return NewObjFile(fields[0], fields[1])
}

func (p *ObjFile) IsEqual(o *ObjFile) bool {
	if p.HFile != o.HFile {
		return false
	}
	if p.OFile != o.OFile {
		return false
	}
	if len(p.Deps) != len(o.Deps) {
		return false
	}
	for i := range p.Deps {
		if p.Deps[i] != o.Deps[i] {
			return false
		}
	}
	return true
}

func (p *ObjFile) String() string {
	return fmt.Sprintf("{OFile:%s, HFile:%s, Deps:%v}", p.OFile, p.HFile, p.Deps)
}

type IncludeList struct {
	include    []string
	absPathMap map[string]struct{}
	relPathMap map[string]struct{}
}

func NewIncludeList() *IncludeList {
	return &IncludeList{include: make([]string, 0), absPathMap: make(map[string]struct{}), relPathMap: make(map[string]struct{})}
}

func (p *IncludeList) AddCflagEntry(index int, entry *CflagEntry) {
	if entry == nil {
		return
	}
	if entry.IsEmpty() {
		return
	}
	for _, objFile := range entry.ObjFiles {
		absPath := filepath.Join(entry.Include, objFile.HFile)
		_, ok := p.absPathMap[absPath]
		if !ok {
			p.absPathMap[absPath] = struct{}{}
			p.AddIncludeForObjFile(objFile, index)
		}
	}
	lenInvalidObjFiles := len(entry.InvalidObjFiles)
	if lenInvalidObjFiles > 0 {
		fmt.Println("Invlid header files:")
		for idx, objFile := range entry.InvalidObjFiles {
			if idx < lenInvalidObjFiles-1 {
				fmt.Printf("\t\t%q,\n", objFile.HFile)
			} else {
				fmt.Printf("\t\t%q\n", objFile.HFile)
			}
		}
	}
}

func (p *IncludeList) AddIncludeForObjFile(objFile *ObjFile, index int) {
	hFile := objFile.HFile
	_, ok := p.relPathMap[objFile.HFile]
	if ok {
		hFile = fmt.Sprintf("%d:%s", index, objFile.HFile)
	}
	p.relPathMap[objFile.HFile] = struct{}{}
	p.include = append(p.include, hFile)
}

type CflagEntry struct {
	Include         string
	ObjFiles        []*ObjFile
	InvalidObjFiles []*ObjFile
}

func (c *CflagEntry) IsEmpty() bool {
	if len(c.Include) == 0 {
		return true
	}
	if len(c.ObjFiles) == 0 && len(c.InvalidObjFiles) == 0 {
		return true
	}
	return false
}

func (c *CflagEntry) String() string {
	return fmt.Sprintf("{Include:%s, ObjFiles:%v}", c.Include, c.ObjFiles)
}
