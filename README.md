# goprune

Prune go project with source files not been used.

## Install

```sh
go install github.com/futurist/goprune@latest
```

## Usage

```sh
goprune abs_path_of_main_or_lib.go ...
```

Below steps will be applied:

1. Find closest `root` dir of the go project with a `go.mod` file.
2. Set `base` to go module path defined in `module` field of `go.mod`.
3. Start from `abs_path_of_main_or_lib.go`, parse all imported packages that start with `base`.
4. Loop step 3 with imported packages until no files left.
5. Remove all go source files that have not been walked (not used).

With `DRY_RUN=1` environment variable been set, will print removed files but not do the removal in step 5.
With `NO_TEST=1` environment variable been set, will also remove go test files in step 5.
