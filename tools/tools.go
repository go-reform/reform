//go:build tools
// +build tools

package tools // import "gopkg.in/reform.v1/tools"

import (
	_ "github.com/AlekSi/gocoverutil"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/quasilyte/go-consistent"
	_ "github.com/reviewdog/reviewdog/cmd/reviewdog"
	_ "golang.org/x/tools/cmd/goimports"
)
