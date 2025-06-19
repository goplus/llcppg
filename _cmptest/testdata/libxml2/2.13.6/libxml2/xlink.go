package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type XlinkHRef *Char
type XlinkRole *Char
type XlinkTitle *Char
type XlinkType c.Int

const (
	XLINK_TYPE_NONE         XlinkType = 0
	XLINK_TYPE_SIMPLE       XlinkType = 1
	XLINK_TYPE_EXTENDED     XlinkType = 2
	XLINK_TYPE_EXTENDED_SET XlinkType = 3
)

type XlinkShow c.Int

const (
	XLINK_SHOW_NONE    XlinkShow = 0
	XLINK_SHOW_NEW     XlinkShow = 1
	XLINK_SHOW_EMBED   XlinkShow = 2
	XLINK_SHOW_REPLACE XlinkShow = 3
)

type XlinkActuate c.Int

const (
	XLINK_ACTUATE_NONE      XlinkActuate = 0
	XLINK_ACTUATE_AUTO      XlinkActuate = 1
	XLINK_ACTUATE_ONREQUEST XlinkActuate = 2
)

// llgo:type C
type XlinkNodeDetectFunc func(c.Pointer, NodePtr)

// llgo:type C
type XlinkSimpleLinkFunk func(c.Pointer, NodePtr, XlinkHRef, XlinkRole, XlinkTitle)

// llgo:type C
type XlinkExtendedLinkFunk func(c.Pointer, NodePtr, c.Int, *XlinkHRef, *XlinkRole, c.Int, *XlinkRole, *XlinkRole, *XlinkShow, *XlinkActuate, c.Int, *XlinkTitle, **Char)

// llgo:type C
type XlinkExtendedLinkSetFunk func(c.Pointer, NodePtr, c.Int, *XlinkHRef, *XlinkRole, c.Int, *XlinkTitle, **Char)

type X_xlinkHandler struct {
	Simple   XlinkSimpleLinkFunk
	Extended XlinkExtendedLinkFunk
	Set      XlinkExtendedLinkSetFunk
}
type XlinkHandler X_xlinkHandler
type XlinkHandlerPtr *XlinkHandler

/*
 * The default detection routine, can be overridden, they call the default
 * detection callbacks.
 */
//go:linkname XlinkGetDefaultDetect C.xlinkGetDefaultDetect
func XlinkGetDefaultDetect() XlinkNodeDetectFunc

//go:linkname XlinkSetDefaultDetect C.xlinkSetDefaultDetect
func XlinkSetDefaultDetect(func_ XlinkNodeDetectFunc)

/*
 * Routines to set/get the default handlers.
 */
//go:linkname XlinkGetDefaultHandler C.xlinkGetDefaultHandler
func XlinkGetDefaultHandler() XlinkHandlerPtr

//go:linkname XlinkSetDefaultHandler C.xlinkSetDefaultHandler
func XlinkSetDefaultHandler(handler XlinkHandlerPtr)

/*
 * Link detection module itself.
 */
//go:linkname XlinkIsLink C.xlinkIsLink
func XlinkIsLink(doc DocPtr, node NodePtr) XlinkType
