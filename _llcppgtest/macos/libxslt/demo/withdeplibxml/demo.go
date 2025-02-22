package main

import (
	"fmt"
	"os"
	"unsafe"

	"libxslt"

	"github.com/goplus/llgo/c"
	libxml2 "github.com/luoliwoshang/llcppg-libxml"
)

func main() {
	libxml2.XmlInitParser()

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
	xmlDoc := libxml2.XmlReadMemory((*int8)(unsafe.Pointer(unsafe.StringData(xml))), c.Int(len(xml)), nil, nil, 0)
	xsltDoc := libxml2.XmlReadMemory((*int8)(unsafe.Pointer(unsafe.StringData(xslt))), c.Int(len(xslt)), nil, nil, 0)

	if xmlDoc == nil || xsltDoc == nil {
		panic("cant read xml or xslt")
	}

	stylesheet := libxslt.XsltParseStylesheetDoc(xsltDoc)
	if stylesheet == nil {
		panic("cant parse xslt")
	}
	result := libxslt.XsltApplyStylesheet(stylesheet, xmlDoc, (**int8)(unsafe.Pointer(uintptr(0))))
	if result == nil {
		panic("cant apply xslt")
	}

	libxslt.XsltSaveResultToFilename(c.Str("output.html"), result, stylesheet, 0)

	libxml2.XmlFreeDoc(xmlDoc)
	libxml2.XmlFreeDoc(result)
	libxslt.XsltFreeStylesheet(stylesheet)

	libxslt.XsltCleanupGlobals()
	libxml2.XmlCleanupParser()

	buf, err := os.ReadFile("./output.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
