package libxml2

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

type EntityType c.Int

const (
	INTERNAL_GENERAL_ENTITY          EntityType = 1
	EXTERNAL_GENERAL_PARSED_ENTITY   EntityType = 2
	EXTERNAL_GENERAL_UNPARSED_ENTITY EntityType = 3
	INTERNAL_PARAMETER_ENTITY        EntityType = 4
	EXTERNAL_PARAMETER_ENTITY        EntityType = 5
	INTERNAL_PREDEFINED_ENTITY       EntityType = 6
)

type EntitiesTable X_xmlHashTable
type EntitiesTablePtr *EntitiesTable

/*
 * External functions:
 */
//go:linkname InitializePredefinedEntities C.xmlInitializePredefinedEntities
func InitializePredefinedEntities()

//go:linkname NewEntity C.xmlNewEntity
func NewEntity(doc DocPtr, name *Char, type_ c.Int, ExternalID *Char, SystemID *Char, content *Char) EntityPtr

//go:linkname FreeEntity C.xmlFreeEntity
func FreeEntity(entity EntityPtr)

//go:linkname AddEntity C.xmlAddEntity
func AddEntity(doc DocPtr, extSubset c.Int, name *Char, type_ c.Int, ExternalID *Char, SystemID *Char, content *Char, out *EntityPtr) c.Int

//go:linkname AddDocEntity C.xmlAddDocEntity
func AddDocEntity(doc DocPtr, name *Char, type_ c.Int, ExternalID *Char, SystemID *Char, content *Char) EntityPtr

//go:linkname AddDtdEntity C.xmlAddDtdEntity
func AddDtdEntity(doc DocPtr, name *Char, type_ c.Int, ExternalID *Char, SystemID *Char, content *Char) EntityPtr

// llgo:link (*Char).GetPredefinedEntity C.xmlGetPredefinedEntity
func (recv_ *Char) GetPredefinedEntity() EntityPtr {
	return nil
}

// llgo:link (*Doc).GetDocEntity C.xmlGetDocEntity
func (recv_ *Doc) GetDocEntity(name *Char) EntityPtr {
	return nil
}

//go:linkname GetDtdEntity C.xmlGetDtdEntity
func GetDtdEntity(doc DocPtr, name *Char) EntityPtr

//go:linkname GetParameterEntity C.xmlGetParameterEntity
func GetParameterEntity(doc DocPtr, name *Char) EntityPtr

//go:linkname EncodeEntities C.xmlEncodeEntities
func EncodeEntities(doc DocPtr, input *Char) *Char

//go:linkname EncodeEntitiesReentrant C.xmlEncodeEntitiesReentrant
func EncodeEntitiesReentrant(doc DocPtr, input *Char) *Char

// llgo:link (*Doc).EncodeSpecialChars C.xmlEncodeSpecialChars
func (recv_ *Doc) EncodeSpecialChars(input *Char) *Char {
	return nil
}

//go:linkname CreateEntitiesTable C.xmlCreateEntitiesTable
func CreateEntitiesTable() EntitiesTablePtr

//go:linkname CopyEntitiesTable C.xmlCopyEntitiesTable
func CopyEntitiesTable(table EntitiesTablePtr) EntitiesTablePtr

//go:linkname FreeEntitiesTable C.xmlFreeEntitiesTable
func FreeEntitiesTable(table EntitiesTablePtr)

//go:linkname DumpEntitiesTable C.xmlDumpEntitiesTable
func DumpEntitiesTable(buf BufferPtr, table EntitiesTablePtr)

//go:linkname DumpEntityDecl C.xmlDumpEntityDecl
func DumpEntityDecl(buf BufferPtr, ent EntityPtr)

//go:linkname CleanupPredefinedEntities C.xmlCleanupPredefinedEntities
func CleanupPredefinedEntities()
