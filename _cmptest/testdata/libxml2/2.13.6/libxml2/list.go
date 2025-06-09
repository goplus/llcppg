package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type X_xmlLink struct {
	Unused [8]uint8
}
type Link X_xmlLink
type LinkPtr *Link

type X_xmlList struct {
	Unused [8]uint8
}
type List X_xmlList
type ListPtr *List

// llgo:type C
type ListDeallocator func(LinkPtr)

// llgo:type C
type ListDataCompare func(c.Pointer, c.Pointer) c.Int

// llgo:type C
type ListWalker func(c.Pointer, c.Pointer) c.Int

/* Creation/Deletion */
//go:linkname ListCreate C.xmlListCreate
func ListCreate(deallocator ListDeallocator, compare ListDataCompare) ListPtr

//go:linkname ListDelete C.xmlListDelete
func ListDelete(l ListPtr)

/* Basic Operators */
//go:linkname ListSearch C.xmlListSearch
func ListSearch(l ListPtr, data c.Pointer) c.Pointer

//go:linkname ListReverseSearch C.xmlListReverseSearch
func ListReverseSearch(l ListPtr, data c.Pointer) c.Pointer

//go:linkname ListInsert C.xmlListInsert
func ListInsert(l ListPtr, data c.Pointer) c.Int

//go:linkname ListAppend C.xmlListAppend
func ListAppend(l ListPtr, data c.Pointer) c.Int

//go:linkname ListRemoveFirst C.xmlListRemoveFirst
func ListRemoveFirst(l ListPtr, data c.Pointer) c.Int

//go:linkname ListRemoveLast C.xmlListRemoveLast
func ListRemoveLast(l ListPtr, data c.Pointer) c.Int

//go:linkname ListRemoveAll C.xmlListRemoveAll
func ListRemoveAll(l ListPtr, data c.Pointer) c.Int

//go:linkname ListClear C.xmlListClear
func ListClear(l ListPtr)

//go:linkname ListEmpty C.xmlListEmpty
func ListEmpty(l ListPtr) c.Int

//go:linkname ListFront C.xmlListFront
func ListFront(l ListPtr) LinkPtr

//go:linkname ListEnd C.xmlListEnd
func ListEnd(l ListPtr) LinkPtr

//go:linkname ListSize C.xmlListSize
func ListSize(l ListPtr) c.Int

//go:linkname ListPopFront C.xmlListPopFront
func ListPopFront(l ListPtr)

//go:linkname ListPopBack C.xmlListPopBack
func ListPopBack(l ListPtr)

//go:linkname ListPushFront C.xmlListPushFront
func ListPushFront(l ListPtr, data c.Pointer) c.Int

//go:linkname ListPushBack C.xmlListPushBack
func ListPushBack(l ListPtr, data c.Pointer) c.Int

/* Advanced Operators */
//go:linkname ListReverse C.xmlListReverse
func ListReverse(l ListPtr)

//go:linkname ListSort C.xmlListSort
func ListSort(l ListPtr)

//go:linkname ListWalk C.xmlListWalk
func ListWalk(l ListPtr, walker ListWalker, user c.Pointer)

//go:linkname ListReverseWalk C.xmlListReverseWalk
func ListReverseWalk(l ListPtr, walker ListWalker, user c.Pointer)

//go:linkname ListMerge C.xmlListMerge
func ListMerge(l1 ListPtr, l2 ListPtr)

//go:linkname ListDup C.xmlListDup
func ListDup(old ListPtr) ListPtr

//go:linkname ListCopy C.xmlListCopy
func ListCopy(cur ListPtr, old ListPtr) c.Int

/* Link operators */
//go:linkname LinkGetData C.xmlLinkGetData
func LinkGetData(lk LinkPtr) c.Pointer
