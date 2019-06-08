package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirkon/complete"
)

var workDir = "gitlab.stageoffice.ru"
var gopathSrc string
var Log = complete.Log

func cutParent(path, parent string) string {
	if !strings.HasPrefix(path, parent) {
		return path
	}
	return path[len(parent)+1:]
}

type workProjects string

func (s workProjects) Predict(args complete.Args) (res []string) {
	defer func() {
		for i, path := range res {
			res[i] = cutParent(path, gopathSrc)
		}
	}()
	var root string
	var dirs []string
	var err error

	switch {
	case args.Last == workDir:
		return complete.PredictDirs(filepath.Join(gopathSrc, args.Last)).Predict(args)
	case strings.HasPrefix(args.Last, workDir):
		var dirs []string
		var err error
		if stat, err := os.Stat(filepath.Join(gopathSrc, args.Last)); err == nil && stat.IsDir() {
			dirs, err = filepath.Glob(filepath.Join(gopathSrc, args.Last, "*"))
		} else {
			dirs, err = filepath.Glob(filepath.Join(gopathSrc, args.Last+"*"))
		}
		if err != nil {
			Log("cant get list of directories in %workDir: %workDir", root, err)
			os.Exit(1)
		}
		if len(dirs) == 1 && strings.Count(cutParent(dirs[0], gopathSrc), "/") == 1 {
			Log("matched %s, discovering deeper into it", dirs[0])
			args.LastCompleted = args.Last
			args.Last = cutParent(dirs[0], gopathSrc)
			return s.Predict(args)
		}
		return dirs
	}

	root = filepath.Join(gopathSrc, string(s))
	dirs, err = filepath.Glob(filepath.Join(root, "*"))
	if err != nil {
		Log("cant get list of directories in %workDir: %workDir", root, err)
		os.Exit(1)
	}
	Log("list of dirs: %workDir", dirs)
	for _, dir := range dirs {
		mask := filepath.Join(dir, args.Last+"*")
		files, err := filepath.Glob(mask)
		if err != nil {
			Log("cannot get list of files with mask %workDir: %workDir", mask, err)
			os.Exit(1)
		}
		res = append(res, files...)
	}
	return res
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

	work := complete.Command{
		TrustedArgs: workProjects(workDir),
	}
	complete.New("work", work).Run()
}
