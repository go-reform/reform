# Contributing Guidelines

## Golden rule

Speak up _before_ writing code. Comment on existing issue or create a new one. Join
[Gitter chat](https://gitter.im/go-reform/reform?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge). Discuss what
you want to implement _before_ implementing it.


## Running tests

First of all, run `make deps deps-check` to install all dependencies. After that, you have two options: use Docker Compose (recommended), or installing database systems directly.


### Docker Compose

If you have Go, Docker and Docker Compose installed, you can run all tests and linters simply by running `make test-dc check`.

You can also set `REFORM_IMAGE_VERSION` and `REFORM_TARGETS` environment variables to test specific combinations.
See [`.travis.yml`](../.travis.yml) for possible values.
You can also set `REFORM_SKIP_PULL` to `1` to avoid running `docker-compose pull` if image is already present.

### Direct

Run `make test` to run basic unit tests. Run `make check` to run linters.
See [`Makefile`](../Makefile) for Make targets for running integration tests and connection parameters.


### Background information

See [#5](https://github.com/go-reform/reform/issues/5), [#63](https://github.com/go-reform/reform/issues/63), and [#135](https://github.com/go-reform/reform/issues/135) for reasons for that design.
