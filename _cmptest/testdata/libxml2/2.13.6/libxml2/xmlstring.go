package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type Char c.Char

/*
 * xmlChar handling
 */
// llgo:link (*Char).Strdup C.xmlStrdup
func (recv_ *Char) Strdup() *Char {
	return nil
}

// llgo:link (*Char).Strndup C.xmlStrndup
func (recv_ *Char) Strndup(len c.Int) *Char {
	return nil
}

//go:linkname CharStrndup C.xmlCharStrndup
func CharStrndup(cur *c.Char, len c.Int) *Char

//go:linkname CharStrdup C.xmlCharStrdup
func CharStrdup(cur *c.Char) *Char

// llgo:link (*Char).Strsub C.xmlStrsub
func (recv_ *Char) Strsub(start c.Int, len c.Int) *Char {
	return nil
}

// llgo:link (*Char).Strchr C.xmlStrchr
func (recv_ *Char) Strchr(val Char) *Char {
	return nil
}

// llgo:link (*Char).Strstr C.xmlStrstr
func (recv_ *Char) Strstr(val *Char) *Char {
	return nil
}

// llgo:link (*Char).Strcasestr C.xmlStrcasestr
func (recv_ *Char) Strcasestr(val *Char) *Char {
	return nil
}

// llgo:link (*Char).Strcmp C.xmlStrcmp
func (recv_ *Char) Strcmp(str2 *Char) c.Int {
	return 0
}

// llgo:link (*Char).Strncmp C.xmlStrncmp
func (recv_ *Char) Strncmp(str2 *Char, len c.Int) c.Int {
	return 0
}

// llgo:link (*Char).Strcasecmp C.xmlStrcasecmp
func (recv_ *Char) Strcasecmp(str2 *Char) c.Int {
	return 0
}

// llgo:link (*Char).Strncasecmp C.xmlStrncasecmp
func (recv_ *Char) Strncasecmp(str2 *Char, len c.Int) c.Int {
	return 0
}

// llgo:link (*Char).StrEqual C.xmlStrEqual
func (recv_ *Char) StrEqual(str2 *Char) c.Int {
	return 0
}

// llgo:link (*Char).StrQEqual C.xmlStrQEqual
func (recv_ *Char) StrQEqual(name *Char, str *Char) c.Int {
	return 0
}

// llgo:link (*Char).Strlen C.xmlStrlen
func (recv_ *Char) Strlen() c.Int {
	return 0
}

// llgo:link (*Char).Strcat C.xmlStrcat
func (recv_ *Char) Strcat(add *Char) *Char {
	return nil
}

// llgo:link (*Char).Strncat C.xmlStrncat
func (recv_ *Char) Strncat(add *Char, len c.Int) *Char {
	return nil
}

// llgo:link (*Char).StrncatNew C.xmlStrncatNew
func (recv_ *Char) StrncatNew(str2 *Char, len c.Int) *Char {
	return nil
}

//go:linkname StrPrintf C.xmlStrPrintf
func StrPrintf(buf *Char, len c.Int, msg *c.Char, __llgo_va_list ...interface{}) c.Int

// llgo:link (*Char).StrVPrintf C.xmlStrVPrintf
func (recv_ *Char) StrVPrintf(len c.Int, msg *c.Char, ap c.VaList) c.Int {
	return 0
}

//go:linkname GetUTF8Char C.xmlGetUTF8Char
func GetUTF8Char(utf *c.Char, len *c.Int) c.Int

//go:linkname CheckUTF8 C.xmlCheckUTF8
func CheckUTF8(utf *c.Char) c.Int

// llgo:link (*Char).UTF8Strsize C.xmlUTF8Strsize
func (recv_ *Char) UTF8Strsize(len c.Int) c.Int {
	return 0
}

// llgo:link (*Char).UTF8Strndup C.xmlUTF8Strndup
func (recv_ *Char) UTF8Strndup(len c.Int) *Char {
	return nil
}

// llgo:link (*Char).UTF8Strpos C.xmlUTF8Strpos
func (recv_ *Char) UTF8Strpos(pos c.Int) *Char {
	return nil
}

// llgo:link (*Char).UTF8Strloc C.xmlUTF8Strloc
func (recv_ *Char) UTF8Strloc(utfchar *Char) c.Int {
	return 0
}

// llgo:link (*Char).UTF8Strsub C.xmlUTF8Strsub
func (recv_ *Char) UTF8Strsub(start c.Int, len c.Int) *Char {
	return nil
}

// llgo:link (*Char).UTF8Strlen C.xmlUTF8Strlen
func (recv_ *Char) UTF8Strlen() c.Int {
	return 0
}

// llgo:link (*Char).UTF8Size C.xmlUTF8Size
func (recv_ *Char) UTF8Size() c.Int {
	return 0
}

// llgo:link (*Char).UTF8Charcmp C.xmlUTF8Charcmp
func (recv_ *Char) UTF8Charcmp(utf2 *Char) c.Int {
	return 0
}
