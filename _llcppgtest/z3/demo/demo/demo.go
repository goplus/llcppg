package main

import (
	"z3"

	"github.com/goplus/lib/c"
)

func main() {
	cfg := z3.MkConfig()
	ctx := z3.MkContext(cfg)
	solver := z3.MkSolver(ctx)
	z3.SolverIncRef(ctx, solver)
	z3.DelConfig(cfg)

	int_sort := z3.MkIntSort(ctx)
	x := z3.MkConst(ctx, z3.MkStringSymbol(ctx, c.Str("x")), int_sort)

	ten := z3.MkInt(ctx, 10, int_sort)
	constraint := z3.MkGt(ctx, x, ten)

	z3.SolverAssert(ctx, solver, constraint)

	if z3.SolverCheck(ctx, solver) == z3.LTRUE {
		model := z3.SolverGetModel(ctx, solver)
		z3.ModelIncRef(ctx, model)

		var val z3.Ast
		z3.ModelEval(ctx, model, x, true, &val)
		c.Printf(c.Str("find solution: x = %s\n"), z3.AstToString(ctx, val))

		z3.ModelDecRef(ctx, model)
	}

	z3.SolverDecRef(ctx, solver)
	z3.DelContext(ctx)
}
