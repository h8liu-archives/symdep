package main

import (
	"fmt"
	"os"

	"go/build"

	"golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/refactor/lexical"
)

func main() {
	ctxt := build.Default

	saved := ctxt.GOPATH
	ctxt.GOPATH = ""
	pkgs := buildutil.AllPackages(&ctxt)
	ctxt.GOPATH = saved

	pkgs = append(pkgs, "lonnie.io/e8vm/asm8")

	conf := loader.Config{
		Build: &ctxt,
		SourceImports: true,
	}
	for _, path := range pkgs {
		if err := conf.ImportWithTests(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	iprog, err := conf.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	for pkg, info := range iprog.AllPackages {
		fmt.Println(pkg.Name())
		/* for _, f := range info.Files {
			// fmt.Println(name)
			for _, unr := range f.Unresolved {
				fmt.Println("- ", unr)
			}
		} */

		lexical.Structure(iprog.Fset, pkg, &info.Info, info.Files)
	}
}
