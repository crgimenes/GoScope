package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type FunctionCall struct {
	Caller string // function that performs the call
	Callee string // called function
	File   string // file where the call occurs
	Line   int    // line number of the call
}

var (
	callGraph            = make(map[string][]FunctionCall)
	functions            = make(map[string]string)
	variableDeclarations = make(map[string]string)
	constantDeclarations = make(map[string]string)
	typeDeclarations     = make(map[string]string)
)

func main() {
	dir := "."
	flag.Parse()
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	fset := token.NewFileSet()

	err := filepath.Walk(dir, func(
		path string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", path, err)
			return nil
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				pos := fset.Position(node.Pos())
				functions[node.Name.Name] = fmt.Sprintf("%s:%d",
					pos.Filename,
					pos.Line,
				)

				if node.Body != nil {
					ast.Inspect(node.Body, func(n2 ast.Node) bool {
						callExpr, ok := n2.(*ast.CallExpr)
						if ok {
							var calleeName string
							switch fun := callExpr.Fun.(type) {
							case *ast.Ident:
								calleeName = fun.Name
							case *ast.SelectorExpr:
								calleeName = fun.Sel.Name
							}
							if calleeName != "" {
								posCall := fset.Position(callExpr.Pos())
								fc := FunctionCall{
									Caller: node.Name.Name,
									Callee: calleeName,
									File:   posCall.Filename,
									Line:   posCall.Line,
								}
								callGraph[calleeName] = append(
									callGraph[calleeName],
									fc,
								)
							}
						}
						return true
					})
				}
			case *ast.GenDecl:
				for _, spec := range node.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						for _, name := range s.Names {
							pos := fset.Position(name.Pos())
							location := fmt.Sprintf("%s:%d",
								pos.Filename,
								pos.Line,
							)
							if node.Tok == token.VAR {
								variableDeclarations[name.Name] = location
								continue
							}

							if node.Tok == token.CONST {
								constantDeclarations[name.Name] = location
								continue
							}
						}
					case *ast.TypeSpec:
						// Handles type declarations.
						pos := fset.Position(s.Pos())
						location := fmt.Sprintf("%s:%d",
							pos.Filename,
							pos.Line,
						)
						typeDeclarations[s.Name.Name] = location
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}

	for name, loc := range functions {
		fmt.Printf("%s %s\n", name, loc)
	}

	for name, loc := range variableDeclarations {
		fmt.Printf("%s %s\n", name, loc)
	}

	for name, loc := range constantDeclarations {
		fmt.Printf("%s %s\n", name, loc)
	}

	for name, loc := range typeDeclarations {
		fmt.Printf("%s %s\n", name, loc)
	}

	for _, callers := range callGraph {
		for _, call := range callers {
			fmt.Printf("%s.%s %s:%d\n",
				call.Caller,
				call.Callee,
				call.File,
				call.Line,
			)
		}
	}
}
