package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirkon/complete"
)

type doesntsuck struct{}

// Predict для реализации интерфейса Predictor
func (doesntsuck) Predict(args complete.Args) (res []string) {
	var err error
	defer func() {
		var output []string
		if err != nil {
			output = append(output, err.Error())
		}
		for _, file := range res {
			output = append(output, file)
		}
		ioutil.WriteFile("/home/emacs/Desktop/output", []byte(strings.Join(output, "\n")), 0666)
	}()
	gopath := os.Getenv("GOPATH")
	base := filepath.Join(gopath, "src")
	path := filepath.Join(base, args.Last)
	items, err := filepath.Glob(path + "*")
	if err != nil {
		return nil
	}
	if len(items) == 0 {
		err = fmt.Errorf("no completion for %s*", path)
		return
	}

	for _, item := range items {
		var subdirs []string
		stat, err := os.Stat(item)
		if err != nil {
			continue
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(item)
			if err != nil {
				continue
			}
			var isProjectDir bool
			var prefix string
			if base == item {
				prefix = ""
			} else {
				prefix = item[len(base)+1:]
			}
			for _, file := range files {
				if file.Name() == "vendor" {
					continue
				}
				if file.IsDir() {
					if strings.HasPrefix(file.Name(), ".") {
						continue
					}
					subdirs = append(subdirs, filepath.Join(prefix, file.Name()))
				} else if strings.HasSuffix(file.Name(), ".go") {
					// проверяем, что есть файлы *.go внутри, если да, то добавляем в результат
					isProjectDir = true
				}
			}
			if isProjectDir {
				res = append(res, prefix)
			}
			res = append(res, subdirs...)
			// добавляем в результат все директории
		}
	}
	return res
}
