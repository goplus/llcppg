package unmarshal

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/xgo-dev/llcppg/ast"
)

type NodeUnmarshaler func(data []byte) (ast.Node, error)

var nodeUnmarshalers map[string]NodeUnmarshaler

func init() {
	nodeUnmarshalers = map[string]NodeUnmarshaler{
		// Not need costum unmarshal
		"Token":       Token,
		"Macro":       Macro,
		"Include":     Include,
		"BasicLit":    BasicLit,
		"BuiltinType": BuiltinType,
		"Ident":       Ident,
		"Variadic":    Variadic,

		"PointerType":      PointerType,
		"LvalueRefType":    LvalueRefType,
		"RvalueRefType":    RvalueRefType,
		"BlockPointerType": BlockPointerType,

		"ArrayType":   ArrayType,
		"Field":       Field,
		"FieldList":   FieldList,
		"ScopingExpr": ScopingExpr,
		"TagExpr":     TagExpr,
		"EnumItem":    EnumItem,
		"EnumType":    EnumType,
		"FuncType":    FuncType,
		"RecordType":  RecordType,
		"TypedefDecl": TypeDefDecl,

		"FuncDecl":     FuncDecl,
		"TypeDecl":     TypeDecl,
		"EnumTypeDecl": EnumTypeDecl,
	}
}

func File(data []byte) (*ast.File, error) {
	type fileTemp struct {
		Decls    []json.RawMessage `json:"decls"`
		Includes []*ast.Include    `json:"includes,omitempty"`
		Macros   []*ast.Macro      `json:"macros,omitempty"`
	}
	var fileData fileTemp
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, newDeserializeError("File", fileData, data, err)
	}

	var decls []ast.Decl

	for _, declData := range fileData.Decls {
		declNode, err := Node(declData)
		if err != nil {
			return nil, newUnmarshalFieldError("File", fileData, "Decls", declData, err)
		}
		decl, ok := declNode.(ast.Decl)
		if !ok {
			return nil, newUnexpectTypeError("File", declNode, "ast.Decl")
		}
		decls = append(decls, decl)
	}

	return &ast.File{
		Includes: fileData.Includes,
		Macros:   fileData.Macros,
		Decls:    decls,
	}, nil
}

func Node(data []byte) (ast.Node, error) {
	type nodeTemp struct {
		Type string `json:"_Type"`
	}
	var nodeData nodeTemp
	if err := json.Unmarshal(data, &nodeData); err != nil {
		return nil, newDeserializeError("Node", nodeData, data, err)
	}

	unmarshaler, ok := nodeUnmarshalers[nodeData.Type]
	if !ok {
		return nil, fmt.Errorf("unknown node type: %s", nodeData.Type)
	}

	return unmarshaler(data)
}

func Token(data []byte) (ast.Node, error) {
	var node ast.Token
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("Token", node, data, err)
	}
	return &node, nil
}

func Macro(data []byte) (ast.Node, error) {
	var node ast.Macro
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("Macro", node, data, err)
	}
	return &node, nil
}

func Include(data []byte) (ast.Node, error) {
	var node ast.Include
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("Include", node, data, err)
	}
	return &node, nil
}

func BasicLit(data []byte) (ast.Node, error) {
	var node ast.BasicLit
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("BasicLit", node, data, err)
	}
	return &node, nil
}

func BuiltinType(data []byte) (ast.Node, error) {
	var node ast.BuiltinType
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("BuiltinType", node, data, err)
	}
	return &node, nil
}

func Ident(data []byte) (ast.Node, error) {
	var node ast.Ident
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("Ident", node, data, err)
	}
	return &node, nil
}

func Variadic(data []byte) (ast.Node, error) {
	var node ast.Variadic
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, newDeserializeError("Variadic", node, data, err)
	}
	return &node, nil
}

func BlockPointerType(data []byte) (ast.Node, error) {
	return XType(data, &ast.BlockPointerType{})
}

func XType(data []byte, xType ast.Node) (ast.Node, error) {
	type XTypeTemp struct {
		X json.RawMessage
	}
	var xTypeData XTypeTemp
	if err := json.Unmarshal(data, &xTypeData); err != nil {
		return nil, newDeserializeError("XType", xTypeData, data, err)
	}

	if !isJSONNull(xTypeData.X) {
		xNode, err := Node(xTypeData.X)
		if err != nil {
			return nil, newUnmarshalFieldError("XType", xTypeData, "X", data, err)
		}
		expr, ok := xNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("XType", xNode, "ast.Expr")
		}
		switch v := xType.(type) {
		case *ast.PointerType:
			v.X = expr
		case *ast.LvalueRefType:
			v.X = expr
		case *ast.RvalueRefType:
			v.X = expr
		case *ast.BlockPointerType:
			fnType, ok := expr.(*ast.FuncType)
			if !ok {
				return nil, newUnexpectTypeError("XType", expr, "*ast.FuncType")
			}
			v.X = fnType
		default:
			return nil, newUnexpectTypeError("XType", xType, "*ast.PointerType, *ast.LvalueRefType, *ast.RvalueRefType")
		}
	}

	return xType, nil
}

