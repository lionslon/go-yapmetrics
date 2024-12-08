package noosexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noosexit",
		Doc:  "check for os.Exit in main.main",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() == "main" {
		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if funcDecl, ok := node.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == "main" {
						reportOsExit(pass, funcDecl)
					}
				}
				return true
			})
		}
	}

	return nil, nil
}

func reportOsExit(pass *analysis.Pass, mainFuncDeclNode ast.Node) {
	ast.Inspect(mainFuncDeclNode, func(node ast.Node) bool {
		if selectorExpr, ok := node.(*ast.SelectorExpr); ok {
			if selIdent, ok := selectorExpr.X.(*ast.Ident); ok {
				if selIdent.Name == "os" && selectorExpr.Sel.Name == "Exit" {
					pass.Reportf(selectorExpr.Pos(), "cannot use os.Exit in main function of package main")
				}
			}
		}

		return true
	})
}
