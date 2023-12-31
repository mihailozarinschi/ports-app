package goimports

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/tools/imports"
)

// Run runs goimports.
// The local prefixes (comma separated) must be defined through the global variable imports.LocalPrefix.
func Run(filename string) ([]byte, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	res, err := imports.Process(filename, src, nil)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(src, res) {
		return nil, nil
	}

	// formatting has changed
	data, err := diff(src, res, filename)
	if err != nil {
		return nil, fmt.Errorf("error computing diff: %s", err)
	}

	return data, nil
}