func PointerType(data []byte) (ast.Node, error) {
	return XType(data, &ast.PointerType{})
}

func LvalueRefType(data []byte) (ast.Node, error) {
	return XType(data, &ast.LvalueRefType{})
}

func RvalueRefType(data []byte) (ast.Node, error) {
	return XType(data, &ast.RvalueRefType{})
}

func ArrayType(data []byte) (ast.Node, error) {
	type arrayTemp struct {
		Elt json.RawMessage
		Len json.RawMessage
	}
	var arrayData arrayTemp
	if err := json.Unmarshal(data, &arrayData); err != nil {
		return nil, newDeserializeError("ArrayType", arrayData, data, err)
	}

	var elt ast.Expr
	var len ast.Expr

	if !isJSONNull(arrayData.Elt) {
		eltNode, err := Node(arrayData.Elt)
		if err != nil {
			return nil, newUnmarshalFieldError("ArrayType", arrayData, "Elt", data, err)
		}
		var ok bool
		elt, ok = eltNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("ArrayType", eltNode, "ast.Expr")
		}
	}

	if !isJSONNull(arrayData.Len) {
		lenNode, err := Node(arrayData.Len)
		if err != nil {
			return nil, newUnmarshalFieldError("ArrayType", arrayData, "Len", data, err)
		}
		var ok bool
		len, ok = lenNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("ArrayType", lenNode, "ast.Expr")
		}
	}

	return &ast.ArrayType{
		Elt: elt,
		Len: len,
	}, nil
}

func Field(data []byte) (ast.Node, error) {
	type fieldTemp struct {
		Type     json.RawMessage
		Doc      *ast.CommentGroup
		Names    []*ast.Ident
		Comment  *ast.CommentGroup
		Access   ast.AccessSpecifier
		IsStatic bool
	}
	var fieldData fieldTemp
	if err := json.Unmarshal(data, &fieldData); err != nil {
		return nil, newDeserializeError("Field", fieldData, data, err)
	}

	var typ ast.Expr
	if !isJSONNull(fieldData.Type) {
		typeNode, err := Node(fieldData.Type)
		if err != nil {
			return nil, newUnmarshalFieldError("Field", fieldData, "Type", data, err)
		}
		var ok bool
		typ, ok = typeNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("Field", typeNode, "ast.Expr")
		}
	}

	return &ast.Field{
		Doc:      fieldData.Doc,
		Names:    fieldData.Names,
		Comment:  fieldData.Comment,
		Access:   fieldData.Access,
		IsStatic: fieldData.IsStatic,
		Type:     typ,
	}, nil
}

func FieldList(data []byte) (ast.Node, error) {
	type fieldListTemp struct {
		List []json.RawMessage
	}
	var fieldListData fieldListTemp
	if err := json.Unmarshal(data, &fieldListData); err != nil {
		return nil, newDeserializeError("FieldList", fieldListData, data, err)
	}

	list := []*ast.Field{}

	for _, fieldData := range fieldListData.List {
		fieldNode, err := Node(fieldData)
		if err != nil {
			return nil, newUnmarshalFieldError("FieldList", fieldListData, "List", data, err)
		}
		field, ok := fieldNode.(*ast.Field)
		if !ok {
			return nil, newUnexpectTypeError("FieldList", fieldNode, &ast.Field{})
		}
		list = append(list, field)
	}

	return &ast.FieldList{
		List: list,
	}, nil
}

func TagExpr(data []byte) (ast.Node, error) {
	type tagExprTemp struct {
		Name json.RawMessage
		Tag  ast.Tag
	}
	var tagExprData tagExprTemp
	if err := json.Unmarshal(data, &tagExprData); err != nil {
		return nil, newDeserializeError("TagExpr", tagExprData, data, err)
	}

	var name ast.Expr
	if !isJSONNull(tagExprData.Name) {
		nameNode, err := Node(tagExprData.Name)
		if err != nil {
			return nil, newUnmarshalFieldError("TagExpr", tagExprData, "Name", data, err)
		}
		var ok bool
		name, ok = nameNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("TagExpr", nameNode, "ast.Expr")
		}
	}

	return &ast.TagExpr{
		Tag:  tagExprData.Tag,
		Name: name,
	}, nil
}

