package main

import (
	"fmt"
	"os"

	"go/parser"
	"go/ast"
	"go/token"
)

func makeUniverse() *ast.Scope {
	return nil
}

func main2() {
	fset := token.NewFileSet()
	path := "/Users/h8liu/gopath/src/lonnie.io/e8vm/asm8"
	pkgs, e := parser.ParseDir(fset, path, nil, 0)
	if e != nil {
		fmt.Fprintf(os.Stderr, "parse: %s", e)
		return
	}

	// universe := makeUniverse()

	for name, p := range pkgs {
		fmt.Println(name)
		// resolve symbols
		p, e = ast.NewPackage(fset, p.Files, nil, nil)
		if e != nil {
			fmt.Fprintf(os.Stderr, "NewPackage: %s", e)
			// return
		}

		for name, f := range p.Files {
			fmt.Println(name)

			for _, unr := range f.Unresolved {
				fmt.Println("- ", unr)
			}
		}
	}
}
