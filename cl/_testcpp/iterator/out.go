package foo

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L/path/foo -lfoo"

type Iterator[_Category any, _Tp any, _Distance any, _Pointer any, _Reference any] struct {
}
type X__IteratorAlias[_Category any, _Tp any, _Distance any, _Pointer any, _Reference any] = Iterator[_Category, _Tp, _Distance, _Pointer, _Reference]
type X__IteratorBase[_Derived any, _Category any, _Tp any, _Distance any, _Pointer any, _Reference any] = Iterator[_Category, _Tp, _Distance, _Pointer, _Reference]
type OutputIteratorTag = c.Int
type BackInsertIterator[_Container any] struct {
	X__IteratorBase[BackInsertIterator[_Container], OutputIteratorTag, c.Void, c.Void, c.Void, c.Void]
	container *_Container
}
type IteratorTraits[_Ip any] struct {
}
type ReverseIterator[_Iter any] struct {
	current _Iter
}
