package libxslt

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname NewLocale C.xsltNewLocale
func NewLocale(langName *libxml2.Char, lowerFirst c.Int) c.Pointer

//go:linkname FreeLocale C.xsltFreeLocale
func FreeLocale(locale c.Pointer)

//go:linkname Strxfrm C.xsltStrxfrm
func Strxfrm(locale c.Pointer, string *libxml2.Char) *libxml2.Char

//go:linkname FreeLocales C.xsltFreeLocales
func FreeLocales()

type Locale c.Pointer
type LocaleChar libxml2.Char

//go:linkname LocaleStrcmp C.xsltLocaleStrcmp
func LocaleStrcmp(locale c.Pointer, str1 *libxml2.Char, str2 *libxml2.Char) c.Int