func ScopingExpr(data []byte) (ast.Node, error) {
	type scopingExprTemp struct {
		Parent json.RawMessage
		X      *ast.Ident
	}
	var scopingExprData scopingExprTemp
	if err := json.Unmarshal(data, &scopingExprData); err != nil {
		return nil, newDeserializeError("ScopingExpr", scopingExprData, data, err)
	}

	var parent ast.Expr
	if !isJSONNull(scopingExprData.Parent) {
		parentNode, err := Node(scopingExprData.Parent)
		if err != nil {
			return nil, newUnmarshalFieldError("ScopingExpr", scopingExprData, "Parent", data, err)
		}
		var ok bool
		parent, ok = parentNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("ScopingExpr", parentNode, "ast.Expr")
		}
	}

	return &ast.ScopingExpr{
		Parent: parent,
		X:      scopingExprData.X,
	}, nil
}

func EnumItem(data []byte) (ast.Node, error) {
	type enumItemTemp struct {
		Name  *ast.Ident
		Value json.RawMessage
	}
	var enumItemData enumItemTemp

	if err := json.Unmarshal(data, &enumItemData); err != nil {
		return nil, newDeserializeError("EnumItem", enumItemData, data, err)
	}

	var value ast.Expr
	if !isJSONNull(enumItemData.Value) {
		valueNode, err := Node(enumItemData.Value)
		if err != nil {
			return nil, newUnmarshalFieldError("EnumItem", enumItemData, "Value", data, err)
		}
		var ok bool
		value, ok = valueNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("EnumItem", valueNode, "ast.Expr")
		}
	}

	return &ast.EnumItem{
		Name:  enumItemData.Name,
		Value: value,
	}, nil
}

func EnumType(data []byte) (ast.Node, error) {
	type enumTypeTemp struct {
		Items []json.RawMessage
	}
	var enumTypeData enumTypeTemp
	if err := json.Unmarshal(data, &enumTypeData); err != nil {
		return nil, newDeserializeError("EnumType", enumTypeData, data, err)
	}

	items := []*ast.EnumItem{}
	for _, itemData := range enumTypeData.Items {
		itemNode, err := Node(itemData)
		if err != nil {
			return nil, newUnmarshalFieldError("EnumType", enumTypeData, "Items", data, err)
		}
		item, ok := itemNode.(*ast.EnumItem)
		if !ok {
			return nil, newUnexpectTypeError("EnumType", itemNode, &ast.EnumItem{})
		}
		items = append(items, item)
	}

	return &ast.EnumType{
		Items: items,
	}, nil
}

func RecordType(data []byte) (ast.Node, error) {
	type recordTypeTemp struct {
		Tag     ast.Tag
		Fields  json.RawMessage
		Methods []json.RawMessage
	}
	var recordTypeData recordTypeTemp
	if err := json.Unmarshal(data, &recordTypeData); err != nil {
		return nil, newDeserializeError("RecordType", recordTypeData, data, err)
	}

	var fields *ast.FieldList
	var methods []*ast.FuncDecl
	if !isJSONNull(recordTypeData.Fields) {
		fieldsNode, err := Node(recordTypeData.Fields)
		if err != nil {
			return nil, newUnmarshalFieldError("RecordType", recordTypeData, "Fields", data, err)
		}
		var ok bool
		fields, ok = fieldsNode.(*ast.FieldList)
		if !ok {
			return nil, newUnexpectTypeError("RecordType", fieldsNode, &ast.FieldList{})
		}
	}

	for _, methodData := range recordTypeData.Methods {
		methodNode, err := Node(methodData)
		if err != nil {
			return nil, newUnmarshalFieldError("RecordType", recordTypeData, "Methods", data, err)
		}
		method, ok := methodNode.(*ast.FuncDecl)
		if !ok {
			return nil, newUnexpectTypeError("RecordType", methodNode, &ast.FuncDecl{})
		}
		methods = append(methods, method)
	}

	return &ast.RecordType{
		Tag:     recordTypeData.Tag,
		Fields:  fields,
		Methods: methods,
	}, nil
}

