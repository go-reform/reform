package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
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
		[]string{"postgres", "pgx"},
		[]string{
			"9.4",
			"9.5",
			"9.6",
			"10",
			"11",
			"12",
		},
		"11",
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
		for _, v := range c.ImageVersions {
			fmt.Fprint(&buf, "    - ")
			fmt.Fprintf(&buf, "REFORM_IMAGE_VERSION=%s ", v)
			fmt.Fprintf(&buf, "REFORM_TARGETS=%s\n", strings.Join(c.Targets, ","))
			count++
		}
	}

	const filename = ".travis.yml"
	const start = "# Generated with 'go run .github/test-dc.go gen'."
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

// testTargets tests specific combinations.
func testTargets() {
	v := os.Getenv("REFORM_IMAGE_VERSION")
	for _, t := range strings.Split(os.Getenv("REFORM_TARGETS"), ",") {
		log.Printf("REFORM_IMAGE_VERSION=%s target=%s ", v, t)

		var commands []string
		if skipPull, _ := strconv.ParseBool(os.Getenv("REFORM_SKIP_PULL")); !skipPull {
			commands = append(commands, fmt.Sprintf("docker-compose --file=.github/docker-compose-%s.yml --project-name=reform pull", t))
		}
		commands = append(commands, fmt.Sprintf("docker-compose --file=.github/docker-compose-%s.yml --project-name=reform up -d --remove-orphans --force-recreate", t))
		commands = append(commands, fmt.Sprintf("make %s", t))
		commands = append(commands, fmt.Sprintf("docker-compose --file=.github/docker-compose-%s.yml --project-name=reform down --remove-orphans --volumes", t))
		for _, c := range commands {
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
}

func test() {
	// Test only specific combinations if REFORM_TARGETS environment variable is defined.
	// Used by Travis CI for matrix builds.
	if os.Getenv("REFORM_TARGETS") != "" {
		testTargets()
		return
	}

	// test default combinations one by one
	for _, c := range configs {
		for _, t := range c.Targets {
			os.Setenv("REFORM_IMAGE_VERSION", c.DefaultImageVersion)
			os.Setenv("REFORM_TARGETS", t)
			testTargets()
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
