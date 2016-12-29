# Contributing Guidelines

## Golden rule

Speak up _before_ writing code. Comment on existing issue or create a new one. Join
[Gitter chat](https://gitter.im/go-reform/reform?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge). Discuss what
you want to implement _before_ implementing it.


## Running tests

To run tests locally you have two options: installing database systems locally, or use Docker and Drone CLI.


### Local install

You should install _some_ versions of PostgreSQL and MySQL.
For Mac with Homebrew this should work:
```
brew update
brew install postgresql
brew services start postgresql
brew install mysql
brew services start mysql
```

Download dependencies with `make download_deps` and install them with `make install_deps`.
Then run tests with `make`. See [`Makefile`](../Makefile) for connection parameters.


### Drone

If you have Docker, you can use [Drone CI](http://readme.drone.io/0.5/) CLI to run tests without installing database
systems. First, you should install [Drone 0.5 CLI](http://readme.drone.io/0.5/install/cli/). Then, run tests with
`make drone`.
