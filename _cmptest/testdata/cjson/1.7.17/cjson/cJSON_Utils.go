package cjson

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

/* Implement RFC6901 (https://tools.ietf.org/html/rfc6901) JSON Pointer spec. */
// llgo:link (*JSON).GetPointer C.cJSONUtils_GetPointer
func (recv_ *JSON) GetPointer(pointer *c.Char) *JSON {
	return nil
}

// llgo:link (*JSON).GetPointerCaseSensitive C.cJSONUtils_GetPointerCaseSensitive
func (recv_ *JSON) GetPointerCaseSensitive(pointer *c.Char) *JSON {
	return nil
}

/* Implement RFC6902 (https://tools.ietf.org/html/rfc6902) JSON Patch spec. */
/* NOTE: This modifies objects in 'from' and 'to' by sorting the elements by their key */
// llgo:link (*JSON).GeneratePatches C.cJSONUtils_GeneratePatches
func (recv_ *JSON) GeneratePatches(to *JSON) *JSON {
	return nil
}

// llgo:link (*JSON).GeneratePatchesCaseSensitive C.cJSONUtils_GeneratePatchesCaseSensitive
func (recv_ *JSON) GeneratePatchesCaseSensitive(to *JSON) *JSON {
	return nil
}

/* Utility for generating patch array entries. */
// llgo:link (*JSON).AddPatchToArray C.cJSONUtils_AddPatchToArray
func (recv_ *JSON) AddPatchToArray(operation *c.Char, path *c.Char, value *JSON) {
}

/* Returns 0 for success. */
// llgo:link (*JSON).ApplyPatches C.cJSONUtils_ApplyPatches
func (recv_ *JSON) ApplyPatches(patches *JSON) c.Int {
	return 0
}

// llgo:link (*JSON).ApplyPatchesCaseSensitive C.cJSONUtils_ApplyPatchesCaseSensitive
func (recv_ *JSON) ApplyPatchesCaseSensitive(patches *JSON) c.Int {
	return 0
}

/* Implement RFC7386 (https://tools.ietf.org/html/rfc7396) JSON Merge Patch spec. */
/* target will be modified by patch. return value is new ptr for target. */
// llgo:link (*JSON).MergePatch C.cJSONUtils_MergePatch
func (recv_ *JSON) MergePatch(patch *JSON) *JSON {
	return nil
}

// llgo:link (*JSON).MergePatchCaseSensitive C.cJSONUtils_MergePatchCaseSensitive
func (recv_ *JSON) MergePatchCaseSensitive(patch *JSON) *JSON {
	return nil
}

/* generates a patch to move from -> to */
/* NOTE: This modifies objects in 'from' and 'to' by sorting the elements by their key */
// llgo:link (*JSON).GenerateMergePatch C.cJSONUtils_GenerateMergePatch
func (recv_ *JSON) GenerateMergePatch(to *JSON) *JSON {
	return nil
}

// llgo:link (*JSON).GenerateMergePatchCaseSensitive C.cJSONUtils_GenerateMergePatchCaseSensitive
func (recv_ *JSON) GenerateMergePatchCaseSensitive(to *JSON) *JSON {
	return nil
}

/* Given a root object and a target object, construct a pointer from one to the other. */
// llgo:link (*JSON).FindPointerFromObjectTo C.cJSONUtils_FindPointerFromObjectTo
func (recv_ *JSON) FindPointerFromObjectTo(target *JSON) *c.Char {
	return nil
}

/* Sorts the members of the object into alphabetical order. */
// llgo:link (*JSON).SortObject C.cJSONUtils_SortObject
func (recv_ *JSON) SortObject() {
}

// llgo:link (*JSON).SortObjectCaseSensitive C.cJSONUtils_SortObjectCaseSensitive
func (recv_ *JSON) SortObjectCaseSensitive() {
}
