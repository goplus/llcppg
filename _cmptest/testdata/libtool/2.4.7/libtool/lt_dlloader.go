package libtool

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const DLLOADER_H = 1

type Dlloader c.Pointer
type Module c.Pointer
type UserData c.Pointer

type Advise struct {
	Unused [8]uint8
}
type Dladvise *Advise

// llgo:type C
type ModuleOpen func(UserData, *c.Char, Dladvise) Module

// llgo:type C
type ModuleClose func(UserData, Module) c.Int

// llgo:type C
type FindSym func(UserData, Module, *c.Char) c.Pointer

// llgo:type C
type DlloaderInit func(UserData) c.Int

// llgo:type C
type DlloaderExit func(UserData) c.Int
type DlloaderPriority c.Int

const (
	DLLOADER_PREPEND DlloaderPriority = 0
	DLLOADER_APPEND  DlloaderPriority = 1
)

/*
This structure defines a module loader, as populated by the get_vtable

	entry point of each loader.
*/
type Dlvtable struct {
	Name         *c.Char
	SymPrefix    *c.Char
	ModuleOpen   *ModuleOpen
	ModuleClose  *ModuleClose
	FindSym      *FindSym
	DlloaderInit *DlloaderInit
	DlloaderExit *DlloaderExit
	DlloaderData UserData
	Priority     DlloaderPriority
}

// llgo:link (*Dlvtable).DlloaderAdd C.lt_dlloader_add
func (recv_ *Dlvtable) DlloaderAdd() c.Int {
	return 0
}

//go:linkname DlloaderNext C.lt_dlloader_next
func DlloaderNext(loader Dlloader) Dlloader

//go:linkname DlloaderRemove C.lt_dlloader_remove
func DlloaderRemove(name *c.Char) *Dlvtable

//go:linkname DlloaderFind C.lt_dlloader_find
func DlloaderFind(name *c.Char) *Dlvtable

//go:linkname DlloaderGet C.lt_dlloader_get
func DlloaderGet(loader Dlloader) *Dlvtable

// llgo:type C
type GetVtable func(UserData) *Dlvtable
