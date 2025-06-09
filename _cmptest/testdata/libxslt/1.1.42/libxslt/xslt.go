package libxslt

import _ "unsafe"

const DEFAULT_VERSION = "1.0"
const DEFAULT_VENDOR = "libxslt"
const DEFAULT_URL = "http://xmlsoft.org/XSLT/"

/*
 * Global initialization function.
 */
//go:linkname Init C.xsltInit
func Init()

/*
 * Global cleanup function.
 */
//go:linkname CleanupGlobals C.xsltCleanupGlobals
func CleanupGlobals()
