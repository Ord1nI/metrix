package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ErrCheckAnalyzer = &analysis.Analyzer{
    Name: "exitcheck",
    Doc:  "check for os.Exit in main function",
    Run:  run,
}


func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		if f.Name.Name == "main" {

			ast.Inspect(f, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl);ok && fd.Name.Name == "main"{
					for _, stm := range fd.Body.List {
						if es, ok := stm.(*ast.ExprStmt); ok {
							if fc, ok := es.X.(*ast.CallExpr); ok {
								if se, ok := fc.Fun.(*ast.SelectorExpr); ok {
									if i, ok := se.X.(*ast.Ident);ok && se.Sel.Name == "Exit" && i.Name == "os" {
										d := analysis.Diagnostic{
											Pos: f.Pos(),
											End: f.End(),
											Message: "Cant use os.Exit in main function.",
										}
										pass.Report(d)
									}
								}
							}
						}
					}
				}
				return true
			})
		}
		break
	}
	return nil, nil
}
