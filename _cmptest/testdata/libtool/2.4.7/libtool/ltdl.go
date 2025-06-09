package libtool

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const LTDL_H = 1

type Handle struct {
	Unused [8]uint8
}
type Dlhandle *Handle

/* Initialisation and finalisation functions for libltdl. */
//go:linkname Dlinit C.lt_dlinit
func Dlinit() c.Int

//go:linkname Dlexit C.lt_dlexit
func Dlexit() c.Int

/* Module search path manipulation.  */
//go:linkname Dladdsearchdir C.lt_dladdsearchdir
func Dladdsearchdir(search_dir *c.Char) c.Int

//go:linkname Dlinsertsearchdir C.lt_dlinsertsearchdir
func Dlinsertsearchdir(before *c.Char, search_dir *c.Char) c.Int

//go:linkname Dlsetsearchpath C.lt_dlsetsearchpath
func Dlsetsearchpath(search_path *c.Char) c.Int

//go:linkname Dlgetsearchpath C.lt_dlgetsearchpath
func Dlgetsearchpath() *c.Char

//go:linkname Dlforeachfile C.lt_dlforeachfile
func Dlforeachfile(search_path *c.Char, func_ func(*c.Char, c.Pointer) c.Int, data c.Pointer) c.Int

/* User module loading advisors.  */
//go:linkname DladviseInit C.lt_dladvise_init
func DladviseInit(advise *Dladvise) c.Int

//go:linkname DladviseDestroy C.lt_dladvise_destroy
func DladviseDestroy(advise *Dladvise) c.Int

//go:linkname DladviseExt C.lt_dladvise_ext
func DladviseExt(advise *Dladvise) c.Int

//go:linkname DladviseResident C.lt_dladvise_resident
func DladviseResident(advise *Dladvise) c.Int

//go:linkname DladviseLocal C.lt_dladvise_local
func DladviseLocal(advise *Dladvise) c.Int

//go:linkname DladviseGlobal C.lt_dladvise_global
func DladviseGlobal(advise *Dladvise) c.Int

//go:linkname DladvisePreload C.lt_dladvise_preload
func DladvisePreload(advise *Dladvise) c.Int

/* Portable libltdl versions of the system dlopen() API. */
//go:linkname Dlopen C.lt_dlopen
func Dlopen(filename *c.Char) Dlhandle

//go:linkname Dlopenext C.lt_dlopenext
func Dlopenext(filename *c.Char) Dlhandle

//go:linkname Dlopenadvise C.lt_dlopenadvise
func Dlopenadvise(filename *c.Char, advise Dladvise) Dlhandle

//go:linkname Dlsym C.lt_dlsym
func Dlsym(handle Dlhandle, name *c.Char) c.Pointer

//go:linkname Dlerror C.lt_dlerror
func Dlerror() *c.Char

//go:linkname Dlclose C.lt_dlclose
func Dlclose(handle Dlhandle) c.Int

/*
A preopened symbol. Arrays of this type comprise the exported

	symbols for a dlpreopened module.
*/
type Dlsymlist struct {
	Name    *c.Char
	Address c.Pointer
}

// llgo:type C
type DlpreloadCallbackFunc func(Dlhandle) c.Int

// llgo:link (*Dlsymlist).Dlpreload C.lt_dlpreload
func (recv_ *Dlsymlist) Dlpreload() c.Int {
	return 0
}

// llgo:link (*Dlsymlist).DlpreloadDefault C.lt_dlpreload_default
func (recv_ *Dlsymlist) DlpreloadDefault() c.Int {
	return 0
}

//go:linkname DlpreloadOpen C.lt_dlpreload_open
func DlpreloadOpen(originator *c.Char, func_ DlpreloadCallbackFunc) c.Int

type DlinterfaceId c.Pointer

// llgo:type C
type DlhandleInterface func(Dlhandle, *c.Char) c.Int

//go:linkname DlinterfaceRegister C.lt_dlinterface_register
func DlinterfaceRegister(id_string *c.Char, iface DlhandleInterface) DlinterfaceId

//go:linkname DlinterfaceFree C.lt_dlinterface_free
func DlinterfaceFree(key DlinterfaceId)

//go:linkname DlcallerSetData C.lt_dlcaller_set_data
func DlcallerSetData(key DlinterfaceId, handle Dlhandle, data c.Pointer) c.Pointer

//go:linkname DlcallerGetData C.lt_dlcaller_get_data
func DlcallerGetData(key DlinterfaceId, handle Dlhandle) c.Pointer

/* Read only information pertaining to a loaded module. */

type Dlinfo struct {
	Filename    *c.Char
	Name        *c.Char
	RefCount    c.Int
	IsResident  c.Uint
	IsSymglobal c.Uint
	IsSymlocal  c.Uint
}

//go:linkname Dlgetinfo C.lt_dlgetinfo
func Dlgetinfo(handle Dlhandle) *Dlinfo

//go:linkname DlhandleIterate C.lt_dlhandle_iterate
func DlhandleIterate(iface DlinterfaceId, place Dlhandle) Dlhandle

//go:linkname DlhandleFetch C.lt_dlhandle_fetch
func DlhandleFetch(iface DlinterfaceId, module_name *c.Char) Dlhandle

//go:linkname DlhandleMap C.lt_dlhandle_map
func DlhandleMap(iface DlinterfaceId, func_ func(Dlhandle, c.Pointer) c.Int, data c.Pointer) c.Int

/* Deprecated module residency management API. */
//go:linkname Dlmakeresident C.lt_dlmakeresident
func Dlmakeresident(handle Dlhandle) c.Int

//go:linkname Dlisresident C.lt_dlisresident
func Dlisresident(handle Dlhandle) c.Int