func FuncType(data []byte) (ast.Node, error) {
	type funcTypeTemp struct {
		Params json.RawMessage
		Ret    json.RawMessage
	}
	var funcTypeData funcTypeTemp
	if err := json.Unmarshal(data, &funcTypeData); err != nil {
		return nil, newDeserializeError("FuncType", funcTypeData, data, err)
	}

	var params *ast.FieldList
	var ret ast.Expr
	if !isJSONNull(funcTypeData.Params) {
		paramsNode, err := Node(funcTypeData.Params)
		if err != nil {
			return nil, newUnmarshalFieldError("FuncType", funcTypeData, "Params", data, err)
		}
		var ok bool
		params, ok = paramsNode.(*ast.FieldList)
		if !ok {
			return nil, newUnexpectTypeError("FuncType", paramsNode, &ast.FieldList{})
		}
	}

	if !isJSONNull(funcTypeData.Ret) {
		retNode, err := Node(funcTypeData.Ret)
		if err != nil {
			return nil, newUnmarshalFieldError("FuncType", funcTypeData, "Ret", data, err)
		}
		var ok bool
		ret, ok = retNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("FuncType", retNode, "ast.Expr")
		}
	}

	return &ast.FuncType{
		Params: params,
		Ret:    ret,
	}, nil
}

func FuncDecl(data []byte) (ast.Node, error) {
	type funcDeclTemp struct {
		MangledName   string
		Type          json.RawMessage
		IsInline      bool
		IsStatic      bool
		IsConst       bool
		IsExplicit    bool
		IsConstructor bool
		IsDestructor  bool
		IsVirtual     bool
		IsOverride    bool
	}
	var funcDeclData funcDeclTemp
	if err := json.Unmarshal(data, &funcDeclData); err != nil {
		return nil, newDeserializeError("FuncDecl", funcDeclData, data, err)
	}

	var typ *ast.FuncType
	if !isJSONNull(funcDeclData.Type) {
		typeNode, err := Node(funcDeclData.Type)
		if err != nil {
			return nil, newUnmarshalFieldError("FuncDecl", funcDeclData, "Type", data, err)
		}
		var ok bool
		typ, ok = typeNode.(*ast.FuncType)
		if !ok {
			return nil, newUnexpectTypeError("FuncDecl", typeNode, &ast.FuncType{})
		}
	}

	declBase, err := declBase(data)
	if err != nil {
		return nil, err
	}

	return &ast.FuncDecl{
		Object:        declBase,
		Type:          typ,
		MangledName:   funcDeclData.MangledName,
		IsInline:      funcDeclData.IsInline,
		IsStatic:      funcDeclData.IsStatic,
		IsConst:       funcDeclData.IsConst,
		IsExplicit:    funcDeclData.IsExplicit,
		IsConstructor: funcDeclData.IsConstructor,
		IsDestructor:  funcDeclData.IsDestructor,
		IsVirtual:     funcDeclData.IsVirtual,
		IsOverride:    funcDeclData.IsOverride,
	}, nil
}

func TypeDecl(data []byte) (ast.Node, error) {
	type typeDeclTemp struct {
		Type json.RawMessage
	}
	var typeDeclData typeDeclTemp
	if err := json.Unmarshal(data, &typeDeclData); err != nil {
		return nil, newDeserializeError("TypeDecl", typeDeclData, data, err)
	}

	var typ *ast.RecordType
	if !isJSONNull(typeDeclData.Type) {
		typeNode, err := Node(typeDeclData.Type)
		if err != nil {
			return nil, newUnmarshalFieldError("TypeDecl", typeDeclData, "Type", data, err)
		}
		var ok bool
		typ, ok = typeNode.(*ast.RecordType)
		if !ok {
			return nil, newUnexpectTypeError("TypeDecl", typeNode, &ast.RecordType{})
		}
	}

	declBase, err := declBase(data)
	if err != nil {
		return nil, err
	}

	return &ast.TypeDecl{
		Object: declBase,
		Type:   typ,
	}, nil
}

func TypeDefDecl(data []byte) (ast.Node, error) {
	type typeDefDeclTemp struct {
		Type json.RawMessage
	}
	var typeDefDeclData typeDefDeclTemp
	if err := json.Unmarshal(data, &typeDefDeclData); err != nil {
		return nil, newDeserializeError("TypeDefDecl", typeDefDeclData, data, err)
	}

	var typ ast.Expr
	if !isJSONNull(typeDefDeclData.Type) {
		typeNode, err := Node(typeDefDeclData.Type)
		if err != nil {
			return nil, newUnmarshalFieldError("TypeDefDecl", typeDefDeclData, "Type", data, err)
		}
		var ok bool
		typ, ok = typeNode.(ast.Expr)
		if !ok {
			return nil, newUnexpectTypeError("TypeDefDecl", typeNode, "ast.Expr")
		}
	}

	declBase, err := declBase(data)
	if err != nil {
		return nil, err
	}

	return &ast.TypedefDecl{
		Object: declBase,
		Type:   typ,
	}, nil
}

