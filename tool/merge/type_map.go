// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in tool/merge/internal/astutil/LICENSE.

package merge

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"hash/maphash"

	"github.com/goplus/llcppg/tool/merge/internal/astutil"
)

// astMap is a hash table that maps structurally equal AST nodes to arbitrary values.
// For example, separately parsed copies of `type T int` share one entry.
//
// Just as with map[K]V, a nil *astMap is a valid empty map.
//
// Read-only map operations ([astMap.at], [astMap.len], and so on) may
// safely be called concurrently.
type astMap struct {
	table  map[uint32][]entry // maps hash to bucket; entry.key==nil means unused
	length int                // number of map entries
	hasher astHasher
}

// entry is an entry (key/value association) in a hash bucket.
type entry struct {
	key   ast.Node
	value any
}

// setHasher has no effect.
//
// It is a relic of an optimization that is no longer profitable. Do
// not use [astHasher], [makeHasher], or [setHasher] in new code.
func (m *astMap) setHasher(astHasher) {}

// delete removes the entry with the given key, if any.
// It returns true if the entry was found.
func (m *astMap) delete(key ast.Node) bool {
	if m != nil && m.table != nil {
		hash := m.hash(key)
		bucket := m.table[hash]
		for i, e := range bucket {
			if e.key != nil && m.equal(key, e.key) {
				// We can't compact the bucket as it
				// would disturb iterators.
				bucket[i] = entry{}
				m.length--
				return true
			}
		}
	}
	return false
}

// at returns the map entry for the given key.
// The result is nil if the entry is not present.
func (m *astMap) at(key ast.Node) any {
	if m != nil && m.table != nil {
		for _, e := range m.table[m.hash(key)] {
			if e.key != nil && m.equal(key, e.key) {
				return e.value
			}
		}
	}
	return nil
}

// set sets the map entry for key to val,
// and returns the previous entry, if any.
func (m *astMap) set(key ast.Node, value any) (prev any) {
	if m.table != nil {
		hash := m.hash(key)
		bucket := m.table[hash]
		var hole *entry
		for i, e := range bucket {
			if e.key == nil {
				hole = &bucket[i]
			} else if m.equal(key, e.key) {
				prev = e.value
				bucket[i].value = value
				return
			}
		}

		if hole != nil {
			*hole = entry{key, value} // overwrite deleted entry
		} else {
			m.table[hash] = append(bucket, entry{key, value})
		}
	} else {
		hash := m.hash(key)
		m.table = map[uint32][]entry{hash: {entry{key, value}}}
	}

	m.length++
	return
}

// len returns the number of map entries.
func (m *astMap) len() int {
	if m != nil {
		return m.length
	}
	return 0
}

// iterate calls function f on each entry in the map in unspecified order.
//
// If f should mutate the map, iterate provides the same guarantees as
// Go maps: if f deletes a map entry that iterate has not yet reached,
// f will not be invoked for it, but if f inserts a map entry that
// iterate has not yet reached, whether or not f will be invoked for
// it is unspecified.
func (m *astMap) iterate(f func(key ast.Node, value any)) {
	if m != nil {
		for _, bucket := range m.table {
			for _, e := range bucket {
				if e.key != nil {
					f(e.key, e.value)
				}
			}
		}
	}
}

// keys returns a new slice containing the set of map keys.
// The order is unspecified.
func (m *astMap) keys() []ast.Node {
	keys := make([]ast.Node, 0, m.len())
	m.iterate(func(key ast.Node, _ any) {
		keys = append(keys, key)
	})
	return keys
}

func (m *astMap) toString(values bool) string {
	if m == nil {
		return "{}"
	}
	var buf bytes.Buffer
	fset := token.NewFileSet()
	fmt.Fprint(&buf, "{")
	sep := ""
	m.iterate(func(key ast.Node, value any) {
		fmt.Fprint(&buf, sep)
		sep = ", "
		start := buf.Len()
		if err := format.Node(&buf, fset, key); err != nil {
			// Preserve the node representation for keys such as a standalone FieldList.
			buf.Truncate(start)
			fmt.Fprint(&buf, key)
		}
		if values {
			fmt.Fprintf(&buf, ": %q", value)
		}
	})
	fmt.Fprint(&buf, "}")
	return buf.String()
}

