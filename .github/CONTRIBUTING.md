# Contributing Guidelines

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

Then run tests with `make`. See [`Makefile`](../Makefile) for connection parameters.


### Drone

If you have Docker, you can use [Drone CI](http://readme.drone.io/0.5/) CLI to run tests without installing database
systems. First, you should install [Drone 0.5 CLI](http://readme.drone.io/0.5/install/cli/). Then, run tests with
`make drone`.
