package unmarshal_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/internal/unmarshal"
)

func TestUnmarshalFile(t *testing.T) {
	files := `{
			"_Type":	"File",
			"decls":	[{
							"_Type":	"FuncDecl",
							"Loc":	{
								"_Type":	"Location",
								"File":	"temp.h"
							},
							"Doc":	{
								"_Type":	"CommentGroup",
								"List":	[]
							},
							"Parent":	null,
							"Name":	{
								"_Type":	"Ident",
								"Name":	"foo"
							},
							"Type":	{
								"_Type":	"FuncType",
								"Params":	{
									"_Type":	"FieldList",
									"List":	[{
											"_Type":	"Field",
											"Type":	{
												"_Type":	"Variadic"
											},
											"Doc":	null,
											"Comment":	null,
											"IsStatic":	false,
											"Access":	0,
											"Names":	[]
										}, {
											"_Type":	"Field",
											"Type":	{
												"_Type":	"BuiltinType",
												"Kind":	6,
												"Flags":	0
											},
											"Doc":	{
												"_Type":	"CommentGroup",
												"List":	[]
											},
											"Comment":	{
												"_Type":	"CommentGroup",
												"List":	[]
											},
											"IsStatic":	false,
											"Access":	0,
											"Names":	[{
													"_Type":	"Ident",
													"Name":	"a"
												}]
										}]
								},
								"Ret":	{
									"_Type":	"BuiltinType",
									"Kind":	0,
									"Flags":	0
								}
							},
							"IsInline":	false,
							"IsStatic":	false,
							"IsConst":	false,
							"IsExplicit":	false,
							"IsConstructor":	false,
							"IsDestructor":	false,
							"IsVirtual":	false,
							"IsOverride":	false
						}
			],
			"includes":	[],
			"macros":	[{
					"_Type":	"Macro",
					"Name":	"OK",
					"Tokens":	[{
							"_Type":	"Token",
							"Token":	3,
							"Lit":	"OK"
						}, {
							"_Type":	"Token",
							"Token":	4,
							"Lit":	"1"
						}]
				}]
		}`

	expected := &ast.File{
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Object: ast.Object{
					Loc: &ast.Location{
						File: "temp.h",
					},
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{},
					},
					Name: &ast.Ident{Name: "foo"},
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type:  &ast.Variadic{},
								Names: []*ast.Ident{},
							}, {
								Doc: &ast.CommentGroup{
									List: []*ast.Comment{},
								},
								Comment: &ast.CommentGroup{
									List: []*ast.Comment{},
								},
								Type: &ast.BuiltinType{
									Kind:  6,
									Flags: 0,
								},
								Names: []*ast.Ident{
									{Name: "a"},
								},
							},
						},
					},
					Ret: &ast.BuiltinType{
						Kind:  0,
						Flags: 0,
					},
				},
			},
		},
		Includes: []*ast.Include{},
		Macros: []*ast.Macro{
			{
				Name: "OK",
				Tokens: []*ast.Token{
					{Token: 3, Lit: "OK"},
					{Token: 4, Lit: "1"},
				},
			},
		},
	}

	fileSet, err := unmarshal.File([]byte(files))

	if err != nil {
		t.Fatalf("UnmarshalNode failed: %v", err)
	}

	resultJSON, err := json.MarshalIndent(fileSet, "", " ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}

	expectedJSON, err := json.MarshalIndent(expected, "", " ")

	if err != nil {
		t.Fatalf("Failed to marshal expected result to JSON: %v", err)
	}

	if string(resultJSON) != string(expectedJSON) {
		t.Errorf("JSON mismatch.\nExpected: %s\nGot: %s", string(expectedJSON), string(resultJSON))
	}
}

