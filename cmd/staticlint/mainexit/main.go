// Package mainexit provides an analysis tool to check for the direct invocation
// of os.Exit in the main function of a Go program.
package mainexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is an analysis.Analyzer that checks for the direct invocation of os.Exit in the main function.
var Analyzer = &analysis.Analyzer{
	Name: "exitInMainCheck",
	Doc:  "check invoke os.Exit function in main function",
	Run:  run,
}

// run is the main function that gets executed for the analysis. It iterates over all the files in the package,
// and if the package is named "main", it looks for the main function to inspect its body for os.Exit calls.
func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if fn, isFn := decl.(*ast.FuncDecl); isFn && fn.Name.Name == "main" && fn.Recv == nil {
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					switch stmt := n.(type) {
					case *ast.GoStmt:
						checkExpr(pass, stmt.Call)
					case *ast.DeferStmt:
						checkExpr(pass, stmt.Call)
					case *ast.CallExpr:
						checkExpr(pass, stmt)
					}
					return true
				})
			}
		}
	}
	return nil, nil
}

// checkExpr inspects a CallExpr node to see if it invokes os.Exit, and if so, it reports it.
func checkExpr(pass *analysis.Pass, call *ast.CallExpr) {
	if fun, isSelector := call.Fun.(*ast.SelectorExpr); isSelector {
		if ident, isIdent := fun.X.(*ast.Ident); isIdent && ident.Name == "os" && fun.Sel.Name == "Exit" {
			pass.Reportf(call.Pos(), "direct use of os.Exit in main function is prohibited")
		}
	}
}
