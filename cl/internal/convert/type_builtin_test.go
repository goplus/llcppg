package convert

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/cltest"
	"github.com/goplus/mod/gopmod"
)

func TestIdentRef(t *testing.T) {
	ct := &TypeConv{}
	_, err := ct.handleIdentRefer(&ast.BuiltinType{Kind: ast.Bool})
	if err == nil {
		t.Fatal("Expect Error")
	}
}

func TestSubstObj(t *testing.T) {
	mod, err := gopmod.Load(".")
	if err != nil {
		panic(err)
	}
	pkg, err := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath: ".",
		},
		Name:       "testpkg",
		GenConf:    &gogen.Config{},
		ConvSym:    cltest.NewConvSym(),
		LibCommand: "${pkg-config --libs xxx}",
		Mod:        mod,
	})
	if err != nil {
		t.Fatal("NewPackage failed")
	}

	corg := types.NewNamed(types.NewTypeName(token.NoPos, nil, "origin", nil), types.Typ[types.Int], nil)
	corg2 := types.NewNamed(types.NewTypeName(token.NoPos, nil, "origin2", nil), types.Typ[types.Int], nil)
	substObj(pkg.p.Types, pkg.p.Types.Scope(), "GoPub", corg.Obj())
	name := gogen.Lookup(pkg.p.Types.Scope(), "GoPub")
	if name == nil {
		t.Fatal("Lookup failed")
	}
	if name.Type().String() != corg.String() {
		t.Fatal("Type not equal")
	}

	// reassign
	substObj(pkg.p.Types, pkg.p.Types.Scope(), "GoPub", corg2.Obj())
	name2 := gogen.Lookup(pkg.p.Types.Scope(), "GoPub")
	if name2 == nil {
		t.Fatal("Lookup failed")
	}
	if name2.Type().String() != corg2.String() {
		t.Fatal("Type not equal")
	}
}
