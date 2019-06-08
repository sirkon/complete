package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirkon/complete"
)

var gopathSrc string
var Log = complete.Log

func cutParent(path, parent string) string {
	if !strings.HasPrefix(path, parent) {
		return path
	}
	return path[len(parent)+1:]
}

type homeProjects string

func (h homeProjects) Predict(args complete.Args) []string {
	var mask string
	switch {
	case args.Last == string(h):
		mask = filepath.Join(gopathSrc, args.Last, "*")
	case strings.HasPrefix(args.Last, string(h)):
		mask = filepath.Join(gopathSrc, args.Last+"*")
	default:
		mask = filepath.Join(gopathSrc, string(h), args.Last+"*")
	}
	prjs, err := filepath.Glob(mask)
	if err != nil {
		Log("cannot get list of files with mask %s: %s", mask, err)
		os.Exit(1)
	}
	for i, path := range prjs {
		prjs[i] = cutParent(path, gopathSrc)
	}
	return prjs
}

func initGlobals() {
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		home, _ := os.UserHomeDir()
		gopath = filepath.Join(home, "go")
	}
	gopathSrc = filepath.Join(gopath, "src")
}

func main() {
	initGlobals()

	home := complete.Command{
		TrustedArgs: homeProjects("github.com/sirkon"),
	}

	// run the command completion, as part of the main() function.
	// this triggers the autocompletion when needed.
	// name must be exactly as the binary that we want to complete.
	complete.New("home", home).Run()
}
