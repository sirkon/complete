package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sirkon/complete"
)

/*
здесь у нас располагется поиск проектов относительно $GOPATH/src
*/

type gopathprojects struct{}

type projectRoot struct {
	root string
}

// Root возращает $GOPATH/src
func (pr *projectRoot) Root() string {
	if len(pr.root) == 0 {
		gopath := os.Getenv("GOPATH")
		if len(gopath) == 0 {
			home, _ := os.UserHomeDir()
			gopath = filepath.Join(home, "go")
		}
		pr.root = filepath.Join(gopath, "src")
	}
	return pr.root
}

// Rel для данного пути path, который должен быть внутри путя $GOPATH/src, возвращает его путь относительно $GOPATH/src
func (pr *projectRoot) Rel(path string) string {
	base := pr.Root()
	if !strings.HasPrefix(path, base) {
		return ""
	}
	if path == base {
		return ""
	}
	return path[len(base)+1:]
}

// Full расширяет данный путь до подпутя в $GOPATH/src
func (pr *projectRoot) Full(path string) string {
	return filepath.Join(pr.Root(), path)
}

// Predict для реализации интерфейса Predictor
func (gopathprojects) Predict(args complete.Args) (res []string) {
	pr := projectRoot{}
	path := pr.Full(args.Last)

	if !isValidDir(path) {
		return
	}

	if stat, err := os.Stat(path); err == nil {
		if stat.IsDir() {
			if isGoDir, _ := checkProjectDir(path); isGoDir {
				res = append(res, pr.Rel(args.Last))
			}
		}
	}

	items, err := filepath.Glob(path + "*")
	if err != nil {
		complete.Log(err.Error())
		return
	}

	var dirs sort.StringSlice
	for _, item := range items {
		if !isValidDir(item) {
			continue
		}
		stat, err := os.Stat(item)
		if err != nil {
			complete.Log("cannot check file: %s", err)
		}
		if stat.IsDir() {
			isGoDir, newDirs := checkProjectDir(item)
			if isGoDir {
				res = append(res, pr.Rel(item))
			}
			for _, newDir := range newDirs {
				newDir = pr.Rel(newDir)
				if !isValidDir(newDir) {
					continue
				}
				dirs = append(dirs, newDir)
			}
		}
	}
	dirs.Sort()
	res = append(res, dirs...)
	return
}

// checkProjectDir выясняет, являетя ли данная директория Go-шным пакетом, а так же возвращает список каталогов внутри
func checkProjectDir(dir string) (yes bool, dirs []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, nil
	}
	for _, item := range files {
		if item.IsDir() {
			dirs = append(dirs, filepath.Join(dir, item.Name()))
		} else if !yes && strings.HasSuffix(item.Name(), ".go") {
			yes = true
		}
	}
	return
}

func isValidDir(name string) bool {
	parts := strings.Split(name, string(os.PathSeparator))
	for _, part := range parts {
		if strings.HasPrefix(part, ".") || part == "vendor" {
			return false
		}
	}
	return true
}
