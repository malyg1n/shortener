package analyze

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var OsExitCheckCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit()",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if checkFile(file, "main") {
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.FuncDecl:
					if x.Name.String() == "main" {
						findInNode(node, pass, "os", "Exit")
					}
					return true
				default:
					// pass
				}
				return true
			})
		}
	}
	return nil, nil
}

func checkFile(file *ast.File, name string) bool {
	return file.Name.String() == name
}

func findInNode(node ast.Node, pass *analysis.Pass, pck, fnc string) bool {
	ast.Inspect(node, func(nodeInMain ast.Node) bool {
		switch xt := nodeInMain.(type) {
		case *ast.CallExpr:
			if cf, ok := xt.Fun.(*ast.SelectorExpr); ok {
				if fmt.Sprintf("%v", cf.X) == pck {
					if fmt.Sprintf("%v", cf.Sel) == fnc {
						pass.Reportf(xt.Pos(), fmt.Sprintf("%s.%s should be removed", pck, fnc))
						return true
					}
				}
			}
		default:
			// pass
		}
		return true
	})

	return true
}
