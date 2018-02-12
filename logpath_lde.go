package complete

import (
	"strings"
)

var fileColon = "file:"

// LogFilePath ...
type LogFilePath struct {
	rest string
	File string
}

// Extract ...
func (p *LogFilePath) Extract(line string) (bool, error) {
	p.rest = line

	// Checks if the rest starts with `"file:"` and pass it
	if strings.HasPrefix(p.rest, fileColon) {
		p.rest = p.rest[len(fileColon):]
	} else {
		return false, nil
	}

	// Take the rest as File(string)
	p.File = p.rest
	p.rest = p.rest[len(p.rest):]
	return true, nil
}
