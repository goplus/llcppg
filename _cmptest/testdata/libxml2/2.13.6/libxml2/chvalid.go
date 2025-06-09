package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlChSRange struct {
	Low  uint16
	High uint16
}
type ChSRange X_xmlChSRange
type ChSRangePtr *ChSRange

type X_xmlChLRange struct {
	Low  c.Uint
	High c.Uint
}
type ChLRange X_xmlChLRange
type ChLRangePtr *ChLRange

type X_xmlChRangeGroup struct {
	NbShortRange c.Int
	NbLongRange  c.Int
	ShortRange   *ChSRange
	LongRange    *ChLRange
}
type ChRangeGroup X_xmlChRangeGroup
type ChRangeGroupPtr *ChRangeGroup

/**
 * Range checking routine
 */
//go:linkname CharInRange C.xmlCharInRange
func CharInRange(val c.Uint, group *ChRangeGroup) c.Int

//go:linkname IsBaseChar C.xmlIsBaseChar
func IsBaseChar(ch c.Uint) c.Int

//go:linkname IsBlank C.xmlIsBlank
func IsBlank(ch c.Uint) c.Int

//go:linkname IsChar C.xmlIsChar
func IsChar(ch c.Uint) c.Int

//go:linkname IsCombining C.xmlIsCombining
func IsCombining(ch c.Uint) c.Int

//go:linkname IsDigit C.xmlIsDigit
func IsDigit(ch c.Uint) c.Int

//go:linkname IsExtender C.xmlIsExtender
func IsExtender(ch c.Uint) c.Int

//go:linkname IsIdeographic C.xmlIsIdeographic
func IsIdeographic(ch c.Uint) c.Int

//go:linkname IsPubidChar C.xmlIsPubidChar
func IsPubidChar(ch c.Uint) c.Int
