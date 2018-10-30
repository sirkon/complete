package logroot

import (
	"github.com/sirkon/complete"
	"os"
	"path/filepath"
	"strings"
)

type LogRoot string

var Log = complete.Log

func (lr LogRoot) Predict(args complete.Args) (result []string) {
	defer func() {
		Log("################### %v", result)
	}()
	if len(args.All) == 1 {
		// first parameter must be some trace id which cannot be completed for obvious reasons
		return nil
	}
	if len(args.All) > 2 {
		// no parameters after that
		Log("Warning many params: %s", args.All)
	}
	if args.Last == string(lr) {
		return expandRoot(string(lr))
	}
	if strings.HasPrefix(args.Last, string(lr)) {
		return expandPath(args.Last)
	} else {
		return expandPath(filepath.Join(string(lr), args.Last))
	}
}

// imply "split logic", i.e. atl/atl will be expanded into atlas/atlas.log if there's 'atlas/atlas.log'
// in a root directory
func (lr LogRoot) splitLogic(args complete.Args) []string {
	return nil
}

func expandRoot(root string) (result []string) {
	defer func() {
		Log("damn %s %v", root, result)
	}()
	dirs, err := filepath.Glob(filepath.Join(root, "*"))
	if err != nil {
		os.Exit(1)
	}
	if len(dirs) == 1 {
		res := expandRoot(dirs[0])
		if len(res) == 0 {
			return dirs
		} else if len(res) == 1 && res[0] == dirs[0] {
			return res
		}
		res = expandPath(dirs[0])
		if len(res) == 0 {
			return dirs
		}
		return res
	}
	return dirs
}

func expandPath(path string) (result []string) {
	defer func() {
		Log("shit %s %v", path, result)
	}()
	items, err := filepath.Glob(path + "*")
	Log("I am here: %v", items)
	if err != nil {
		Log("Cannot expand path %s: %s", path, items)
		os.Exit(1)
	}
	if len(items) == 1 {
		res := expandRoot(items[0])
		if len(res) == 0 {
			return items
		} else {
			return res
		}
		res = expandPath(items[0])
		if len(res) == 0 {
			return items
		}
		return res
	}
	return items
}
