package parser

import (
	"github.com/goplus/llcppg/ast"
)

func XMarshalDeclList(list []ast.Decl) []map[string]any {
	var root []map[string]any
	for _, item := range list {
		root = append(root, XMarshalASTDecl(item))
	}
	return root
}

func XMarshalFieldList(list []*ast.Field) []map[string]any {
	if list == nil {
		return nil
	}
	var root []map[string]any

	for _, item := range list {
		root = append(root, XMarshalASTExpr(item))
	}
	return root
}

func XMarshalIncludeList(list []*ast.Include) []map[string]any {
	var root []map[string]any
	for _, item := range list {
		root = append(root, map[string]any{
			"_Type": "Include",
			"Path":  item.Path,
		})
	}
	return root
}

func XMarshalMacroList(list []*ast.Macro) []map[string]any {
	var root []map[string]any

	for _, item := range list {
		root = append(root, map[string]any{
			"_Type":  "Macro",
			"Loc":    XMarshalLocation(item.Loc),
			"Name":   item.Name,
			"Tokens": XMarshalTokenList(item.Tokens),
		})
	}
	return root
}

func XMarshalTokenList(list []*ast.Token) []map[string]any {
	if list == nil {
		return nil
	}
	var root []map[string]any
	for _, item := range list {
		root = append(root, XMarshalToken(item))
	}
	return root
}

func XMarshalIdentList(list []*ast.Ident) []map[string]any {
	if list == nil {
		return nil
	}
	var root []map[string]any

	for _, item := range list {
		root = append(root, XMarshalASTExpr(item))
	}
	return root
}

func XMarshalASTFile(file *ast.File) map[string]any {
	return map[string]any{
		"_Type":    "File",
		"decls":    XMarshalDeclList(file.Decls),
		"includes": XMarshalIncludeList(file.Includes),
		"macros":   XMarshalMacroList(file.Macros),
	}
}

func XMarshalToken(tok *ast.Token) map[string]any {
	return map[string]any{
		"_Type": "Token",
		"Token": uint(tok.Token),
		"Lit":   tok.Lit,
	}
}

func XMarshalASTDecl(decl ast.Decl) map[string]any {
	if decl == nil {
		return nil
	}
	root := make(map[string]any)

	switch d := decl.(type) {
	case *ast.EnumTypeDecl:
		root["_Type"] = "EnumTypeDecl"
		XMarshalObject(d.Object, root)
		root["Type"] = XMarshalASTExpr(d.Type)
	case *ast.TypedefDecl:
		root["_Type"] = "TypedefDecl"
		XMarshalObject(d.Object, root)
		root["Type"] = XMarshalASTExpr(d.Type)
	case *ast.FuncDecl:
		root["_Type"] = "FuncDecl"
		XMarshalObject(d.Object, root)
		root["MangledName"] = d.MangledName
		root["Type"] = XMarshalASTExpr(d.Type)
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
		XMarshalObject(d.Object, root)
		root["Type"] = XMarshalASTExpr(d.Type)
	}
	return root
}

func XMarshalObject(decl ast.Object, root map[string]any) {
	root["Loc"] = XMarshalLocation(decl.Loc)
	root["Doc"] = XMarshalASTExpr(decl.Doc)
	root["Parent"] = XMarshalASTExpr(decl.Parent)
	root["Name"] = XMarshalASTExpr(decl.Name)
}

func XMarshalLocation(loc *ast.Location) map[string]any {
	if loc == nil {
		return nil
	}
	root := make(map[string]any)
	root["_Type"] = "Location"
	root["File"] = loc.File
	return root
}

func XMarshalASTExpr(t ast.Expr) map[string]any {
	if t == nil {
		return nil
	}

	root := make(map[string]any)

	switch d := t.(type) {
	case *ast.EnumType:
		root["_Type"] = "EnumType"
		var items []map[string]any
		for _, e := range d.Items {
			items = append(items, XMarshalASTExpr(e))
		}
		root["Items"] = items
	case *ast.EnumItem:
		root["_Type"] = "EnumItem"
		root["Name"] = XMarshalASTExpr(d.Name)
		root["Value"] = XMarshalASTExpr(d.Value)
	case *ast.RecordType:
		root["_Type"] = "RecordType"
		root["Tag"] = uint(d.Tag)
		root["Fields"] = XMarshalASTExpr(d.Fields)
		var methods []map[string]any
		for _, m := range d.Methods {
			methods = append(methods, XMarshalASTDecl(m))
		}
		root["Methods"] = methods
	case *ast.FuncType:
		root["_Type"] = "FuncType"
		root["Params"] = XMarshalASTExpr(d.Params)
		root["Ret"] = XMarshalASTExpr(d.Ret)
	case *ast.FieldList:
		root["_Type"] = "FieldList"
		if d == nil {
			return nil
		}
		root["List"] = XMarshalFieldList(d.List)
	case *ast.Field:
		root["_Type"] = "Field"
		root["Type"] = XMarshalASTExpr(d.Type)
		root["Doc"] = XMarshalASTExpr(d.Doc)
		root["Comment"] = XMarshalASTExpr(d.Comment)
		root["IsStatic"] = d.IsStatic
		root["Access"] = uint(d.Access)
		root["Names"] = XMarshalIdentList(d.Names)
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
		root["Name"] = XMarshalASTExpr(d.Name)
		root["Tag"] = uint(d.Tag)
	case *ast.BasicLit:
		root["_Type"] = "BasicLit"
		root["Kind"] = uint(d.Kind)
		root["Value"] = d.Value
	case *ast.LvalueRefType:
		root["_Type"] = "LvalueRefType"
		root["X"] = XMarshalASTExpr(d.X)
	case *ast.RvalueRefType:
		root["_Type"] = "RvalueRefType"
		root["X"] = XMarshalASTExpr(d.X)
	case *ast.PointerType:
		root["_Type"] = "PointerType"
		root["X"] = XMarshalASTExpr(d.X)
	case *ast.BlockPointerType:
		root["_Type"] = "BlockPointerType"
		root["X"] = XMarshalASTExpr(d.X)
	case *ast.ArrayType:
		root["_Type"] = "ArrayType"
		root["Elt"] = XMarshalASTExpr(d.Elt)
		root["Len"] = XMarshalASTExpr(d.Len)
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
			list = append(list, XMarshalASTExpr(c))
		}
		root["List"] = list
	case *ast.ScopingExpr:
		root["_Type"] = "ScopingExpr"
		root["X"] = XMarshalASTExpr(d.X)
		root["Parent"] = XMarshalASTExpr(d.Parent)
	default:
		return nil
	}
	return root
}
