package main

import (
	"github.com/sirkon/complete"
)

func main() {

	home := complete.Command{
		TrustedArgs: gopathprojects{},
	}

	// run the command completion, as part of the main() function.
	// this triggers the autocompletion when needed.
	// name must be exactly as the binary that we want to complete.
	complete.New("goland", home).Run()
}
