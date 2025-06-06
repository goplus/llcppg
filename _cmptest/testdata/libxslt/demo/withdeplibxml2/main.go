package main

import (
	"fmt"
	"os"
	"unsafe"

	"libxslt"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
)

func main() {
	libxml2.InitParser()

	xml :=
		`<?xml version='1.0'?>
		<root>
			<person>
				<name>Alice</name>
				<age>25</age>
			</person>
		</root>`
	xslt := `<?xml version="1.0" encoding="UTF-8"?>
		<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
			<xsl:template match="/">
				<html>
					<body>
						<h1>个人信息</h1>
						<p>姓名: <xsl:value-of select="root/person/name"/></p>
						<p>年龄: <xsl:value-of select="root/person/age"/></p>
					</body>
				</html>
			</xsl:template>
		</xsl:stylesheet>
	`
	xmlDoc := libxml2.ReadMemory((*int8)(unsafe.Pointer(unsafe.StringData(xml))), c.Int(len(xml)), nil, nil, 0)
	xsltDoc := libxml2.ReadMemory((*int8)(unsafe.Pointer(unsafe.StringData(xslt))), c.Int(len(xslt)), nil, nil, 0)

	if xmlDoc == nil || xsltDoc == nil {
		panic("cant read xml or xslt")
	}

	stylesheet := libxslt.ParseStylesheetDoc(xsltDoc)
	if stylesheet == nil {
		panic("cant parse xslt")
	}
	result := libxslt.ApplyStylesheet(stylesheet, xmlDoc, (**int8)(unsafe.Pointer(uintptr(0))))
	if result == nil {
		panic("cant apply xslt")
	}

	libxslt.SaveResultToFilename(c.Str("output.html"), result, stylesheet, 0)

	libxml2.FreeDoc(xmlDoc)
	libxml2.FreeDoc(result)
	libxslt.FreeStylesheet(stylesheet)

	libxslt.CleanupGlobals()
	libxml2.CleanupParser()

	buf, err := os.ReadFile("./output.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
