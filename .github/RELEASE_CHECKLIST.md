# Release Checklist

1. Check tests, linters.
2. Check issues and pull requests, update milestones.
3. Bump version in `doc.go`.
4. Update CHANGELOG.md.
5. Make tag.
6. Push it!
7. Make [release](https://github.com/go-reform/reform/releases).
8. Refresh
   * `env GO111MODULE=on GOPROXY=https://proxy.golang.org go get -v gopkg.in/reform.v1@v1.M.P`
   * `env GO111MODULE=on GOPROXY=https://gocenter.io      go get -v gopkg.in/reform.v1@v1.M.P`
   * https://pkg.go.dev/gopkg.in/reform.v1@v1.M.P
   * https://godoc.org/gopkg.in/reform.v1
   * https://goreportcard.com/report/gopkg.in/reform.v1
