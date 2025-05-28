package parser

import (
	"github.com/goplus/llcppg/ast"
)

func MarshalDeclList(list []ast.Decl) []map[string]any {
	var root []map[string]any
	for _, item := range list {
		root = append(root, MarshalASTDecl(item))
	}
	return root
}

func MarshalFieldList(list []*ast.Field) []map[string]any {
	if list == nil {
		return nil
	}
	var root []map[string]any

	for _, item := range list {
		root = append(root, MarshalASTExpr(item))
	}
	return root
}

func MarshalIncludeList(list []*ast.Include) []map[string]any {
	var root []map[string]any
	for _, item := range list {
		root = append(root, map[string]any{
			"_Type": "Include",
			"Path":  item.Path,
		})
	}
	return root
}

func MarshalMacroList(list []*ast.Macro) []map[string]any {
	var root []map[string]any

	for _, item := range list {
		root = append(root, map[string]any{
			"_Type":  "Macro",
			"Loc":    MarshalLocation(item.Loc),
			"Name":   item.Name,
			"Tokens": MarshalTokenList(item.Tokens),
		})
	}
	return root
}

func MarshalTokenList(list []*ast.Token) []map[string]any {
	if list == nil {
		return nil
	}
	var root []map[string]any
	for _, item := range list {
		root = append(root, MarshalToken(item))
	}
	return root
}

func MarshalIdentList(list []*ast.Ident) []map[string]any {
	if list == nil {
		return nil
	}
	var root []map[string]any

	for _, item := range list {
		root = append(root, MarshalASTExpr(item))
	}
	return root
}

func MarshalASTFile(file *ast.File) map[string]any {
	root := make(map[string]any)
	root["_Type"] = "File"
	root["decls"] = MarshalDeclList(file.Decls)
	root["includes"] = MarshalIncludeList(file.Includes)
	root["macros"] = MarshalMacroList(file.Macros)
	return root
}

func MarshalToken(tok *ast.Token) map[string]any {
	root := make(map[string]any)
	root["_Type"] = "Token"
	root["Token"] = uint(tok.Token)
	root["Lit"] = tok.Lit
	return root
}

func MarshalASTDecl(decl ast.Decl) map[string]any {
	if decl == nil {
		return nil
	}
	root := make(map[string]any)

	switch d := decl.(type) {
	case *ast.EnumTypeDecl:
		root["_Type"] = "EnumTypeDecl"
		MarshalObject(d.Object, root)
		root["Type"] = MarshalASTExpr(d.Type)
	case *ast.TypedefDecl:
		root["_Type"] = "TypedefDecl"
		MarshalObject(d.Object, root)
		root["Type"] = MarshalASTExpr(d.Type)
	case *ast.FuncDecl:
		root["_Type"] = "FuncDecl"
		MarshalObject(d.Object, root)
		root["MangledName"] = d.MangledName
		root["Type"] = MarshalASTExpr(d.Type)
		root["IsInline"] = d.IsInline
		root["IsStatic"] = d.IsStatic
		root["IsConst"] = d.IsConst
		root["IsExplicit"] = d.IsExplicit
		root["IsConstructor"] = d.IsConstructor
		root["IsDestructor"] = d.IsDestructor
		root["IsVirtual"] = d.IsVirtual
		root["IsOverride"] = d.IsOverride
	case *ast.TypeDecl:
		root["_Type"] = "TypeDecl"
		MarshalObject(d.Object, root)
		root["Type"] = MarshalASTExpr(d.Type)
	}
	return root
}

func MarshalObject(decl ast.Object, root map[string]any) {
	root["Loc"] = MarshalLocation(decl.Loc)
	root["Doc"] = MarshalASTExpr(decl.Doc)
	root["Parent"] = MarshalASTExpr(decl.Parent)
	root["Name"] = MarshalASTExpr(decl.Name)
}

