package main

import (
	"fmt"
	"os"
	"sort"

	"go/build"
	"path/filepath"

	// "golang.org/x/tools/go/buildutil"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/refactor/lexical"
)

func main() {
	target := "lonnie.io/e8vm/arch8"

	ctxt := build.Default

	pkgs := []string{target}

	conf := loader.Config{
		Build:         &ctxt,
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

	var deps []string
	fileDeps := make(map[string]bool)

	for pkg, info := range iprog.AllPackages {
		name := pkg.Path()
		if name != target {
			deps = append(deps, name)
			// fmt.Println(name) // depending packages
			continue
		}

		res := lexical.Structure(iprog.Fset, pkg, &info.Info, info.Files)
		fset := iprog.Fset

		for obj, refs := range res.Refs {
			pack := obj.Pkg()
			if pack != pkg {
				continue
			}

			fdef := filepath.Base(fset.Position(obj.Pos()).Filename)

			for _, ref := range refs {
				fused := filepath.Base(fset.Position(ref.Id.NamePos).Filename)
				if fused != fdef {
					s := fdef + " <- " + fused
					fileDeps[s] = true
				}
			}
		}
	}

	var strs []string
	for m := range fileDeps {
		strs = append(strs, m)
	}
	sort.Strings(deps)
	sort.Strings(strs)

	for _, dep := range deps {
		fmt.Println(dep)
	}
	for _, s := range strs {
		fmt.Println(s)
	}
}
