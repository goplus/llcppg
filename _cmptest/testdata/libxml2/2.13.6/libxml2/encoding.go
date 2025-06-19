package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type CharEncError c.Int

const (
	ENC_ERR_SUCCESS  CharEncError = 0
	ENC_ERR_SPACE    CharEncError = -1
	ENC_ERR_INPUT    CharEncError = -2
	ENC_ERR_PARTIAL  CharEncError = -3
	ENC_ERR_INTERNAL CharEncError = -4
	ENC_ERR_MEMORY   CharEncError = -5
)

type CharEncoding c.Int

const (
	CHAR_ENCODING_ERROR     CharEncoding = -1
	CHAR_ENCODING_NONE      CharEncoding = 0
	CHAR_ENCODING_UTF8      CharEncoding = 1
	CHAR_ENCODING_UTF16LE   CharEncoding = 2
	CHAR_ENCODING_UTF16BE   CharEncoding = 3
	CHAR_ENCODING_UCS4LE    CharEncoding = 4
	CHAR_ENCODING_UCS4BE    CharEncoding = 5
	CHAR_ENCODING_EBCDIC    CharEncoding = 6
	CHAR_ENCODING_UCS4_2143 CharEncoding = 7
	CHAR_ENCODING_UCS4_3412 CharEncoding = 8
	CHAR_ENCODING_UCS2      CharEncoding = 9
	CHAR_ENCODING_8859_1    CharEncoding = 10
	CHAR_ENCODING_8859_2    CharEncoding = 11
	CHAR_ENCODING_8859_3    CharEncoding = 12
	CHAR_ENCODING_8859_4    CharEncoding = 13
	CHAR_ENCODING_8859_5    CharEncoding = 14
	CHAR_ENCODING_8859_6    CharEncoding = 15
	CHAR_ENCODING_8859_7    CharEncoding = 16
	CHAR_ENCODING_8859_8    CharEncoding = 17
	CHAR_ENCODING_8859_9    CharEncoding = 18
	CHAR_ENCODING_2022_JP   CharEncoding = 19
	CHAR_ENCODING_SHIFT_JIS CharEncoding = 20
	CHAR_ENCODING_EUC_JP    CharEncoding = 21
	CHAR_ENCODING_ASCII     CharEncoding = 22
)

// llgo:type C
type CharEncodingInputFunc func(*c.Char, *c.Int, *c.Char, *c.Int) c.Int

// llgo:type C
type CharEncodingOutputFunc func(*c.Char, *c.Int, *c.Char, *c.Int) c.Int

type X_xmlCharEncodingHandler struct {
	Name   *c.Char
	Input  CharEncodingInputFunc
	Output CharEncodingOutputFunc
}
type CharEncodingHandler X_xmlCharEncodingHandler
type CharEncodingHandlerPtr *CharEncodingHandler

/*
 * Interfaces for encoding handlers.
 */
//go:linkname InitCharEncodingHandlers C.xmlInitCharEncodingHandlers
func InitCharEncodingHandlers()

//go:linkname CleanupCharEncodingHandlers C.xmlCleanupCharEncodingHandlers
func CleanupCharEncodingHandlers()

//go:linkname RegisterCharEncodingHandler C.xmlRegisterCharEncodingHandler
func RegisterCharEncodingHandler(handler CharEncodingHandlerPtr)

// llgo:link CharEncoding.LookupCharEncodingHandler C.xmlLookupCharEncodingHandler
func (recv_ CharEncoding) LookupCharEncodingHandler(out *CharEncodingHandlerPtr) c.Int {
	return 0
}

//go:linkname OpenCharEncodingHandler C.xmlOpenCharEncodingHandler
func OpenCharEncodingHandler(name *c.Char, output c.Int, out *CharEncodingHandlerPtr) c.Int

// llgo:link CharEncoding.GetCharEncodingHandler C.xmlGetCharEncodingHandler
func (recv_ CharEncoding) GetCharEncodingHandler() CharEncodingHandlerPtr {
	return nil
}

//go:linkname FindCharEncodingHandler C.xmlFindCharEncodingHandler
func FindCharEncodingHandler(name *c.Char) CharEncodingHandlerPtr

//go:linkname NewCharEncodingHandler C.xmlNewCharEncodingHandler
func NewCharEncodingHandler(name *c.Char, input CharEncodingInputFunc, output CharEncodingOutputFunc) CharEncodingHandlerPtr

/*
 * Interfaces for encoding names and aliases.
 */
//go:linkname AddEncodingAlias C.xmlAddEncodingAlias
func AddEncodingAlias(name *c.Char, alias *c.Char) c.Int

//go:linkname DelEncodingAlias C.xmlDelEncodingAlias
func DelEncodingAlias(alias *c.Char) c.Int

//go:linkname GetEncodingAlias C.xmlGetEncodingAlias
func GetEncodingAlias(alias *c.Char) *c.Char

//go:linkname CleanupEncodingAliases C.xmlCleanupEncodingAliases
func CleanupEncodingAliases()

//go:linkname ParseCharEncoding C.xmlParseCharEncoding
func ParseCharEncoding(name *c.Char) CharEncoding

// llgo:link CharEncoding.GetCharEncodingName C.xmlGetCharEncodingName
func (recv_ CharEncoding) GetCharEncodingName() *c.Char {
	return nil
}

/*
 * Interfaces directly used by the parsers.
 */
//go:linkname DetectCharEncoding C.xmlDetectCharEncoding
func DetectCharEncoding(in *c.Char, len c.Int) CharEncoding

/** DOC_ENABLE */
// llgo:link (*CharEncodingHandler).CharEncOutFunc C.xmlCharEncOutFunc
func (recv_ *CharEncodingHandler) CharEncOutFunc(out *X_xmlBuffer, in *X_xmlBuffer) c.Int {
	return 0
}

// llgo:link (*CharEncodingHandler).CharEncInFunc C.xmlCharEncInFunc
func (recv_ *CharEncodingHandler) CharEncInFunc(out *X_xmlBuffer, in *X_xmlBuffer) c.Int {
	return 0
}

// llgo:link (*CharEncodingHandler).CharEncFirstLine C.xmlCharEncFirstLine
func (recv_ *CharEncodingHandler) CharEncFirstLine(out *X_xmlBuffer, in *X_xmlBuffer) c.Int {
	return 0
}

// llgo:link (*CharEncodingHandler).CharEncCloseFunc C.xmlCharEncCloseFunc
func (recv_ *CharEncodingHandler) CharEncCloseFunc() c.Int {
	return 0
}

/*
 * Export a few useful functions
 */
//go:linkname UTF8Toisolat1 C.UTF8Toisolat1
func UTF8Toisolat1(out *c.Char, outlen *c.Int, in *c.Char, inlen *c.Int) c.Int

//go:linkname Isolat1ToUTF8 C.isolat1ToUTF8
func Isolat1ToUTF8(out *c.Char, outlen *c.Int, in *c.Char, inlen *c.Int) c.Int
