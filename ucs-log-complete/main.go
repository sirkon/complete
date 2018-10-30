package main

import (
	"github.com/sirkon/complete"
	"github.com/sirkon/complete/internal/logroot"
)

func main() {
	home := complete.Command{
		TrustedArgs: logroot.LogRoot("/var/log/ucs/dev"),
	}

	// run the command completion, as part of the main() function.
	// this triggers the autocompletion when needed.
	// name must be exactly as the binary that we want to complete.
	complete.New("ucs-log", home).Run()
}
