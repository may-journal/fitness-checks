package main

import (
	"go/ast"
	"go/token"
)

// testingNames is the set of names a function parameter of type *testing.T
// has anywhere in f.
func testingNames(f *ast.File) map[string]bool {
	names := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		if field, ok := n.(*ast.Field); ok && isTestingT(field.Type) {
			for _, id := range field.Names {
				names[id.Name] = true
			}
		}
		return true
	})
	return names
}

// isTestingT reports whether e is the type *testing.T.
func isTestingT(e ast.Expr) bool {
	star, ok := e.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "testing" && sel.Sel.Name == "T"
}

// isRun reports whether call is t.Run("literal", f) on one of the testing
// names.
func isRun(call *ast.CallExpr, names map[string]bool) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Run" || len(call.Args) != 2 {
		return false
	}
	return onTesting(sel.X, names) && isStringLit(call.Args[0])
}

// onTesting reports whether e is an identifier among the testing names.
func onTesting(e ast.Expr, names map[string]bool) bool {
	id, ok := e.(*ast.Ident)
	return ok && names[id.Name]
}

// isStringLit reports whether e is a string literal.
func isStringLit(e ast.Expr) bool {
	lit, ok := e.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING
}
