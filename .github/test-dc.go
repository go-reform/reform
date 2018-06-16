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
	Targets             []string // make targets, also parts for Docker Compose file name
	ImageVersions       []string
	DefaultImageVersion string
}

var configs = []config{
	// https://www.postgresql.org/support/versioning/
	{
		[]string{"postgres"},
		[]string{
			"9.3",
			"9.4",
			"9.5",
			"9.6",
			"10",
		},
		"10",
	},

	// https://www.mysql.com/support/supportedplatforms/database.html
	{
		[]string{"mysql", "mysql-traditional"},
		[]string{
			"5.5",
			"5.6",
			"5.7",
			"8.0",
		},
		"5.7",
	},

	{
		[]string{"sqlite3"},
		[]string{
			"dummy",
		},
		"dummy",
	},

	{
		[]string{"mssql", "sqlserver"},
		[]string{
			"latest",
		},
		"latest",
	},
}

// gen updates matrix configuration in .travis.yml file
func gen() {
	var buf bytes.Buffer
	for _, c := range configs {
		s := fmt.Sprintf("  #   %s: %s\n", strings.Join(c.Targets, ", "), strings.Join(c.ImageVersions, ", "))
		buf.WriteString(s)
	}

	buf.WriteString("  matrix:")

	var count int
	for _, c := range configs {
		fmt.Fprint(&buf, "\n")
		for _, t := range c.Targets {
			for _, v := range c.ImageVersions {
				fmt.Fprint(&buf, "    - ")
				fmt.Fprintf(&buf, "REFORM_TARGET=%s ", t)
				fmt.Fprintf(&buf, "REFORM_IMAGE_VERSION=%s\n", v)
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
	t := os.Getenv("REFORM_TARGET")
	v := os.Getenv("REFORM_IMAGE_VERSION")
	log.Printf("REFORM_TARGET=%s REFORM_IMAGE_VERSION=%s", t, v)

	for _, c := range []string{
		fmt.Sprintf("docker-compose --file=.github/docker-compose-%s.yml --project-name=reform pull", t),
		fmt.Sprintf("docker-compose --file=.github/docker-compose-%s.yml --project-name=reform up -d --remove-orphans --force-recreate", t),
		fmt.Sprintf("make %s", t),
		fmt.Sprintf("docker-compose --file=.github/docker-compose-%s.yml --project-name=reform down --remove-orphans --volumes", t),
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
	// Test only a single configuration if REFORM_TARGET environment variable is defined.
	// Used by Travis CI for matrix builds.
	if os.Getenv("REFORM_TARGET") != "" {
		testOne()
		return
	}

	// test all configurations one by one
	for _, c := range configs {
		for _, t := range c.Targets {
			os.Setenv("REFORM_TARGET", t)
			os.Setenv("REFORM_IMAGE_VERSION", c.DefaultImageVersion)
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
