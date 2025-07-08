package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"os"

	_ "embed"

	"github.com/goccy/go-yaml"
	flag "github.com/spf13/pflag"
	"github.com/zbrd/usage"
)

// Globals
// -------

var (
	prog usage.Program

	opts = &Options{
		input:  "-",
		output: "-",
	}

	matterDelim = []byte("---\n")
	matterMark  = matterDelim[0:3]
)

//go:embed usage.txt
var usageTpl string

// Types
// -----

type Options struct {
	input, output string
}

type Split struct {
	matter, content []byte
}

// Main
// ----

func init() {
	prog = usage.Prog(flag.CommandLine)
	flag.Usage = func() { doUsage() }
	flag.StringVarP(&opts.output, "out", "o", opts.output, "")
	flag.Lookup("out").Usage = "Output file `PATH`"
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		opts.input = flag.Arg(0)
	}

	if in, err := openInput(); err != nil {
		exit("open input", err)
	} else if split, err := splitMatter(in); err != nil {
		exit("split matter", err)
	} else if meta, err := parseYAML(split.matter); err != nil {
		exit("parse YAML matter", err)
	} else if err := writeOutput(meta, split.content); err != nil {
		exit("output result", err)
	}
}

func doUsage() {
	data := usage.Data{
		"Contact": "zbrd@duck.com",
	}
	if err := prog.PrintUsage(usageTpl, data); err != nil {
		exit("print usage", err)
	}
}

// Funcs
// -----

func exit(op string, err error) {
	fmt.Fprintf(os.Stderr, "Failed to %s: %s", op, err)
	os.Exit(1)
}

func openInput() (io.Reader, error) {
	if opts.input == "-" {
		return os.Stdin, nil
	} else {
		return os.Open(opts.input)
	}
}

func splitMatter(r io.Reader) (Split, error) {
	var (
		split        Split
		buff         = bufio.NewReader(r)
		magic, front = readMagic(buff)
	)

	if front {
		// file contains valid frontmatter,
		// read it and store in split.matter
		for line := range matterLines(buff) {
			split.matter = append(split.matter, line...)
		}
	} else {
		// file has no valid frontmatter,
		// re-consume assumed 'magic' bytes into split.content
		split.content = append(split.content, magic...)
	}

	// read the rest of the file into split.content
	all, err := io.ReadAll(buff)
	split.content = append(split.content, all...)
	return split, err
}

func readMagic(b *bufio.Reader) ([]byte, bool) {
	var (
		n   int
		err error
	)

	magic := make([]byte, len(matterDelim))

	if n, err = io.ReadFull(b, magic); err != nil {
		switch err {
		case io.ErrUnexpectedEOF:
			// file ended before reading len(magic) bytes,
			// meaning file *definitely* has no frontmatter
			return magic[0:n], false
		case io.EOF:
			// file is empty
			return nil, false
		default:
			// unexpected error type
			// TODO: panic? return error?
			return nil, false
		}
	} else {
		// we read exactly len(magic) bytes,
		// file *might* start with magic string
		fmm := magic[0:max(0, n-1)] // n-1: remove '\n'
		return magic[0:n], bytes.Equal(fmm, matterMark)
	}
}

func matterLines(b *bufio.Reader) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for {
			line, err := b.ReadBytes('\n')

			if bytes.Equal(line, matterMark) ||
				bytes.Equal(line, matterDelim) {
				return
			}
			if !yield(line) {
				return
			}
			if err != nil {
				return
			}
		}
	}
}

func parseYAML(b []byte) (map[string]any, error) {
	var meta map[string]any

	if j, err := yaml.YAMLToJSON(b); err != nil {
		return nil, err
	} else if err := json.Unmarshal(j, &meta); err != nil {
		return nil, err
	} else {
		return meta, nil
	}
}

func openOutput() (io.WriteCloser, error) {
	if opts.output == "-" {
		return os.Stdout, nil
	} else {
		return os.OpenFile(opts.output, os.O_WRONLY|os.O_CREATE, 0644)
	}
}

func writeOutput(meta map[string]any, content []byte) error {
	outmap := map[string]any{
		"path":    opts.input,
		"meta":    meta,
		"content": string(content),
	}

	if out, err := openOutput(); err != nil {
		return err
	} else if data, err := json.Marshal(outmap); err != nil {
		return err
	} else {
		out.Write(data)
		return nil
	}
}