func EnumTypeDecl(data []byte) (ast.Node, error) {
	type enumTypeDeclTemp struct {
		Type json.RawMessage
	}
	var enumTypeDeclData enumTypeDeclTemp
	if err := json.Unmarshal(data, &enumTypeDeclData); err != nil {
		return nil, newDeserializeError("EnumTypeDecl", enumTypeDeclData, data, err)
	}

	var typ *ast.EnumType
	if !isJSONNull(enumTypeDeclData.Type) {
		typeNode, err := Node(enumTypeDeclData.Type)
		if err != nil {
			return nil, newUnmarshalFieldError("EnumTypeDecl", enumTypeDeclData, "Type", data, err)
		}
		var ok bool
		typ, ok = typeNode.(*ast.EnumType)
		if !ok {
			return nil, newUnexpectTypeError("EnumTypeDecl", typeNode, &ast.EnumType{})
		}
	}

	declBase, err := declBase(data)
	if err != nil {
		return nil, err
	}

	return &ast.EnumTypeDecl{
		Object: declBase,
		Type:   typ,
	}, nil
}

func declBase(data []byte) (ast.Object, error) {
	type declBaseTemp struct {
		Loc    *ast.Location
		Doc    *ast.CommentGroup
		Name   *ast.Ident
		Parent json.RawMessage
	}
	var declBaseData declBaseTemp
	if err := json.Unmarshal(data, &declBaseData); err != nil {
		return ast.Object{}, newDeserializeError("declBase", declBaseData, data, err)
	}

	var parent ast.Expr
	if !isJSONNull(declBaseData.Parent) {
		parentNode, err := Node(declBaseData.Parent)
		if err != nil {
			return ast.Object{}, newUnmarshalFieldError("declBase", declBaseData, "Parent", data, err)
		}
		var ok bool
		parent, ok = parentNode.(ast.Expr)
		if !ok {
			return ast.Object{}, newUnexpectTypeError("declBase", parentNode, "ast.Expr")
		}
	}

	return ast.Object{
		Loc:    declBaseData.Loc,
		Doc:    declBaseData.Doc,
		Name:   declBaseData.Name,
		Parent: parent,
	}, nil
}

func isJSONNull(data json.RawMessage) bool {
	return len(data) == 4 && string(data) == "null"
}

// DeserializeError represents an error that occurs during json.Unmarshal.
// It provides context about where the error occurred and what was being unmarshaled.
type DeserializeError struct {
	Func       string // function name
	TargetType any
	Field      string // optional, only for unmarshal node in a struct
	Data       string // origin raw json data
	Err        error  // unmarshal error message
}

func (e *DeserializeError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("unmarshal error in %s when converting %s of %v: %s\ninput: %s",
			e.Func, e.Field, reflect.TypeOf(e.TargetType), e.Err.Error(), e.Data)
	}
	return fmt.Sprintf("unmarshal error in %s into %v: %s\ninput: %s",
		e.Func, reflect.TypeOf(e.TargetType), e.Err.Error(), e.Data)
}

func newDeserializeError(funcName string, targetType any, data []byte, err error) *DeserializeError {
	return &DeserializeError{
		Func:       funcName,
		TargetType: targetType,
		Data:       string(data),
		Err:        err,
	}
}

func newUnmarshalFieldError(funcName string, targetType any, field string, data []byte, err error) *DeserializeError {
	return &DeserializeError{
		Func:       funcName,
		TargetType: targetType,
		Field:      field,
		Data:       string(data),
		Err:        err,
	}
}

type UnexpectTypeError struct {
	Func     string
	GotType  any
	WantType any
}

func (e *UnexpectTypeError) Error() string {
	if reflect.TypeOf(e.WantType).Kind() == reflect.String {
		return fmt.Sprintf("unmarshal error in %s: got %v, want %s", e.Func, reflect.TypeOf(e.GotType), e.WantType)
	}
	return fmt.Sprintf("unmarshal error in %s: got %v, want %v", e.Func, reflect.TypeOf(e.GotType), reflect.TypeOf(e.WantType))
}

func newUnexpectTypeError(funcName string, gotType any, wantType any) *UnexpectTypeError {
	return &UnexpectTypeError{
		Func:     funcName,
		GotType:  gotType,
		WantType: wantType,
	}
}
