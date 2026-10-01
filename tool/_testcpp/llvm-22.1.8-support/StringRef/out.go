package StringRef

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

const XGoPackage = true
const LLGoPackage = "link: -L$(llvm-config --libdir) -lLLVM; -lLLVM"

type APInt struct {
}
type HashCode struct {
}

// StringRef - Represent a constant reference to a string, i.e. a character
// array and a length, which need not be null terminated.
//
// This class does not own the string data, it is expected to be used in
// situations where the character data resides in some other buffer, whose
// lifetime extends past that of the StringRef. For this reason, it is not in
// general safe to store a StringRef.
type StringRef struct {
	Data   *c.Char
	Length c.SizeT
}
type StringRefIterator = *c.Char
type StringRefConstIterator = *c.Char
type StringRefSizeType = c.SizeT
type StringRefValueType = c.Char

// A wrapper around a string literal that serves as a proxy for constructing
// global tables of StringRefs with the length computed at compile time.
// In order to avoid the invocation of a global constructor, StringLiteral
// should *only* be used in a constexpr context, as such:
//
// constexpr StringLiteral S("test");
type StringLiteral struct {
	StringRef
}

// Helper functions for StringRef::getAsInteger.
//
//go:linkname GetAsUnsignedInteger C._ZN4llvm20getAsUnsignedIntegerENS_9StringRefEjRy
func GetAsUnsignedInteger(Str StringRef, Radix c.Uint, Result *c.UlongLong) bool

//go:linkname GetAsSignedInteger C._ZN4llvm18getAsSignedIntegerENS_9StringRefEjRx
func GetAsSignedInteger(Str StringRef, Radix c.Uint, Result *c.LongLong) bool

//go:linkname GetAutoSenseRadix C._ZN4llvm17getAutoSenseRadixERNS_9StringRefE
func GetAutoSenseRadix(Str *StringRef) c.Uint

//go:linkname ConsumeUnsignedInteger C._ZN4llvm22consumeUnsignedIntegerERNS_9StringRefEjRy
func ConsumeUnsignedInteger(Str *StringRef, Radix c.Uint, Result *c.UlongLong) bool

//go:linkname ConsumeSignedInteger C._ZN4llvm20consumeSignedIntegerERNS_9StringRefEjRx
func ConsumeSignedInteger(Str *StringRef, Radix c.Uint, Result *c.LongLong) bool

//go:linkname StringRefNpos C._ZN4llvm9StringRef4nposE
var StringRefNpos c.SizeT

// Compare two strings, ignoring case.
//
// llgo:link (*StringRef).CompareInsensitive C._ZNK4llvm9StringRef19compare_insensitiveES0_
func (this *StringRef) CompareInsensitive(RHS StringRef) c.Int {
	return 0
}

// compare_numeric - Compare two strings, treating sequences of digits as
// numbers.
//
// llgo:link (*StringRef).CompareNumeric C._ZNK4llvm9StringRef15compare_numericES0_
func (this *StringRef) CompareNumeric(RHS StringRef) c.Int {
	return 0
}

// Determine the edit distance between this string and another
// string.
//
// \param Other the string to compare this string against.
//
// \param AllowReplacements whether to allow character
// replacements (change one character into another) as a single
// operation, rather than as two operations (an insertion and a
// removal).
//
// \param MaxEditDistance If non-zero, the maximum edit distance that
// this routine is allowed to compute. If the edit distance will exceed
// that maximum, returns \c MaxEditDistance+1.
//
// \returns the minimum number of character insertions, removals,
// or (if \p AllowReplacements is \c true) replacements needed to
// transform one of the given strings into the other. If zero,
// the strings are identical.
//
// llgo:link (*StringRef).EditDistance C._ZNK4llvm9StringRef13edit_distanceES0_bj
func (this *StringRef) EditDistance(Other StringRef, AllowReplacements bool, MaxEditDistance c.Uint) c.Uint {
	return 0
}

// llgo:link (*StringRef).EditDistanceInsensitive C._ZNK4llvm9StringRef25edit_distance_insensitiveES0_bj
func (this *StringRef) EditDistanceInsensitive(Other StringRef, AllowReplacements bool, MaxEditDistance c.Uint) c.Uint {
	return 0
}

// Check if this string starts with the given \p Prefix, ignoring case.
//
// llgo:link (*StringRef).StartsWithInsensitive C._ZNK4llvm9StringRef23starts_with_insensitiveES0_
func (this *StringRef) StartsWithInsensitive(Prefix StringRef) bool {
	return false
}

// Check if this string ends with the given \p Suffix, ignoring case.
//
// llgo:link (*StringRef).EndsWithInsensitive C._ZNK4llvm9StringRef21ends_with_insensitiveES0_
func (this *StringRef) EndsWithInsensitive(Suffix StringRef) bool {
	return false
}

// Search for the first character \p C in the string, ignoring case.
//
// \returns The index of the first occurrence of \p C, or npos if not
// found.
//
// llgo:link (*StringRef).FindInsensitive__0 C._ZNK4llvm9StringRef16find_insensitiveEcm
func (this *StringRef) FindInsensitive__0(C c.Char, From c.SizeT) c.SizeT {
	return 0
}

// Search for the first string \p Str in the string.
//
// \returns The index of the first occurrence of \p Str, or npos if not
// found.
//
// llgo:link (*StringRef).Find__1 C._ZNK4llvm9StringRef4findES0_m
func (this *StringRef) Find__1(Str StringRef, From c.SizeT) c.SizeT {
	return 0
}

// Search for the first string \p Str in the string, ignoring case.
//
// \returns The index of the first occurrence of \p Str, or npos if not
// found.
//
// llgo:link (*StringRef).FindInsensitive__1 C._ZNK4llvm9StringRef16find_insensitiveES0_m
func (this *StringRef) FindInsensitive__1(Str StringRef, From c.SizeT) c.SizeT {
	return 0
}

