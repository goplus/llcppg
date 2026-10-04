package foo

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -L/path/foo -lfoo"

type Iterator[_Category any, _Tp any, _Distance any, _Pointer any, _Reference any] struct {
}
type X_IteratorAlias[_Category any, _Tp any, _Distance any, _Pointer any, _Reference any] = Iterator[_Category, _Tp, _Distance, _Pointer, _Reference]
type X_IteratorBase[_Derived any, _Category any, _Tp any, _Distance any, _Pointer any, _Reference any] = Iterator[_Category, _Tp, _Distance, _Pointer, _Reference]
type OutputIteratorTag = c.Int
type BackInsertIterator[_Container any] struct {
	Iterator
	container *_Container
}
