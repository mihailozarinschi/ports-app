//go:build tools
// +build tools

package ports_app

// These imports are needed to have the tools versioned within our vendor
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
