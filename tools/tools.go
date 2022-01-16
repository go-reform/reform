//go:build tools
// +build tools

package tools // import "gopkg.in/reform.v1/tools"

import (
	_ "github.com/AlekSi/gocoverutil"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/quasilyte/go-consistent"
	_ "github.com/reviewdog/reviewdog/cmd/reviewdog"
	_ "mvdan.cc/gofumpt"
)

//go:generate go build -v -o ../bin/gocoverutil github.com/AlekSi/gocoverutil
//go:generate go build -v -o ../bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go build -v -o ../bin/go-consistent github.com/quasilyte/go-consistent
//go:generate go build -v -o ../bin/reviewdog github.com/reviewdog/reviewdog/cmd/reviewdog
//go:generate go build -v -o ../bin/gofumpt mvdan.cc/gofumpt
