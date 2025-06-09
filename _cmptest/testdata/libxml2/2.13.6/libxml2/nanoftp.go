package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

// llgo:type C
type FtpListCallback func(c.Pointer, *c.Char, *c.Char, *c.Char, *c.Char, c.Ulong, c.Int, c.Int, *c.Char, c.Int, c.Int, c.Int)

// llgo:type C
type FtpDataCallback func(c.Pointer, *c.Char, c.Int)

/*
 * Init
 */
//go:linkname NanoFTPInit C.xmlNanoFTPInit
func NanoFTPInit()

//go:linkname NanoFTPCleanup C.xmlNanoFTPCleanup
func NanoFTPCleanup()

/*
 * Creating/freeing contexts.
 */
//go:linkname NanoFTPNewCtxt C.xmlNanoFTPNewCtxt
func NanoFTPNewCtxt(URL *c.Char) c.Pointer

//go:linkname NanoFTPFreeCtxt C.xmlNanoFTPFreeCtxt
func NanoFTPFreeCtxt(ctx c.Pointer)

//go:linkname NanoFTPConnectTo C.xmlNanoFTPConnectTo
func NanoFTPConnectTo(server *c.Char, port c.Int) c.Pointer

/*
 * Opening/closing session connections.
 */
//go:linkname NanoFTPOpen C.xmlNanoFTPOpen
func NanoFTPOpen(URL *c.Char) c.Pointer

//go:linkname NanoFTPConnect C.xmlNanoFTPConnect
func NanoFTPConnect(ctx c.Pointer) c.Int

//go:linkname NanoFTPClose C.xmlNanoFTPClose
func NanoFTPClose(ctx c.Pointer) c.Int

//go:linkname NanoFTPQuit C.xmlNanoFTPQuit
func NanoFTPQuit(ctx c.Pointer) c.Int

//go:linkname NanoFTPScanProxy C.xmlNanoFTPScanProxy
func NanoFTPScanProxy(URL *c.Char)

//go:linkname NanoFTPProxy C.xmlNanoFTPProxy
func NanoFTPProxy(host *c.Char, port c.Int, user *c.Char, passwd *c.Char, type_ c.Int)

//go:linkname NanoFTPUpdateURL C.xmlNanoFTPUpdateURL
func NanoFTPUpdateURL(ctx c.Pointer, URL *c.Char) c.Int

/*
 * Rather internal commands.
 */
//go:linkname NanoFTPGetResponse C.xmlNanoFTPGetResponse
func NanoFTPGetResponse(ctx c.Pointer) c.Int

//go:linkname NanoFTPCheckResponse C.xmlNanoFTPCheckResponse
func NanoFTPCheckResponse(ctx c.Pointer) c.Int

/*
 * CD/DIR/GET handlers.
 */
//go:linkname NanoFTPCwd C.xmlNanoFTPCwd
func NanoFTPCwd(ctx c.Pointer, directory *c.Char) c.Int

//go:linkname NanoFTPDele C.xmlNanoFTPDele
func NanoFTPDele(ctx c.Pointer, file *c.Char) c.Int

//go:linkname NanoFTPGetConnection C.xmlNanoFTPGetConnection
func NanoFTPGetConnection(ctx c.Pointer) c.Int

//go:linkname NanoFTPCloseConnection C.xmlNanoFTPCloseConnection
func NanoFTPCloseConnection(ctx c.Pointer) c.Int

//go:linkname NanoFTPList C.xmlNanoFTPList
func NanoFTPList(ctx c.Pointer, callback FtpListCallback, userData c.Pointer, filename *c.Char) c.Int

//go:linkname NanoFTPGetSocket C.xmlNanoFTPGetSocket
func NanoFTPGetSocket(ctx c.Pointer, filename *c.Char) c.Int

//go:linkname NanoFTPGet C.xmlNanoFTPGet
func NanoFTPGet(ctx c.Pointer, callback FtpDataCallback, userData c.Pointer, filename *c.Char) c.Int

//go:linkname NanoFTPRead C.xmlNanoFTPRead
func NanoFTPRead(ctx c.Pointer, dest c.Pointer, len c.Int) c.Int
