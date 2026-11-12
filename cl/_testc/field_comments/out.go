package foo

import (
	"github.com/goplus/lib/c"
	"unsafe"
)

const XGoPackage = true

type Record struct {
	// The horizontal coordinate.
	X     c.Int // In pixels.
	Y     c.Int // The vertical coordinate.
	Plain c.Int
	// A documented value.
	// Kept on two lines.
	Multiline c.Int
	Values    [2]c.Int // Two stored values.
	Total     c.Int    // A field renamed by configuration.
	// A block starting on the field's line.
	// Its continuation stays with this field.
	InlineBlock c.Int
	AfterBlock  c.Int
	// The field's leading documentation.
	// A second description.
	// Kept with the leading documentation.
	BothMultiline c.Int
	AfterBoth     c.Int
}
type Trailing struct {
	Value c.Int
}
type Positioned struct {
	First c.Int
	// Only the second field's documentation.
	// Its second documentation line.
	Second c.Int // The second field's inline comment.
	Third  c.Int
}
type Nested struct {
	Child          _llcppg_anon_0    // The named record.
	Entries        [2]_llcppg_anon_1 // The record array.
	_llcppg_anon_2                   // The embedded struct.
	_llcppg_anon_3                   // The embedded union.
	Tail           c.Int             // After the embedded records.
}
type Value struct {
	_xgo_union [3]uint32
}
type _llcppg_anon_0 struct {
	Value c.Int // A named record's member.
}
type _llcppg_anon_1 struct {
	Value c.Int // An array element's member.
}

// The embedded struct.
type _llcppg_anon_2 struct {
	Promoted c.Int // A promoted struct member.
}

// The embedded union.
type _llcppg_anon_3 struct {
	_xgo_union [1]uint32
}
type _llcppg_anon_4 struct {
	// The first member's leading documentation.
	First  c.Int // The first promoted member.
	Second c.Int // The second promoted member.
	_llcppg_anon_5
}
type _llcppg_anon_5 struct {
	_xgo_union [1]uint32
}

// A promoted union member.
func (p *_llcppg_anon_3) XGof_ref_code() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}
func (p *_llcppg_anon_3) XGof_ref_value() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}

// A member promoted through both records.
func (p *_llcppg_anon_5) XGof_ref_nested() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// An integer value.
// In units of one.
func (p *Value) XGof_ref_number() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// A floating-point value.
func (p *Value) XGof_ref_real() *c.Float {
	return (*c.Float)(unsafe.Pointer(p))
}
func (p *Value) XGof_ref_plain() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

type _llcppg_anon_6 struct {
	Tag c.Int // The named view's tag.
}

// The named struct view.
func (p *Value) XGof_ref_named() *_llcppg_anon_6 {
	return (*_llcppg_anon_6)(unsafe.Pointer(p))
}

// A member promoted through an anonymous union.
func (p *Value) XGof_ref_other() *c.Int {
	return (*c.Int)(unsafe.Pointer(p))
}

// The first member's leading documentation.
// The first promoted member.
func (p *Value) XGof_ref_first() *c.Int {
	return &(*_llcppg_anon_4)(unsafe.Pointer(p)).First
}

// The second promoted member.
func (p *Value) XGof_ref_second() *c.Int {
	return &(*_llcppg_anon_4)(unsafe.Pointer(p)).Second
}

// A member promoted through both records.
func (p *Value) XGof_ref_nested() *c.Int {
	return (*_llcppg_anon_4)(unsafe.Pointer(p)).XGof_ref_nested()
}
