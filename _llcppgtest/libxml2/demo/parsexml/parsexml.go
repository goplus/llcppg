package main

import (
	"libxml2"
	"unsafe"

	"github.com/goplus/lib/c"
)

func main() {
	libxml2.XmlInitParser()
	xml := "<?xml version='1.0'?><root><person><name>Alice</name><age>25</age></person></root>"
	doc := libxml2.XmlReadMemory((*int8)(unsafe.Pointer(unsafe.StringData(xml))), c.Int(len(xml)), nil, nil, 0)
	if doc == nil {
		panic("Failed to parse XML")
	}
	docPtr := (*libxml2.XmlDoc)(unsafe.Pointer(doc))
	root := docPtr.XmlDocGetRootElement()
	c.Printf(c.Str("Root element: %s\n"), root.Name)
	libxml2.XmlFreeDoc(doc)
	libxml2.XmlCleanupParser()
}