// Search for the last character \p C in the string, ignoring case.
//
// \returns The index of the last occurrence of \p C, or npos if not
// found.
//
// llgo:link (*StringRef).RfindInsensitive__1 C._ZNK4llvm9StringRef17rfind_insensitiveEcm
func (this *StringRef) RfindInsensitive__1(C c.Char, From c.SizeT) c.SizeT {
	return 0
}

// Search for the last string \p Str in the string.
//
// \returns The index of the last occurrence of \p Str, or npos if not
// found.
//
// llgo:link (*StringRef).Rfind__0 C._ZNK4llvm9StringRef5rfindES0_
func (this *StringRef) Rfind__0(Str StringRef) c.SizeT {
	return 0
}

// Search for the last string \p Str in the string, ignoring case.
//
// \returns The index of the last occurrence of \p Str, or npos if not
// found.
//
// llgo:link (*StringRef).RfindInsensitive__0 C._ZNK4llvm9StringRef17rfind_insensitiveES0_
func (this *StringRef) RfindInsensitive__0(Str StringRef) c.SizeT {
	return 0
}

// Find the first character in the string that is in \p Chars, or npos if
// not found.
//
// Complexity: O(size() + Chars.size())
//
// llgo:link (*StringRef).FindFirstOf__1 C._ZNK4llvm9StringRef13find_first_ofES0_m
func (this *StringRef) FindFirstOf__1(Chars StringRef, From c.SizeT) c.SizeT {
	return 0
}

// Find the first character in the string that is not \p C or npos if not
// found.
//
// llgo:link (*StringRef).FindFirstNotOf__0 C._ZNK4llvm9StringRef17find_first_not_ofEcm
func (this *StringRef) FindFirstNotOf__0(C c.Char, From c.SizeT) c.SizeT {
	return 0
}

// Find the first character in the string that is not in the string
// \p Chars, or npos if not found.
//
// Complexity: O(size() + Chars.size())
//
// llgo:link (*StringRef).FindFirstNotOf__1 C._ZNK4llvm9StringRef17find_first_not_ofES0_m
func (this *StringRef) FindFirstNotOf__1(Chars StringRef, From c.SizeT) c.SizeT {
	return 0
}

// Find the last character in the string that is in \p C, or npos if not
// found.
//
// Complexity: O(size() + Chars.size())
//
// llgo:link (*StringRef).FindLastOf__1 C._ZNK4llvm9StringRef12find_last_ofES0_m
func (this *StringRef) FindLastOf__1(Chars StringRef, From c.SizeT) c.SizeT {
	return 0
}

// Find the last character in the string that is not \p C, or npos if not
// found.
//
// llgo:link (*StringRef).FindLastNotOf__0 C._ZNK4llvm9StringRef16find_last_not_ofEcm
func (this *StringRef) FindLastNotOf__0(C c.Char, From c.SizeT) c.SizeT {
	return 0
}

// Find the last character in the string that is not in \p Chars, or
// npos if not found.
//
// Complexity: O(size() + Chars.size())
//
// llgo:link (*StringRef).FindLastNotOf__1 C._ZNK4llvm9StringRef16find_last_not_ofES0_m
func (this *StringRef) FindLastNotOf__1(Chars StringRef, From c.SizeT) c.SizeT {
	return 0
}

// Return the number of non-overlapped occurrences of \p Str in
// the string.
//
// llgo:link (*StringRef).Count__1 C._ZNK4llvm9StringRef5countES0_
func (this *StringRef) Count__1(Str StringRef) c.SizeT {
	return 0
}

// Parse the current string as an integer of the specified \p Radix, or of
// an autosensed radix if the \p Radix given is 0.  The current value in
// \p Result is discarded, and the storage is changed to be wide enough to
// store the parsed integer.
//
// \returns true if the string does not solely consist of a valid
// non-empty number in the appropriate base.
//
// APInt::fromString is superficially similar but assumes the
// string is well-formed in the given radix.
//
// llgo:link (*StringRef).GetAsInteger__1 C._ZNK4llvm9StringRef12getAsIntegerEjRNS_5APIntE
func (this *StringRef) GetAsInteger__1(Radix c.Uint, Result *APInt) bool {
	return false
}

// Parse the current string as an integer of the specified \p Radix.  If
// \p Radix is specified as zero, this does radix autosensing using
// extended C rules: 0 is octal, 0x is hex, 0b is binary.
//
// If the string does not begin with a number of the specified radix,
// this returns true to signify the error. The string is considered
// erroneous if empty.
// The portion of the string representing the discovered numeric value
// is removed from the beginning of the string.
//
// llgo:link (*StringRef).ConsumeInteger__1 C._ZN4llvm9StringRef14consumeIntegerEjRNS_5APIntE
func (this *StringRef) ConsumeInteger__1(Radix c.Uint, Result *APInt) bool {
	return false
}

// Parse the current string as an IEEE double-precision floating
// point value.  The string must be a well-formed double.
//
// If \p AllowInexact is false, the function will fail if the string
// cannot be represented exactly.  Otherwise, the function only fails
// in case of an overflow or underflow, or an invalid floating point
// representation.
//
// llgo:link (*StringRef).GetAsDouble C._ZNK4llvm9StringRef11getAsDoubleERdb
func (this *StringRef) GetAsDouble(Result *c.Double, AllowInexact bool) bool {
	return false
}

// Compute a hash_code for a StringRef.
//
//go:linkname HashValue C._ZN4llvm10hash_valueENS_9StringRefE
func HashValue(S StringRef) HashCode
