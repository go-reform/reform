# Contributing Guidelines

## Golden rule

Speak up _before_ writing code. Comment on existing issue or create a new one. Discuss what
you want to implement _before_ implementing it.


## Running tests

First of all, run `env GO111MODULE=on go get -v ./...` to install versioned dependencies.

Then run all tests with Docker Compose: `make test`.
