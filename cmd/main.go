package main

import (
	"fmt"
	"os"
	"path"

	"github.com/futurist/goprune/lib"
	"golang.org/x/mod/modfile"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goprune abs_path_of_main_or_lib.go")
		return
	}
	file := os.Args[1]
	if stat, err := os.Stat(file); err != nil {
		panic(err)
	} else if stat.IsDir() {
		file = path.Join(file, "main.go")
	}
	root := lib.FindRoot(file)

	gomod := path.Join(root, "go.mod")
	b, err := os.ReadFile(gomod)
	if err != nil {
		panic(err)
	}
	f, err := modfile.Parse(gomod, b, nil)
	if err != nil {
		panic(err)
	}

	base := f.Module.Mod.Path
	dryRun := lib.IsTruthy(os.Getenv("DRY_RUN"))
	fmt.Println("root:", root, "module:", base, "dryRun:", dryRun)

	state := lib.State{
		DryRun: dryRun,
		Root:   root,
		Base:   base,
	}
	lib.Process(file, &state)

	lib.Prune(&state)
}
