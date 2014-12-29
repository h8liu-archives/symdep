package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	// "sort"

	"go/build"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
	"golang.org/x/tools/refactor/lexical"
)

// File contains dependency info of a file
type File struct {
	Name    string
	Depends map[string]*File
	UsedBy  map[string]*File

	hits  map[string]bool
	level int
}

// Package contains dependency info of a package.
type Package struct {
	Depends []string         // depending packages
	Files   map[string]*File // files

	LevelsTop [][]*File
	LevelsBot [][]*File
}

type reverse bool

func (r reverse) inputs(f *File) map[string]*File {
	if !r {
		return f.Depends
	}
	return f.UsedBy
}

func (r reverse) outputs(f *File) map[string]*File {
	if !r {
		return f.UsedBy
	}
	return f.Depends
}

func newFile(name string) *File {
	ret := new(File)
	ret.Name = name
	ret.Depends = make(map[string]*File)
	ret.UsedBy = make(map[string]*File)

	return ret
}

func (p *Package) addFile(s string) {
	p.Files[s] = newFile(s)
}

func (p *Package) depends(fdef, fused string) {
	def := p.Files[fdef]
	use := p.Files[fused]
	if def == nil || use == nil {
		panic("file not found")
	}

	if def.UsedBy[fused] == nil {
		def.UsedBy[fused] = use
	}
	if use.Depends[fdef] == nil {
		use.Depends[fdef] = def
	}
}

// PrintDepends prints the list of dependencies.
func (p *Package) PrintDepends() {
	for _, dep := range p.Depends {
		fmt.Println(dep)
	}
}

func (p *Package) PrintFileDeps() {
	for _, f := range p.Files {
		fmt.Println(f.Name)

		for _, dep := range f.Depends {
			fmt.Println("    ", dep.Name)
		}
	}
}

func printLevels(lvls [][]*File) {
	for i, lvl := range lvls {
		fmt.Printf("%d:", i)
		for _, f := range lvl {
			fmt.Printf(" %s", f.Name)
		}
		fmt.Println()
	}
}

func (p *Package) PrintLevelsTop() { printLevels(p.LevelsTop) }
func (p *Package) PrintLevelsBot() {
	n := len(p.LevelsBot)
	lvls := make([][]*File, n)
	for i, lvl := range p.LevelsBot {
		lvls[n-1-i] = lvl
	}
	printLevels(lvls)
}

var errCircDep = errors.New("has circular dependency")

func (p *Package) makeLevels(r reverse) ([][]*File, error) {
	var cur []*File
	var next []*File
	var ret [][]*File

	for _, f := range p.Files {
		if len(r.inputs(f)) == 0 {
			cur = append(cur, f)
			f.level = 0
		} else {
			f.hits = make(map[string]bool)
			f.level = -1
		}
	}

	ntotal := 0
	level := 1
	for len(cur) > 0 {
		ret = append(ret, cur)
		ntotal += len(cur)

		for _, f := range cur {
			outputs := r.outputs(f)
			for _, out := range outputs {
				if r.inputs(out)[f.Name] == nil {
					panic("bug")
				}

				if out.level >= 0 && out.level < level {
					return nil, errCircDep
				}

				wasHit := out.hits[f.Name]
				out.hits[f.Name] = true

				if !wasHit && len(out.hits) == len(r.inputs(out)) {
					next = append(next, out)
					out.level = level
				}
			}
		}

		cur = next
		next = nil
		level++
	}

	if ntotal != len(p.Files) {
		return nil, errCircDep
	}

	return ret, nil
}

func (p *Package) MakeLevelsTop() error {
	var e error
	p.LevelsTop, e = p.makeLevels(reverse(false))
	return e
}

func (p *Package) MakeLevelsBot() error {
	var e error
	p.LevelsBot, e = p.makeLevels(reverse(true))
	return e
}

// NewPackage analyzes the dependency info of a package.
func NewPackage(packPath string) (*Package, error) {
	ctxt := build.Default
	conf := loader.Config{
		Build:         &ctxt,
		SourceImports: true,
	}

	if err := conf.ImportWithTests(packPath); err != nil {
		return nil, err
	}

	iprog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	ret := new(Package)
	ret.Files = make(map[string]*File)

	var thePkg *types.Package
	for pkg := range iprog.AllPackages {
		name := pkg.Path()
		if name != packPath {
			ret.Depends = append(ret.Depends, name)
			continue
		}
		thePkg = pkg
	}

	fset := iprog.Fset
	fname := func(p token.Pos) string {
		ret := filepath.Base(fset.Position(p).Filename)
		if strings.HasSuffix(ret, "_test.go") {
			return strings.TrimSuffix(ret, "_test.go")
		} else if strings.HasSuffix(ret, ".go") {
			return strings.TrimSuffix(ret, ".go")
		}
		return ret
	}

	info := iprog.AllPackages[thePkg]
	for _, f := range info.Files {
		ret.addFile(fname(f.Package))
	}

	res := lexical.Structure(iprog.Fset, thePkg, &info.Info, info.Files)
	for obj, refs := range res.Refs {
		pack := obj.Pkg()
		if pack != thePkg {
			// we only care about internal references
			continue
		}

		fdef := fname(obj.Pos())
		for _, ref := range refs {
			fused := fname(ref.Id.NamePos)
			if fused != fdef {
				ret.depends(fdef, fused)
			}
		}
	}

	return ret, nil
}

func ne(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "error:", e)
		os.Exit(-1)
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "requires exacly one package name")
		os.Exit(-1)
	}

	p, e := NewPackage(args[0])
	ne(e)

	e = p.MakeLevelsBot()
	if e != nil {
		p.PrintFileDeps()
	}
	ne(e)

	p.PrintLevelsBot()
}