func TestUnmarshalNode(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected ast.Node
	}{
		{
			name:     "Token",
			json:     `{ "_Type":	"Token","Token":	3,"Lit":	"DEBUG"}`,
			expected: &ast.Token{Token: 3, Lit: "DEBUG"},
		},
		{
			name:     "BuiltinType",
			json:     `{"_Type": "BuiltinType", "Kind": 6, "Flags": 0}`,
			expected: &ast.BuiltinType{Kind: ast.TypeKind(6), Flags: ast.TypeFlag(0)},
		},
		{
			name: "PointerType",
			json: `{
					"_Type":	"PointerType",
					"X":	{
						"_Type":	"BuiltinType",
						"Kind":	2,
						"Flags":	1
					}
				}`,
			expected: &ast.PointerType{
				X: &ast.BuiltinType{
					Kind:  2,
					Flags: 1,
				},
			},
		},
		{
			name: "PointerType",
			json: `{
					"_Type":	"PointerType",
					"X":	{
						"_Type":	"PointerType",
						"X":	{
							"_Type":	"PointerType",
							"X":	{
								"_Type":	"BuiltinType",
								"Kind":	2,
								"Flags":	1
							}
						}
					}
				}`,
			expected: &ast.PointerType{
				X: &ast.PointerType{
					X: &ast.PointerType{
						X: &ast.BuiltinType{
							Kind:  2,
							Flags: 1,
						},
					},
				}},
		},
		{
			name: "ArrayType",
			json: `{
					"_Type":	"ArrayType",
					"Elt":	{
						"_Type":	"BuiltinType",
						"Kind":	2,
						"Flags":	1
					},
					"Len":	null
				}`,
			expected: &ast.ArrayType{
				Elt: &ast.BuiltinType{
					Kind:  2,
					Flags: 1,
				},
				Len: nil,
			},
		},
		{
			name: "ArrayType",
			json: `{
					"_Type":	"ArrayType",
					"Elt":	{
						"_Type":	"BuiltinType",
						"Kind":	2,
						"Flags":	1
					},
					"Len":	{
						"_Type":	"BasicLit",
						"Kind":	0,
						"Value":	"10"
					}
				}`,
			expected: &ast.ArrayType{
				Elt: &ast.BuiltinType{
					Kind:  2,
					Flags: 1,
				},
				Len: &ast.BasicLit{
					Kind:  0,
					Value: "10",
				},
			},
		},
		{
			name: "ArrayType",
			json: `{
					"_Type":	"ArrayType",
					"Elt":	{
						"_Type":	"ArrayType",
						"Elt":	{
							"_Type":	"BuiltinType",
							"Kind":	2,
							"Flags":	1
						},
						"Len":	{
							"_Type":	"BasicLit",
							"Kind":	0,
							"Value":	"4"
						}
					},
					"Len":	{
						"_Type":	"BasicLit",
						"Kind":	0,
						"Value":	"3"
					}
				}`,
			expected: &ast.ArrayType{
				Elt: &ast.ArrayType{
					Elt: &ast.BuiltinType{
						Kind:  2,
						Flags: 1,
					},
					Len: &ast.BasicLit{
						Kind:  0,
						Value: "4",
					},
				},
				Len: &ast.BasicLit{
					Kind:  0,
					Value: "3",
				},
			},
		},
		{
			name:     "Variadic",
			json:     `{"_Type": "Variadic"}`,
			expected: &ast.Variadic{},
		},
		{
			name: "LvalueRefType",
			json: `{
						"_Type":	"LvalueRefType",
						"X":	{
							"_Type":	"BuiltinType",
							"Kind":	6,
							"Flags":	0
						}
					}`,
			expected: &ast.LvalueRefType{
				X: &ast.BuiltinType{
					Kind:  6,
					Flags: 0,
				},
			},
		},
		{
			name: "RvalueRefType",
			json: `{
						"_Type":	"RvalueRefType",
						"X":	{
							"_Type":	"BuiltinType",
							"Kind":	6,
							"Flags":	0
						}
					}`,
			expected: &ast.RvalueRefType{
				X: &ast.BuiltinType{
					Kind:  6,
					Flags: 0,
				},
			},
		},
		{
			name: "Ident",
			json: `{
						"_Type":	"Ident",
						"Name":	"Foo"
					}`,
			expected: &ast.Ident{
				Name: "Foo",
			},
		},
		{
			name: "ScopingExpr",
			json: `{
					"_Type":	"ScopingExpr",
					"X":	{
						"_Type":	"Ident",
						"Name":	"b"
					},
					"Parent":	{
						"_Type":	"Ident",
						"Name":	"a"
					}
				}`,
			expected: &ast.ScopingExpr{
				X: &ast.Ident{
					Name: "b",
				},
				Parent: &ast.Ident{
					Name: "a",
				},
			},
		},
		{
			name: "ScopingExpr",
			json: `{
					"_Type":	"ScopingExpr",
					"X":	{
						"_Type":	"Ident",
						"Name":	"c"
					},
					"Parent":	{
						"_Type":	"ScopingExpr",
						"X":	{
							"_Type":	"Ident",
							"Name":	"b"
						},
						"Parent":	{
							"_Type":	"Ident",
							"Name":	"a"
						}
					}
				}`,
			expected: &ast.ScopingExpr{
				X: &ast.Ident{
					Name: "c",
				},
				Parent: &ast.ScopingExpr{
					X: &ast.Ident{
						Name: "b",
					},
					Parent: &ast.Ident{
						Name: "a",
					},
				},
			},
		},
		{
			name: "TagExpr",
			json: `{
					"_Type":	"TagExpr",
					"Name":	{
						"_Type":	"Ident",
						"Name":	"Foo"
					},
					"Tag":	0
				}`,
			expected: &ast.TagExpr{
				Tag: 0,
				Name: &ast.Ident{
					Name: "Foo",
				},
			},
		},
		{
			name: "Field",
			json: `{
                "_Type": "Field",
                "Type": {"_Type": "Variadic"},
                "Doc": {
					"_Type":	"CommentGroup",
					"List":	[{
							"_Type":	"Comment",
							"Text":	"/// doc"
						}]
				},
                "Comment": null,
                "IsStatic": false,
                "Access": 0,
                	"Names":	[{
						"_Type":	"Ident",
					"Name":	"a"
				}]
            }`,
			expected: &ast.Field{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{{Text: "/// doc"}},
				},
				Type:     &ast.Variadic{},
				IsStatic: false,
				Access:   ast.AccessSpecifier(0),
				Names:    []*ast.Ident{{Name: "a"}},
			},
		},
		{
			name: "FieldList",
			json: `{
					"_Type":	"FieldList",
					"List":	[{
							"_Type":	"Field",
							"Type":	{
								"_Type":	"BuiltinType",
								"Kind":	6,
								"Flags":	0
							},
							"Doc":	{
								"_Type":	"CommentGroup",
								"List":	[]
							},
							"Comment":	{
								"_Type":	"CommentGroup",
								"List":	[]
							},
							"IsStatic":	false,
							"Access":	1,
							"Names":	[{
									"_Type":	"Ident",
									"Name":	"x"
								}]
						}]
				}`,
			expected: &ast.FieldList{
				List: []*ast.Field{
					{
						Doc: &ast.CommentGroup{
							List: []*ast.Comment{},
						},
						Comment: &ast.CommentGroup{
							List: []*ast.Comment{},
						},
						Type: &ast.BuiltinType{
							Kind:  6,
							Flags: 0,
						},
						Access: ast.AccessSpecifier(1),
						Names:  []*ast.Ident{{Name: "x"}},
					},
				},
			},
		},
		{
			name: "FuncType",
			json: `{
					"_Type":	"FuncType",
					"Params":	{
						"_Type":	"FieldList",
						"List":	[{
								"_Type":	"Field",
								"Type":	{
									"_Type":	"Variadic"
								},
								"Doc":	null,
								"Comment":	null,
								"IsStatic":	false,
								"Access":	0,
								"Names":	[]
							}, {
								"_Type":	"Field",
								"Type":	{
									"_Type":	"BuiltinType",
									"Kind":	6,
									"Flags":	0
								},
								"Doc":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"Comment":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"IsStatic":	false,
								"Access":	0,
								"Names":	[{
										"_Type":	"Ident",
										"Name":	"a"
									}]
							}]
					},
					"Ret":	{
						"_Type":	"BuiltinType",
						"Kind":	0,
						"Flags":	0
					}
				}`,
			expected: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type:  &ast.Variadic{},
							Names: []*ast.Ident{},
						}, {
							Doc: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							Comment: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							Type: &ast.BuiltinType{
								Kind:  6,
								Flags: 0,
							},
							Names: []*ast.Ident{
								{Name: "a"},
							},
						},
					},
				},
				Ret: &ast.BuiltinType{
					Kind:  0,
					Flags: 0,
				},
			},
		},
		{
			name: "FuncDecl",
			json: `{
				"_Type":	"FuncDecl",
				"Loc":	{
					"_Type":	"Location",
					"File":	"temp.h"
				},
				"Doc":	{
					"_Type":	"CommentGroup",
					"List":	[]
				},
				"Parent":	null,
				"Name":	{
					"_Type":	"Ident",
					"Name":	"foo"
				},
				"Type":	{
					"_Type":	"FuncType",
					"Params":	{
						"_Type":	"FieldList",
						"List":	[{
								"_Type":	"Field",
								"Type":	{
									"_Type":	"Variadic"
								},
								"Doc":	null,
								"Comment":	null,
								"IsStatic":	false,
								"Access":	0,
								"Names":	[]
							}, {
								"_Type":	"Field",
								"Type":	{
									"_Type":	"BuiltinType",
									"Kind":	6,
									"Flags":	0
								},
								"Doc":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"Comment":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"IsStatic":	false,
								"Access":	0,
								"Names":	[{
										"_Type":	"Ident",
										"Name":	"a"
									}]
							}]
					},
					"Ret":	{
						"_Type":	"BuiltinType",
						"Kind":	0,
						"Flags":	0
					}
				},
				"IsInline":	false,
				"IsStatic":	false,
				"IsConst":	false,
				"IsExplicit":	false,
				"IsConstructor":	false,
				"IsDestructor":	false,
				"IsVirtual":	false,
				"IsOverride":	false
			}`,
			expected: &ast.FuncDecl{
				Object: ast.Object{
					Loc: &ast.Location{
						File: "temp.h",
					},
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{},
					},
					Name: &ast.Ident{Name: "foo"},
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Type:  &ast.Variadic{},
								Names: []*ast.Ident{},
							}, {
								Doc: &ast.CommentGroup{
									List: []*ast.Comment{},
								},
								Comment: &ast.CommentGroup{
									List: []*ast.Comment{},
								},
								Type: &ast.BuiltinType{
									Kind:  6,
									Flags: 0,
								},
								Names: []*ast.Ident{
									{Name: "a"},
								},
							},
						},
					},
					Ret: &ast.BuiltinType{
						Kind:  0,
						Flags: 0,
					},
				},
			},
		},
		{
			name: "RecordType",
			json: `{
					"_Type":	"RecordType",
					"Tag":	3,
					"Fields":	{
						"_Type":	"FieldList",
						"List":	[{
								"_Type":	"Field",
								"Type":	{
									"_Type":	"BuiltinType",
									"Kind":	6,
									"Flags":	0
								},
								"Doc":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"Comment":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"IsStatic":	true,
								"Access":	1,
								"Names":	[{
										"_Type":	"Ident",
										"Name":	"a"
									}]
							}, {
								"_Type":	"Field",
								"Type":	{
									"_Type":	"BuiltinType",
									"Kind":	6,
									"Flags":	0
								},
								"Doc":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"Comment":	{
									"_Type":	"CommentGroup",
									"List":	[]
								},
								"IsStatic":	false,
								"Access":	1,
								"Names":	[{
										"_Type":	"Ident",
										"Name":	"b"
									}]
							}]
					},
					"Methods":	[{
							"_Type":	"FuncDecl",
							"Loc":	{
								"_Type":	"Location",
								"File":	"temp.h"
							},
							"Doc":	{
								"_Type":	"CommentGroup",
								"List":	[]
							},
							"Parent":	{
								"_Type":	"Ident",
								"Name":	"A"
							},
							"Name":	{
								"_Type":	"Ident",
								"Name":	"foo"
							},
							"Type":	{
								"_Type":	"FuncType",
								"Params":	{
									"_Type":	"FieldList",
									"List":	[{
											"_Type":	"Field",
											"Type":	{
												"_Type":	"BuiltinType",
												"Kind":	6,
												"Flags":	0
											},
											"Doc":	{
												"_Type":	"CommentGroup",
												"List":	[]
											},
											"Comment":	{
												"_Type":	"CommentGroup",
												"List":	[]
											},
											"IsStatic":	false,
											"Access":	0,
											"Names":	[{
													"_Type":	"Ident",
													"Name":	"a"
												}]
										}, {
											"_Type":	"Field",
											"Type":	{
												"_Type":	"BuiltinType",
												"Kind":	8,
												"Flags":	16
											},
											"Doc":	{
												"_Type":	"CommentGroup",
												"List":	[]
											},
											"Comment":	{
												"_Type":	"CommentGroup",
												"List":	[]
											},
											"IsStatic":	false,
											"Access":	0,
											"Names":	[{
													"_Type":	"Ident",
													"Name":	"b"
												}]
										}]
								},
								"Ret":	{
									"_Type":	"BuiltinType",
									"Kind":	8,
									"Flags":	0
								}
							},
							"IsInline":	false,
							"IsStatic":	false,
							"IsConst":	false,
							"IsExplicit":	false,
							"IsConstructor":	false,
							"IsDestructor":	false,
							"IsVirtual":	false,
							"IsOverride":	false
						}]
				}`,
			expected: &ast.RecordType{
				Tag: 3,
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.BuiltinType{
								Kind:  6,
								Flags: 0,
							},
							Doc: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							Comment: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							IsStatic: true,
							Access:   1,
							Names: []*ast.Ident{
								{Name: "a"},
							},
						},
						{
							Type: &ast.BuiltinType{
								Kind:  6,
								Flags: 0,
							},
							Doc: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							Comment: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							IsStatic: false,
							Access:   1,
							Names: []*ast.Ident{
								{Name: "b"},
							},
						},
					},
				},
				Methods: []*ast.FuncDecl{
					{
						Object: ast.Object{
							Loc: &ast.Location{
								File: "temp.h",
							},
							Doc: &ast.CommentGroup{
								List: []*ast.Comment{},
							},
							Name:   &ast.Ident{Name: "foo"},
							Parent: &ast.Ident{Name: "A"},
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Type: &ast.BuiltinType{
											Kind:  6,
											Flags: 0,
										},
										Doc: &ast.CommentGroup{
											List: []*ast.Comment{},
										},
										Comment: &ast.CommentGroup{
											List: []*ast.Comment{},
										},
										IsStatic: false,
										Access:   0,
										Names: []*ast.Ident{
											{Name: "a"},
										},
									},
									{
										Type: &ast.BuiltinType{
											Kind:  8,
											Flags: 16,
										},
										Doc: &ast.CommentGroup{
											List: []*ast.Comment{},
										},
										Comment: &ast.CommentGroup{
											List: []*ast.Comment{},
										},
										IsStatic: false,
										Access:   0,
										Names: []*ast.Ident{
											{Name: "b"},
										},
									},
								},
							},
							Ret: &ast.BuiltinType{
								Kind:  8,
								Flags: 0,
							},
						},
						IsInline:      false,
						IsStatic:      false,
						IsConst:       false,
						IsExplicit:    false,
						IsConstructor: false,
						IsDestructor:  false,
						IsVirtual:     false,
						IsOverride:    false,
					},
				},
			},
		},
		{
			name: "TypedefDecl",
			json: `{
				"_Type":	"TypedefDecl",
				"Loc":	{
					"_Type":	"Location",
					"File":	"temp.h"
				},
				"Doc":	{
					"_Type":	"CommentGroup",
					"List":	[]
				},
				"Parent":	null,
				"Name":	{
					"_Type":	"Ident",
					"Name":	"INT"
				},
				"Type":	{
					"_Type":	"BuiltinType",
					"Kind":	6,
					"Flags":	0
				}
			}`,
			expected: &ast.TypedefDecl{
				Object: ast.Object{
					Loc: &ast.Location{
						File: "temp.h",
					},
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{},
					},
					Name: &ast.Ident{
						Name: "INT",
					},
					Parent: nil,
				},
				Type: &ast.BuiltinType{
					Kind:  6,
					Flags: 0,
				},
			},
		},
		{
			name: "EnumItem",
			json: `{
						"_Type":	"EnumItem",
						"Name":	{
							"_Type":	"Ident",
							"Name":	"a"
						},
						"Value":	{
							"_Type":	"BasicLit",
							"Kind":	0,
							"Value":	"0"
						}
					}`,
			expected: &ast.EnumItem{
				Name: &ast.Ident{
					Name: "a",
				},
				Value: &ast.BasicLit{
					Kind:  0,
					Value: "0",
				},
			},
		},
		{
			name: "EnumTypeDecl",
			json: `{
				"_Type":	"EnumTypeDecl",
				"Loc":	{
					"_Type":	"Location",
					"File":	"temp.h"
				},
				"Doc":	{
					"_Type":	"CommentGroup",
					"List":	[]
				},
				"Parent":	null,
				"Name":	{
					"_Type":	"Ident",
					"Name":	"Foo"
				},
				"Type":	{
					"_Type":	"EnumType",
					"Items":	[{
							"_Type":	"EnumItem",
							"Name":	{
								"_Type":	"Ident",
								"Name":	"a"
							},
							"Value":	{
								"_Type":	"BasicLit",
								"Kind":	0,
								"Value":	"0"
							}
						}, {
							"_Type":	"EnumItem",
							"Name":	{
								"_Type":	"Ident",
								"Name":	"b"
							},
							"Value":	{
								"_Type":	"BasicLit",
								"Kind":	0,
								"Value":	"1"
							}
						}, {
							"_Type":	"EnumItem",
							"Name":	{
								"_Type":	"Ident",
								"Name":	"c"
							},
							"Value":	{
								"_Type":	"BasicLit",
								"Kind":	0,
								"Value":	"2"
							}
						}]
				}
			}`,
			expected: &ast.EnumTypeDecl{
				Object: ast.Object{
					Loc: &ast.Location{
						File: "temp.h",
					},
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{},
					},
					Name:   &ast.Ident{Name: "Foo"},
					Parent: nil,
				},
				Type: &ast.EnumType{
					Items: []*ast.EnumItem{
						{
							Name: &ast.Ident{
								Name: "a",
							},
							Value: &ast.BasicLit{
								Kind:  0,
								Value: "0",
							},
						},
						{
							Name: &ast.Ident{
								Name: "b",
							},
							Value: &ast.BasicLit{
								Kind:  0,
								Value: "1",
							},
						},
						{
							Name: &ast.Ident{
								Name: "c",
							},
							Value: &ast.BasicLit{
								Kind:  0,
								Value: "2",
							},
						},
					},
				},
			},
		},
		{
			name: "Macro",
			json: `{
						"_Type":	"Macro",
						"Name":	"OK",
						"Tokens":	[{
								"_Type":	"Token",
								"Token":	3,
								"Lit":	"OK"
							}, {
								"_Type":	"Token",
								"Token":	4,
								"Lit":	"1"
							}]
					}`,
			expected: &ast.Macro{
				Name: "OK",
				Tokens: []*ast.Token{
					{Token: 3, Lit: "OK"},
					{Token: 4, Lit: "1"},
				},
			},
		},
		{
			name: "Include",
			json: `{
				"_Type":	"Include",
				"Path":	"foo.h"
			}`,
			expected: &ast.Include{
				Path: "foo.h",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := unmarshal.Node([]byte(tc.json))
			if err != nil {
				t.Fatalf("UnmarshalNode failed: %v", err)
			}

			resultJSON, err := json.MarshalIndent(node, "", " ")
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			expectedJSON, err := json.MarshalIndent(tc.expected, "", " ")

			if err != nil {
				t.Fatalf("Failed to marshal expected result to JSON: %v", err)
			}

			if string(resultJSON) != string(expectedJSON) {
				t.Errorf("JSON mismatch.\nExpected: %s\nGot: %s", string(expectedJSON), string(resultJSON))
			}
		})
	}
}

