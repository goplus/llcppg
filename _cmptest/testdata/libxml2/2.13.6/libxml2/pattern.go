package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlPattern struct {
	Unused [8]uint8
}
type Pattern X_xmlPattern
type PatternPtr *Pattern
type PatternFlags c.Int

const (
	PATTERN_DEFAULT PatternFlags = 0
	PATTERN_XPATH   PatternFlags = 1
	PATTERN_XSSEL   PatternFlags = 2
	PATTERN_XSFIELD PatternFlags = 4
)

//go:linkname FreePattern C.xmlFreePattern
func FreePattern(comp PatternPtr)

//go:linkname FreePatternList C.xmlFreePatternList
func FreePatternList(comp PatternPtr)

// llgo:link (*Char).Patterncompile C.xmlPatterncompile
func (recv_ *Char) Patterncompile(dict *Dict, flags c.Int, namespaces **Char) PatternPtr {
	return nil
}

// llgo:link (*Char).PatternCompileSafe C.xmlPatternCompileSafe
func (recv_ *Char) PatternCompileSafe(dict *Dict, flags c.Int, namespaces **Char, patternOut *PatternPtr) c.Int {
	return 0
}

//go:linkname PatternMatch C.xmlPatternMatch
func PatternMatch(comp PatternPtr, node NodePtr) c.Int

type X_xmlStreamCtxt struct {
	Unused [8]uint8
}
type StreamCtxt X_xmlStreamCtxt
type StreamCtxtPtr *StreamCtxt

//go:linkname PatternStreamable C.xmlPatternStreamable
func PatternStreamable(comp PatternPtr) c.Int

//go:linkname PatternMaxDepth C.xmlPatternMaxDepth
func PatternMaxDepth(comp PatternPtr) c.Int

//go:linkname PatternMinDepth C.xmlPatternMinDepth
func PatternMinDepth(comp PatternPtr) c.Int

//go:linkname PatternFromRoot C.xmlPatternFromRoot
func PatternFromRoot(comp PatternPtr) c.Int

//go:linkname PatternGetStreamCtxt C.xmlPatternGetStreamCtxt
func PatternGetStreamCtxt(comp PatternPtr) StreamCtxtPtr

//go:linkname FreeStreamCtxt C.xmlFreeStreamCtxt
func FreeStreamCtxt(stream StreamCtxtPtr)

//go:linkname StreamPushNode C.xmlStreamPushNode
func StreamPushNode(stream StreamCtxtPtr, name *Char, ns *Char, nodeType c.Int) c.Int

//go:linkname StreamPush C.xmlStreamPush
func StreamPush(stream StreamCtxtPtr, name *Char, ns *Char) c.Int

//go:linkname StreamPushAttr C.xmlStreamPushAttr
func StreamPushAttr(stream StreamCtxtPtr, name *Char, ns *Char) c.Int

//go:linkname StreamPop C.xmlStreamPop
func StreamPop(stream StreamCtxtPtr) c.Int

//go:linkname StreamWantsAnyNode C.xmlStreamWantsAnyNode
func StreamWantsAnyNode(stream StreamCtxtPtr) c.Int
