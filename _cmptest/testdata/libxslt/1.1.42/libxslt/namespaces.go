package libxslt

import (
	"github.com/goplus/llcppg/_cmptest/testdata/libxml2/2.13.6/libxml2"
	_ "unsafe"
)

//go:linkname NamespaceAlias C.xsltNamespaceAlias
func NamespaceAlias(style StylesheetPtr, node libxml2.NodePtr)

//go:linkname GetNamespace C.xsltGetNamespace
func GetNamespace(ctxt TransformContextPtr, cur libxml2.NodePtr, ns libxml2.NsPtr, out libxml2.NodePtr) libxml2.NsPtr

//go:linkname GetPlainNamespace C.xsltGetPlainNamespace
func GetPlainNamespace(ctxt TransformContextPtr, cur libxml2.NodePtr, ns libxml2.NsPtr, out libxml2.NodePtr) libxml2.NsPtr

//go:linkname GetSpecialNamespace C.xsltGetSpecialNamespace
func GetSpecialNamespace(ctxt TransformContextPtr, cur libxml2.NodePtr, URI *libxml2.Char, prefix *libxml2.Char, out libxml2.NodePtr) libxml2.NsPtr

//go:linkname CopyNamespace C.xsltCopyNamespace
func CopyNamespace(ctxt TransformContextPtr, elem libxml2.NodePtr, ns libxml2.NsPtr) libxml2.NsPtr

//go:linkname CopyNamespaceList C.xsltCopyNamespaceList
func CopyNamespaceList(ctxt TransformContextPtr, node libxml2.NodePtr, cur libxml2.NsPtr) libxml2.NsPtr

//go:linkname FreeNamespaceAliasHashes C.xsltFreeNamespaceAliasHashes
func FreeNamespaceAliasHashes(style StylesheetPtr)
