package lib

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type State struct {
	DryRun bool
	Base   string
	Root   string
	Files  []string
}

func IsTruthy(s string) bool {
	return s == "1" || s == "true"
}

func FindRoot(file string) string {
	base := path.Dir(file)
	if _, err := os.Stat(path.Join(base, "go.mod")); errors.Is(err, os.ErrNotExist) {
		base = FindRoot(base)
	}
	return base
}

func Process(file string, state *State) {
	src, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.ImportsOnly)
	if err != nil {
		panic(err)
	}

	state.Files = append(state.Files, file)

	// inspect the AST and process all imports.
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.ImportSpec:
			s = x.Path.Value
			s = s[1 : len(s)-1]
		}

		if !strings.HasPrefix(s, state.Base) {
			return true
		}

		// if s != "" {
		// 	fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		// }
		nextFolder := strings.ReplaceAll(s, state.Base, state.Root)

		arr, err := os.ReadDir(nextFolder)
		if err != nil {
			panic(err)
		}
		for _, de := range arr {
			if de.IsDir() || !strings.HasSuffix(de.Name(), ".go") {
				continue
			}
			nextFile := path.Join(nextFolder, de.Name())
			if slices.Contains(state.Files, nextFile) {
				continue
			}
			Process(nextFile, state)
		}
		return true
	})
}

func Prune(state *State) {
	err := filepath.Walk(state.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") || info.IsDir() {
			return nil
		}
		if !slices.Contains(state.Files, path) {
			if !state.DryRun {
				os.Remove(path)
			}
			fmt.Println("removed:", path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
