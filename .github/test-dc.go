package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type config struct {
	ComposeName         string
	ImageVersions       []string
	DefaultImageVersion string
	MakeTargets         []string
}

var configs = []config{
	{
		"cockroach",
		[]string{
			"latest",
		},
		"latest",
		[]string{"cockroach-postgres", "cockroach-pgx"},
	},

	// https://www.postgresql.org/support/versioning/
	{
		"postgres",
		[]string{
			"9.3",
			"9.4",
			"9.5",
			"9.6",
			"10",
		},
		"10",
		[]string{"postgres", "pgx"},
	},

	// https://www.mysql.com/support/supportedplatforms/database.html
	{
		"mysql",
		[]string{
			"5.5",
			"5.6",
			"5.7",
			"8.0",
		},
		"5.7",
		[]string{"mysql", "mysql-traditional"},
	},

	{
		"sqlite3",
		[]string{
			"dummy",
		},
		"dummy",
		[]string{"sqlite3"},
	},

	{
		"mssql",
		[]string{
			"latest",
		},
		"latest",
		[]string{"mssql", "sqlserver"},
	},
}

// gen updates matrix configuration in .travis.yml file
func gen() {
	var buf bytes.Buffer
	for _, c := range configs {
		s := fmt.Sprintf("  #   %s: %s\n", strings.Join(c.MakeTargets, ", "), strings.Join(c.ImageVersions, ", "))
		buf.WriteString(s)
	}

	buf.WriteString("  matrix:")

	var count int
	for _, c := range configs {
		fmt.Fprint(&buf, "\n")
		for _, v := range c.ImageVersions {
			for _, t := range c.MakeTargets {
				fmt.Fprint(&buf, "    - ")
				fmt.Fprintf(&buf, "REFORM_COMPOSE_NAME=%s ", c.ComposeName)
				fmt.Fprintf(&buf, "REFORM_IMAGE_VERSION=%s ", v)
				fmt.Fprintf(&buf, "REFORM_MAKE_TARGET=%s\n", t)
				count++
			}
		}
	}

	const filename = ".travis.yml"
	const start = "# Generated with 'go run .github/test-dc.go'."
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	i := bytes.Index(b, []byte(start))
	if i < 0 {
		log.Fatalf("failed to find %q in %s", start, filename)
	}
	b = append(b[:i], start...)
	b = append(b, []byte(fmt.Sprintf("\n  # %d combinations:\n", count))...)
	b = append(b, buf.Bytes()...)
	if err = ioutil.WriteFile(filename, b, 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s matrix updated with %d combinations.", filename, count)
}

// testOne tests a single configuration
func testOne() {
	n := os.Getenv("REFORM_COMPOSE_NAME")
	v := os.Getenv("REFORM_IMAGE_VERSION")
	t := os.Getenv("REFORM_MAKE_TARGET")
	log.Printf("REFORM_COMPOSE_NAME=%s REFORM_IMAGE_VERSION=%s REFORM_MAKE_TARGET=%s", n, v, t)
	f := fmt.Sprintf(".github/docker-compose-%s.yml", n)

	for _, c := range []string{
		fmt.Sprintf("docker-compose --file=%s --project-name=reform pull", f),
		fmt.Sprintf("docker-compose --file=%s --project-name=reform up -d --remove-orphans --force-recreate", f),
		fmt.Sprintf("make %s", t),
		fmt.Sprintf("docker-compose --file=%s --project-name=reform down --remove-orphans --volumes", f),
	} {
		log.Print(c)
		args := strings.Split(c, " ")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func test() {
	// Test only a single configuration if REFORM_MAKE_TARGET environment variable is defined.
	// Used by Travis CI for matrix builds.
	if os.Getenv("REFORM_MAKE_TARGET") != "" {
		testOne()
		return
	}

	// test all configurations one by one
	for _, c := range configs {
		for _, t := range c.MakeTargets {
			os.Setenv("REFORM_COMPOSE_NAME", c.ComposeName)
			os.Setenv("REFORM_IMAGE_VERSION", c.DefaultImageVersion)
			os.Setenv("REFORM_MAKE_TARGET", t)
			testOne()
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("test-dc: ")

	commands := []string{"gen", "test"}
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Expected exactly one argument, one of: %s.", strings.Join(commands, ", "))
	}

	switch flag.Arg(0) {
	case "gen":
		gen()
	case "test":
		test()
	default:
		log.Fatal("unexpected command")
	}
}