func MarshalLocation(loc *ast.Location) map[string]any {
	if loc == nil {
		return nil
	}
	root := make(map[string]any)
	root["_Type"] = "Location"
	root["File"] = loc.File
	return root
}

func MarshalASTExpr(t ast.Expr) map[string]any {
	if t == nil {
		return nil
	}

	root := make(map[string]any)

	switch d := t.(type) {
	case *ast.EnumType:
		root["_Type"] = "EnumType"
		var items []map[string]any
		for _, e := range d.Items {
			items = append(items, MarshalASTExpr(e))
		}
		root["Items"] = items
	case *ast.EnumItem:
		root["_Type"] = "EnumItem"
		root["Name"] = MarshalASTExpr(d.Name)
		root["Value"] = MarshalASTExpr(d.Value)
	case *ast.RecordType:
		root["_Type"] = "RecordType"
		root["Tag"] = uint(d.Tag)
		root["Fields"] = MarshalASTExpr(d.Fields)
		var methods []map[string]any
		for _, m := range d.Methods {
			methods = append(methods, MarshalASTDecl(m))
		}
		root["Methods"] = methods
	case *ast.FuncType:
		root["_Type"] = "FuncType"
		root["Params"] = MarshalASTExpr(d.Params)
		root["Ret"] = MarshalASTExpr(d.Ret)
	case *ast.FieldList:
		root["_Type"] = "FieldList"
		root["List"] = MarshalFieldList(d.List)
	case *ast.Field:
		root["_Type"] = "Field"
		root["Type"] = MarshalASTExpr(d.Type)
		root["Doc"] = MarshalASTExpr(d.Doc)
		root["Comment"] = MarshalASTExpr(d.Comment)
		root["IsStatic"] = d.IsStatic
		root["Access"] = uint(d.Access)
		root["Names"] = MarshalIdentList(d.Names)
	case *ast.Variadic:
		root["_Type"] = "Variadic"
	case *ast.Ident:
		root["_Type"] = "Ident"
		if d == nil {
			return nil
		}
		root["Name"] = d.Name
	case *ast.TagExpr:
		root["_Type"] = "TagExpr"
		root["Name"] = MarshalASTExpr(d.Name)
		root["Tag"] = uint(d.Tag)
	case *ast.BasicLit:
		root["_Type"] = "BasicLit"
		root["Kind"] = uint(d.Kind)
		root["Value"] = d.Value
	case *ast.LvalueRefType:
		root["_Type"] = "LvalueRefType"
		root["X"] = MarshalASTExpr(d.X)
	case *ast.RvalueRefType:
		root["_Type"] = "RvalueRefType"
		root["X"] = MarshalASTExpr(d.X)
	case *ast.PointerType:
		root["_Type"] = "PointerType"
		root["X"] = MarshalASTExpr(d.X)
	case *ast.BlockPointerType:
		root["_Type"] = "BlockPointerType"
		root["X"] = MarshalASTExpr(d.X)
	case *ast.ArrayType:
		root["_Type"] = "ArrayType"
		root["Elt"] = MarshalASTExpr(d.Elt)
		root["Len"] = MarshalASTExpr(d.Len)
	case *ast.BuiltinType:
		root["_Type"] = "BuiltinType"
		root["Kind"] = uint(d.Kind)
		root["Flags"] = uint(d.Flags)
	case *ast.Comment:
		root["_Type"] = "Comment"
		if d == nil {
			return nil
		}
		root["Text"] = d.Text
	case *ast.CommentGroup:
		root["_Type"] = "CommentGroup"
		if d == nil {
			return nil
		}
		var list []map[string]any
		for _, c := range d.List {
			list = append(list, MarshalASTExpr(c))
		}
		root["List"] = list
	case *ast.ScopingExpr:
		root["_Type"] = "ScopingExpr"
		root["X"] = MarshalASTExpr(d.X)
		root["Parent"] = MarshalASTExpr(d.Parent)
	default:
		return nil
	}
	return root
}
