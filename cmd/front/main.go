package main

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	flag "github.com/spf13/pflag"
)

var version string

//go:embed usage.txt
var usageFmt string

var versionFmt = "%[1]s %[2]s\n"

func main() {
	var (
		delim  = "---"
		input  = "-"
		output = "-"
		run    = func() { runProgram(delim, input, output) }
	)

	flag.Usage = printUsage

	flag.BoolFuncP("help", "h", "Show help", func(string) error {
		run = printUsage
		return nil
	})

	flag.BoolFuncP("version", "v", "Show version", func(string) error {
		run = printVersion
		return nil
	})

	flag.StringVarP(&delim, "delim", "d", delim,
		"Set frontmatter delimiter to `DELIM`")

	flag.StringVarP(&output, "out", "o", output,
		"Output to `FILE`")

	flag.Parse()

	if flag.NArg() > 0 {
		input = flag.Arg(0)
	}

	if input == "" {
		input = "-"
	}

	run()
}

func die(err error, op ...string) {
	fmt.Fprintf(os.Stderr,
		"failed to %s: %s", strings.Join(op, " "), err)
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), usageFmt,
		flag.CommandLine.Name(), flag.CommandLine.FlagUsagesWrapped(55))
}

func printVersion() {
	fmt.Fprintf(flag.CommandLine.Output(), versionFmt,
		path.Base(flag.CommandLine.Name()), version)
}
