# Contributing Guidelines

## Golden rule

Speak up _before_ writing code. Comment on existing issue or create a new one. Discuss what
you want to implement _before_ implementing it.


## Running tests

1. Reform uses Go modules to version dependencies. Make sure they are not disabled in your environment.
2. Run `make` without arguments to see all Makefile targets.
3. Run `make env-up` or `make env-up-detach` to start databases with Docker Compose.
4. Run `make init` to install development tools.
5. Run `make test` to run all tests.