// string returns a string representation of the map's entries.
// Values are printed using fmt.Sprintf("%v", v).
// Order is unspecified.
func (m *astMap) string() string {
	return m.toString(true)
}

// keysString returns a string representation of the map's key set.
// Order is unspecified.
func (m *astMap) keysString() string {
	return m.toString(false)
}

// -- astHasher --

// Equal AST keys must hash to the same bucket, even for different import aliases.
func (m *astMap) hash(t ast.Node) uint32 {
	return m.hasher.hash(t)
}

// An astHasher maps AST nodes to hash values.
type astHasher struct {
	identName func(*ast.Ident) (string, bool)
}

var theHasher astHasher

// makeHasher returns astHasher{}.
func makeHasher() astHasher { return theHasher }

func (h astHasher) hash(t ast.Node) uint32 {
	if t == nil {
		panic("cannot hash a nil AST key")
	}
	return h.hashNode(t)
}

// hashString computes the Fowler–Noll–Vo hash of s.
func hashString(s string) uint32 {
	var h uint32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// Use the AST expressions directly so array lengths such as [N]int need no type checker.
func (h astHasher) hashNode(t ast.Node) uint32 {
	switch t := t.(type) {
	case *ast.Ident:
		return h.hashTypeName(t)

	case *ast.BasicLit:
		return uint32(t.Kind) + 3*hashString(t.Value)

	case *ast.TypeSpec:
		hash := 2*h.hashTypeName(t.Name) + 3*h.hashNode(t.Type) + 7*h.hashTypeParams(t.TypeParams)
		if t.Assign.IsValid() {
			hash++
		}
		return hash

	case *ast.ArrayType:
		if t.Len == nil {
			return 9049 + 2*h.hashNode(t.Elt)
		}
		return 9043 + 2*h.hashNode(t.Len) + 3*h.hashNode(t.Elt)

	case *ast.StructType:
		var hash uint32 = 9059
		for _, f := range t.Fields.List {
			if len(f.Names) == 0 {
				hash += 8861
			}
			if f.Tag != nil {
				hash += hashString(f.Tag.Value)
			}
			for _, name := range f.Names {
				hash += h.hashTypeName(name)
			}
			hash += h.hashNode(f.Type)
		}
		return hash

	case *ast.StarExpr:
		return 9067 + 2*h.hashNode(t.X)

	case *ast.FuncType:
		var hash uint32 = 9091
		if isASTVariadic(t.Params) {
			hash *= 8863
		}
		hash += 7 * h.hashTypeParams(t.TypeParams)
		return hash + 3*h.hashTuple(t.Params) + 5*h.hashTuple(t.Results)

	case *ast.BinaryExpr:
		if t.Op == token.OR {
			return h.hashUnion(t)
		}
		return uint32(t.Op) + 3*h.hashNode(t.X) + 5*h.hashNode(t.Y)

	case *ast.InterfaceType:
		var hash uint32 = 9103
		for _, m := range t.Methods.List {
			if len(m.Names) == 0 {
				hash += h.hashTermSet([]ast.Expr{m.Type})
				continue
			}
			for _, name := range m.Names {
				hash += 3*h.hashTypeName(name) + 5*h.shallowHash(m.Type)
			}
		}
		return hash

	case *ast.MapType:
		return 9109 + 2*h.hashNode(t.Key) + 3*h.hashNode(t.Value)

	case *ast.ChanType:
		return 9127 + 2*uint32(t.Dir) + 3*h.hashNode(t.Value)

	case *ast.SelectorExpr:
		return 2*h.hashNode(t.X) + 3*h.hashTypeName(t.Sel)

	case *ast.IndexExpr:
		return h.hashNode(t.X) + 2*h.hashNode(t.Index)

	case *ast.IndexListExpr:
		hash := h.hashNode(t.X)
		for _, targ := range t.Indices {
			hash += 2 * h.hashNode(targ)
		}
		return hash

	case *ast.FieldList:
		return h.hashTuple(t)

	case *ast.ParenExpr:
		return h.hashNode(t.X)

	case *ast.UnaryExpr:
		return uint32(t.Op) + 3*h.hashNode(t.X)

	case *ast.Ellipsis:
		return 9049 + 2*h.hashNode(t.Elt)
	}

	// Other expressions, such as len(s) in [len(s)]int, can share a bucket.
	// Equal compares their full AST when the bucket is searched.
	return hashString(fmt.Sprintf("%T", t))
}

func (h astHasher) hashTuple(tuple *ast.FieldList) uint32 {
	n := tuple.NumFields()
	hash := 9137 + 2*uint32(n)
	if tuple != nil {
		for _, f := range tuple.List {
			hash += 3 * uint32(max(1, len(f.Names))) * h.hashNode(f.Type)
		}
	}
	return hash
}

func (h astHasher) hashUnion(t *ast.BinaryExpr) uint32 {
	return h.hashTermSet([]ast.Expr{t.X, t.Y})
}

func (h astHasher) hashTermSet(terms []ast.Expr) uint32 {
	hash := 9157 + 2*uint32(len(terms))
	for _, term := range terms {
		termHash := h.hashNode(term)
		if term, ok := term.(*ast.UnaryExpr); ok && term.Op == token.TILDE {
			termHash *= 9161
		}
		hash += 3 * termHash
	}
	return hash
}

// Type parameters such as T in func[T ~int]() carry their constraints in a FieldList.
func (h astHasher) hashTypeParams(params *ast.FieldList) uint32 {
	var hash uint32
	if params != nil {
		for _, p := range params.List {
			hash += uint32(len(p.Names)) * h.hashNode(p.Type)
		}
	}
	return hash
}

var theSeed = maphash.MakeSeed()

// Separately parsed identifiers share a hash when their names or import paths agree.
func (h astHasher) hashTypeName(tname *ast.Ident) uint32 {
	name, _ := h.nameOf(tname)
	hash := maphash.String(theSeed, name)
	return uint32(hash ^ (hash >> 32))
}

// Hash an interface method's outer type shapes, e.g. the slice in func() []T.
func (h astHasher) shallowHash(t ast.Node) uint32 {
	switch t := t.(type) {
	case *ast.FuncType:
		var hash uint32 = 604171
		if isASTVariadic(t.Params) {
			hash *= 971767
		}
		return hash + 1062599*h.shallowHash(t.Params) + 1282529*h.shallowHash(t.Results)

	case *ast.FieldList:
		n := t.NumFields()
		hash := 9137 + 2*uint32(n)
		if t != nil {
			for _, f := range t.List {
				hash += 53471161 * uint32(max(1, len(f.Names))) * h.shallowHash(f.Type)
			}
		}
		return hash

	case *ast.ArrayType:
		if t.Len == nil {
			return 2690201
		}
		return 1524181 + 2*h.hashNode(t.Len)

	case *ast.StructType:
		return 3326489

	case *ast.StarExpr:
		return 4393139

	case *ast.BinaryExpr:
		if t.Op == token.OR {
			return 562448657
		}

	case *ast.InterfaceType:
		return 2124679

	case *ast.MapType:
		return 9109

	case *ast.ChanType:
		return 9127

	case *ast.Ident:
		return h.hashTypeName(t)

	case *ast.IndexExpr:
		return h.shallowHash(t.X)

	case *ast.IndexListExpr:
		return h.shallowHash(t.X)

	case *ast.ParenExpr:
		return h.shallowHash(t.X)

	case *ast.Ellipsis:
		return 2690201
	}
	return h.hashNode(t)
}

func isASTVariadic(params *ast.FieldList) bool {
	if params == nil || len(params.List) == 0 {
		return false
	}
	_, ok := params.List[len(params.List)-1].Type.(*ast.Ellipsis)
	return ok
}

type importPath string

func equalASTKey(x, y ast.Node) bool {
	return (&astMap{}).equal(x, y)
}

func (m *astMap) equal(x, y ast.Node) bool {
	return astutil.Equal(x, y, func(a, b *ast.Ident) bool {
		aName, aImported := m.hasher.nameOf(a)
		bName, bImported := m.hasher.nameOf(b)
		return aImported == bImported && aName == bName
	})
}

func (h astHasher) nameOf(ident *ast.Ident) (string, bool) {
	if h.identName != nil {
		return h.identName(ident)
	}
	return astIdentName(ident)
}

// A qualified reference uses its import path so c.Int and libc.Int can compare equal.
func astIdentName(ident *ast.Ident) (name string, imported bool) {
	if ident.Obj != nil && ident.Obj.Kind == ast.Pkg {
		if path, ok := ident.Obj.Data.(importPath); ok {
			return string(path), true
		}
	}
	return ident.Name, false
}