func TestUnmarshalErrors(t *testing.T) {
	testCases := []struct {
		name        string
		fn          func([]byte) (ast.Node, error)
		input       string
		expectedErr string
	}{
		// UnmarshalNode errors
		{
			name:        "UnmarshalNode - Invalid JSON",
			fn:          unmarshal.Node,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Node into unmarshal.nodeTemp",
		},
		{
			name:        "UnmarshalNode - Unknown type",
			fn:          unmarshal.Node,
			input:       `{"_Type": "UnknownType"}`,
			expectedErr: "unknown node type: UnknownType",
		},

		// unmarshalToken errors
		{
			name:        "unmarshalToken - Invalid JSON",
			fn:          unmarshal.Token,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Token into ast.Token",
		},

		// unmarshalMacro errors
		{
			name:        "unmarshalMacro - Invalid JSON",
			fn:          unmarshal.Macro,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Macro into ast.Macro",
		},

		// unmarshalInclude errors
		{
			name:        "unmarshalInclude - Invalid JSON",
			fn:          unmarshal.Include,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Include into ast.Include",
		},

		// unmarshalBasicLit errors
		{
			name:        "unmarshalBasicLit - Invalid JSON",
			fn:          unmarshal.BasicLit,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in BasicLit into ast.BasicLit",
		},

		// unmarshalBuiltinType errors
		{
			name:        "unmarshalBuiltinType - Invalid JSON",
			fn:          unmarshal.BuiltinType,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in BuiltinType into ast.BuiltinType",
		},

		// unmarshalIdent errors
		{
			name:        "unmarshalIdent - Invalid JSON",
			fn:          unmarshal.Ident,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Ident into ast.Ident",
		},

		// unmarshalVariadic errors
		{
			name:        "unmarshalVariadic - Invalid JSON",
			fn:          unmarshal.Variadic,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Variadic into ast.Variadic",
		},

		// unmarshalXType errors
		{
			name:        "unmarshalXType - Invalid JSON",
			fn:          unmarshal.PointerType,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in XType into unmarshal.XTypeTemp",
		},
		{
			name:        "unmarshalXType - Invalid X field",
			fn:          unmarshal.PointerType,
			input:       `{"X": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in XType when converting X of unmarshal.XTypeTemp",
		},
		{
			name: "unmarshalXType - Unexpected type",
			fn: func(data []byte) (ast.Node, error) {
				return unmarshal.XType(data, &ast.BasicLit{})
			},
			input:       `{"X": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in XType: got *ast.Token, want ast.Expr",
		},
		{
			name: "unmarshalXType - Unexpected type",
			fn: func(data []byte) (ast.Node, error) {
				return unmarshal.XType(data, &ast.BasicLit{})
			},
			input:       `{"X": {"_Type": "Ident", "Name": "test"}}`,
			expectedErr: "unmarshal error in XType: got *ast.BasicLit, want *ast.PointerType, *ast.LvalueRefType, *ast.RvalueRefType",
		},
		// unmarshalArrayType errors
		{
			name:        "unmarshalArrayType - Invalid JSON",
			fn:          unmarshal.ArrayType,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in ArrayType into unmarshal.arrayTemp",
		},
		{
			name:        "unmarshalArrayType - Invalid Elt",
			fn:          unmarshal.ArrayType,
			input:       `{"Elt": {"_Type": "InvalidType"}, "Len": null}`,
			expectedErr: "unmarshal error in ArrayType when converting Elt of unmarshal.arrayTemp",
		},
		{
			name:        "unmarshalArrayType - Unexpect Elt",
			fn:          unmarshal.ArrayType,
			input:       `{"Elt": {"_Type": "Token", "Token": 1, "Lit": "test"}, "Len": null}`,
			expectedErr: "unmarshal error in ArrayType: got *ast.Token, want ast.Expr",
		},
		{
			name:        "unmarshalArrayType - Invalid Len",
			fn:          unmarshal.ArrayType,
			input:       `{"Elt": {"_Type": "BuiltinType", "Kind": 1}, "Len": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in ArrayType when converting Len of *ast.ArrayType",
		},
		{
			name:        "unmarshalArrayType - Unexpect Len",
			fn:          unmarshal.ArrayType,
			input:       `{"Elt": {"_Type": "BuiltinType", "Kind": 1}, "Len": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in ArrayType: got *ast.Token, want ast.Expr",
		},

		// unmarshalField errors
		{
			name:        "unmarshalField - Invalid JSON",
			fn:          unmarshal.Field,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in Field into unmarshal.fieldTemp",
		},
		{
			name:        "unmarshalField - Invalid Type",
			fn:          unmarshal.Field,
			input:       `{"Type": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in Field when converting Type of unmarshal.fieldTemp",
		},

		// unmarshalFieldList errors
		{
			name:        "unmarshalFieldList - Invalid JSON",
			fn:          unmarshal.FieldList,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in FieldList into unmarshal.fieldListTemp",
		},
		{
			name:        "unmarshalFieldList - Invalid field",
			fn:          unmarshal.FieldList,
			input:       `{"List": [{"_Type": "InvalidType"}]}`,
			expectedErr: "unmarshal error in FieldList when converting List of unmarshal.fieldListTemp",
		},
		{
			name:        "unmarshalFieldList - Unexpected field",
			fn:          unmarshal.FieldList,
			input:       `{"List": [{"_Type": "Token", "Token": 1, "Lit": "test"}]}`,
			expectedErr: "unmarshal error in FieldList: got *ast.Token, want *ast.Field",
		},
		// unmarshalTagExpr errors
		{
			name:        "unmarshalTagExpr - Invalid JSON",
			fn:          unmarshal.TagExpr,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in TagExpr into unmarshal.tagExprTemp",
		},
		{
			name:        "unmarshalTagExpr - Invalid Name",
			fn:          unmarshal.TagExpr,
			input:       `{"Name": {"_Type": "InvalidType"}, "Tag": 0}`,
			expectedErr: "unmarshal error in TagExpr when converting Name of unmarshal.tagExprTemp",
		},
		{
			name:        "unmarshalTagExpr - Unexpected Name",
			fn:          unmarshal.TagExpr,
			input:       `{"Name": {"_Type": "Token", "Token": 1, "Lit": "test"}, "Tag": 0}`,
			expectedErr: "unmarshal error in TagExpr: got *ast.Token, want ast.Expr",
		},

		// unmarshalScopingExpr errors
		{
			name:        "unmarshalScopingExpr - Invalid JSON",
			fn:          unmarshal.ScopingExpr,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in ScopingExpr into unmarshal.scopingExprTemp",
		},
		{
			name:        "unmarshalScopingExpr - Invalid Parent",
			fn:          unmarshal.ScopingExpr,
			input:       `{"Parent": {"_Type": "InvalidType"}, "X": {"_Type": "Ident", "Name": "test"}}`,
			expectedErr: "unmarshal error in ScopingExpr when converting Parent of unmarshal.scopingExprTemp",
		},
		{
			name:        "unmarshalScopingExpr - Unexpected Parent",
			fn:          unmarshal.ScopingExpr,
			input:       `{"Parent": {"_Type": "Token", "Token": 1, "Lit": "test"}, "X": {"_Type": "Ident", "Name": "test"}}`,
			expectedErr: "unmarshal error in ScopingExpr: got *ast.Token, want ast.Expr",
		},
		{
			name:        "unmarshalScopingExpr - Invalid X",
			fn:          unmarshal.ScopingExpr,
			input:       `{"Parent": {"_Type": "Ident", "Name": "test"}, "X": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in ScopingExpr when converting X of unmarshal.scopingExprTemp",
		},
		{
			name:        "unmarshalScopingExpr - Unexpected X",
			fn:          unmarshal.ScopingExpr,
			input:       `{"Parent": {"_Type": "Ident", "Name": "test"}, "X": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in ScopingExpr: got *ast.Token, want ast.Ident",
		},

		// unmarshalEnumItem errors
		{
			name:        "unmarshalEnumItem - Invalid JSON",
			fn:          unmarshal.EnumItem,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in EnumItem into unmarshal.enumItemTemp",
		},
		{
			name:        "unmarshalEnumItem - Invalid Value",
			fn:          unmarshal.EnumItem,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Value": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in EnumItem when converting Value of unmarshal.enumItemTemp",
		},
		{
			name:        "unmarshalEnumItem - Unexpected Value",
			fn:          unmarshal.EnumItem,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Value": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in EnumItem: got *ast.Token, want ast.Expr",
		},

		// unmarshalEnumType errors
		{
			name:        "unmarshalEnumType - Invalid JSON",
			fn:          unmarshal.EnumType,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in EnumType into unmarshal.enumTypeTemp",
		},
		{
			name:        "unmarshalEnumType - Invalid Item",
			fn:          unmarshal.EnumType,
			input:       `{"Items": [{"_Type": "InvalidType"}]}`,
			expectedErr: "unmarshal error in EnumType when converting Items of unmarshal.enumTypeTemp",
		},
		{
			name:        "unmarshalEnumType - Unexpected Item",
			fn:          unmarshal.EnumType,
			input:       `{"Items": [{"_Type": "Token", "Token": 1, "Lit": "test"}]}`,
			expectedErr: "unmarshal error in EnumType: got *ast.Token, want *ast.EnumItem",
		},

		// unmarshalRecordType errors
		{
			name:        "unmarshalRecordType - Invalid JSON",
			fn:          unmarshal.RecordType,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in RecordType into unmarshal.recordTypeTemp",
		},
		{
			name:        "unmarshalRecordType - Invalid Fields",
			fn:          unmarshal.RecordType,
			input:       `{"Tag": 0, "Fields": {"_Type": "InvalidType"}, "Methods": []}`,
			expectedErr: "unmarshal error in RecordType when converting Fields of unmarshal.recordTypeTemp",
		},
		{
			name:        "unmarshalRecordType - Unexpected Fields",
			fn:          unmarshal.RecordType,
			input:       `{"Tag": 0, "Fields": {"_Type": "Token", "Token": 1, "Lit": "test"}, "Methods": []}`,
			expectedErr: "unmarshal error in RecordType: got *ast.Token, want *ast.FieldList",
		},
		{
			name:        "unmarshalRecordType - Invalid Method",
			fn:          unmarshal.RecordType,
			input:       `{"Tag": 0, "Fields": {"_Type": "FieldList", "List": []}, "Methods": [{"_Type": "InvalidType"}]}`,
			expectedErr: "unmarshal error in RecordType when converting Methods of unmarshal.recordTypeTemp",
		},
		{
			name:        "unmarshalRecordType - Unexpected Method",
			fn:          unmarshal.RecordType,
			input:       `{"Tag": 0, "Fields": {"_Type": "FieldList", "List": []}, "Methods": [{"_Type": "Token", "Token": 1, "Lit": "test"}]}`,
			expectedErr: "unmarshal error in RecordType: got *ast.Token, want *ast.FuncDecl",
		},

		// unmarshalFuncType errors
		{
			name:        "unmarshalFuncType - Invalid JSON",
			fn:          unmarshal.FuncType,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in FuncType into unmarshal.funcTypeTemp",
		},
		{
			name:        "unmarshalFuncType - Invalid Params",
			fn:          unmarshal.FuncType,
			input:       `{"Params": {"_Type": "InvalidType"}, "Ret": {"_Type": "BuiltinType", "Kind": 1}}`,
			expectedErr: "unmarshal error in FuncType when converting Params of unmarshal.funcTypeTemp",
		},
		{
			name:        "unmarshalFuncType - Unexpected Params",
			fn:          unmarshal.FuncType,
			input:       `{"Params": {"_Type": "Token", "Token": 1, "Lit": "test"}, "Ret": {"_Type": "BuiltinType", "Kind": 1}}`,
			expectedErr: "unmarshal error in FuncType: got *ast.Token, want *ast.FieldList",
		},
		{
			name:        "unmarshalFuncType - Invalid Ret",
			fn:          unmarshal.FuncType,
			input:       `{"Params": {"_Type": "FieldList", "List": []}, "Ret": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in FuncType when converting Ret of unmarshal.funcTypeTemp",
		},
		{
			name:        "unmarshalFuncType - Unexpected Ret",
			fn:          unmarshal.FuncType,
			input:       `{"Params": {"_Type": "FieldList", "List": []}, "Ret": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in FuncType: got *ast.Token, want ast.Expr",
		},

		// unmarshalFuncDecl errors
		{
			name:        "unmarshalFuncDecl - Invalid JSON",
			fn:          unmarshal.FuncDecl,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in FuncDecl into unmarshal.funcDeclTemp",
		},
		{
			name:        "unmarshalFuncDecl - Invalid Type",
			fn:          unmarshal.FuncDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in FuncDecl when converting Type of unmarshal.funcDeclTemp",
		},
		{
			name:        "unmarshalFuncDecl - Unexpected Type",
			fn:          unmarshal.FuncDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in FuncDecl: got *ast.Token, want *ast.FuncType",
		},
		{
			name:        "unmarshalFuncDecl - Invalid DeclBase",
			fn:          unmarshal.FuncDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "FuncType", "Params": {"_Type": "FieldList", "List": []}, "Ret": {"_Type": "BuiltinType", "Kind": 1}}, "Loc": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in declBase when converting Parent of unmarshal.declBaseTemp",
		},
		{
			name:        "unmarshalFuncDecl - Invalid DeclBase",
			fn:          unmarshal.FuncDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"},"Loc":false, "Type": {"_Type": "FuncType", "Params": {"_Type": "FieldList", "List": []}, "Ret": {"_Type": "BuiltinType", "Kind": 1}}, "Loc": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in declBase into unmarshal.declBaseTemp",
		},
		{
			name:        "unmarshalFuncDecl - Unexpected DeclBase",
			fn:          unmarshal.FuncDecl,
			input:       `{"Loc":{"_Type":"Location","File":"temp.h"},"Parent":{"_Type":"Token","Token":1,"Lit":"test"},"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "FuncType", "Params": {"_Type": "FieldList", "List": []}, "Ret": {"_Type": "BuiltinType", "Kind": 1}}, "Loc": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in declBase: got *ast.Token, want ast.Expr",
		},
		// unmarshalTypeDecl errors
		{
			name:        "unmarshalTypeDecl - Invalid JSON",
			fn:          unmarshal.TypeDecl,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in TypeDecl into unmarshal.typeDeclTemp",
		},
		{
			name:        "unmarshalTypeDecl - Invalid Type",
			fn:          unmarshal.TypeDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in TypeDecl when converting Type of unmarshal.typeDeclTemp",
		},
		{
			name:        "unmarshalTypeDecl - Unexpected Type",
			fn:          unmarshal.TypeDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in TypeDecl: got *ast.Token, want *ast.RecordType",
		},
		{
			name:        "unmarshalTypeDecl - Invalid DeclBase",
			fn:          unmarshal.TypeDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "RecordType", "Tag": 0, "Fields": {"_Type": "FieldList", "List": []}, "Methods": []}, "Loc": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in declBase when converting Parent of unmarshal.declBaseTemp",
		},
		// unmarshalTypeDefDecl errors
		{
			name:        "unmarshalTypeDefDecl - Invalid JSON",
			fn:          unmarshal.TypeDefDecl,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in TypeDefDecl into unmarshal.typeDefDeclTemp",
		},
		{
			name:        "unmarshalTypeDefDecl - Invalid Type",
			fn:          unmarshal.TypeDefDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in TypeDefDecl when converting Type of unmarshal.typeDefDeclTemp",
		},
		{
			name:        "unmarshalTypeDefDecl - Unexpected Type",
			fn:          unmarshal.TypeDefDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in TypeDefDecl: got *ast.Token, want ast.Expr",
		},
		{
			name:        "unmarshalTypeDefDecl - Invalid DeclBase",
			fn:          unmarshal.TypeDefDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "BuiltinType", "Kind": 1}, "Loc": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in declBase when converting Parent of unmarshal.declBaseTemp",
		},
		// unmarshalEnumTypeDecl errors
		{
			name:        "unmarshalEnumTypeDecl - Invalid JSON",
			fn:          unmarshal.EnumTypeDecl,
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in EnumTypeDecl into unmarshal.enumTypeDeclTemp",
		},
		{
			name:        "unmarshalEnumTypeDecl - Invalid Type",
			fn:          unmarshal.EnumTypeDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in EnumTypeDecl when converting Type of unmarshal.enumTypeDeclTemp",
		},
		{
			name:        "unmarshalEnumTypeDecl - Unexpected Type",
			fn:          unmarshal.EnumTypeDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "Token", "Token": 1, "Lit": "test"}}`,
			expectedErr: "unmarshal error in EnumTypeDecl: got *ast.Token, want *ast.EnumType",
		},
		{
			name:        "unmarshalEnumTypeDecl - Invalid DeclBase",
			fn:          unmarshal.EnumTypeDecl,
			input:       `{"Name": {"_Type": "Ident", "Name": "test"}, "Type": {"_Type": "EnumType", "Items": []}, "Loc": {"_Type": "InvalidType"}}`,
			expectedErr: "unmarshal error in declBase when converting Parent of unmarshal.declBaseTemp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.fn([]byte(tc.input))
			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				} else {
					return
				}
			}
			if err == nil {
				t.Errorf("Expected error, but got nil")
			} else if !strings.Contains(err.Error(), tc.expectedErr) {
				t.Errorf("Expected error containing %q, but got %q", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestUnmarshalFileErrors(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr string
	}{
		// unmarshalFile errors
		{
			name:        "unmarshalFile - Invalid JSON",
			input:       `{"invalid": "json"`,
			expectedErr: "unmarshal error in File into unmarshal.fileTemp",
		},
		{
			name:        "unmarshalFile - Invalid JSON",
			input:       `{"_Type": "File", "Includes": [], "Macros": [], "Decls": [{"_Type": "InvalidType"}]}`,
			expectedErr: "",
		},
		{
			name:        "unmarshalFile - Unexpected Decl",
			input:       `{"_Type": "File", "Includes": [], "Macros": [], "Decls": [{"_Type": "Token", "Token": 1, "Lit": "test"}]}`,
			expectedErr: "unmarshal error in File: got *ast.Token, want ast.Decl",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := unmarshal.File([]byte(tc.input))
			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing %q, but got nil", tc.expectedErr)
				} else if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("Expected error containing %q, but got: %v", tc.expectedErr, err)
				}
			}
		})
	}
}
