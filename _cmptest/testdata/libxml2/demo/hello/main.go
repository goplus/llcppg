package main

import (
	"unsafe"

	"libxml2"

	"github.com/goplus/lib/c"
)

func main() {
	libxml2.InitParser()
	xml := "<?xml version='1.0'?><root><person><name>Alice</name><age>25</age></person></root>"
	doc := libxml2.ReadMemory((*int8)(unsafe.Pointer(unsafe.StringData(xml))), c.Int(len(xml)), nil, nil, 0)
	if doc == nil {
		panic("Failed to parse XML")
	}
	docPtr := (*libxml2.Doc)(unsafe.Pointer(doc))
	root := docPtr.DocGetRootElement()
	c.Printf(c.Str("Root element: %s\n"), root.Name)
	libxml2.FreeDoc(doc)
	libxml2.CleanupParser()
}
