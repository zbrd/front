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
	prog    usage.Program
	version string

	opts = &Options{
		Input:  "-",
		Output: "-",
	}

	info = usage.Data{
		"Author":  "Zaim B.",
		"Contact": "zbrd@duck.com",
		"Version": version,
	}

	matterDelim = []byte("---\n")
	matterMark  = matterDelim[0:3]
)

//go:embed usage.txt
var usageTpl string

//go:embed version.txt
var versionTpl string

// Types
// -----

type Options struct {
	Input       string
	Output      string
	ShowUsage   bool
	ShowVersion bool
}

func (o Options) Reader() (r io.Reader, err error) {
	if opts.Input == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(o.Input)
	}
	return
}

func (o Options) Writer() (w io.WriteCloser, err error) {
	if opts.Output == "-" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(o.Output, os.O_WRONLY|os.O_CREATE, 0644)
	}
	return
}

type Output struct {
	Meta    map[string]any `json:"meta"`
	Path    string         `json:"path"`
	Content string         `json:"content"`
}

func (o Output) SetContent(c []byte) Output {
	o.Content = string(c)
	return o
}

// Main
// ----

func init() {
	prog = usage.Prog(flag.CommandLine)
	flag.Usage = showUsage
	flag.CommandLine.SetOutput(os.Stdout)

	flag.StringVarP(&opts.Output, "out", "o", opts.Output, "")
	flag.BoolVarP(&opts.ShowUsage, "help", "h", false, "")
	flag.BoolVarP(&opts.ShowVersion, "version", "v", false, "")

	flag.Lookup("out").Usage = "Output file `PATH`"
	flag.Lookup("help").Usage = "Show help information"
	flag.Lookup("version").Usage = "Show version information"
}

func main() {
	flag.Parse()

	switch {
	case opts.ShowUsage:
		showUsage()
		return
	case opts.ShowVersion:
		showVersion()
		return
	}

	if flag.NArg() > 0 {
		opts.Input = flag.Arg(0)
	}

	if r, err := opts.Reader(); err != nil {
		exit("open input", err)
	} else if w, err := opts.Writer(); err != nil {
		exit("open output", err)
	} else if err := parseFront(opts.Input, r, w); err != nil {
		exit("parse frontmatter", err)
	}
}

func showUsage() {
	if err := prog.PrintUsage(usageTpl, info); err != nil {
		exit("print usage", err)
	}
}

func showVersion() {
	if err := prog.PrintUsage(versionTpl, info); err != nil {
		exit("print version", err)
	}
}

// Funcs
// -----

func exit(op string, err error) {
	fmt.Fprintf(os.Stderr, "Failed to %s: %s", op, err)
	os.Exit(1)
}

func splitMatter(r io.Reader) ([]byte, []byte, error) {
	var (
		matter, content []byte
		buff            = bufio.NewReader(r)
		magic, front    = readMagic(buff)
	)

	if front {
		// file contains valid frontmatter,
		// read it and store in split.matter
		for line := range matterLines(buff) {
			matter = append(matter, line...)
		}
	} else {
		// file has no valid frontmatter,
		// re-consume assumed 'magic' bytes into split.content
		content = append(content, magic...)
	}

	// read the rest of the file into split.content
	all, err := io.ReadAll(buff)
	content = append(content, all...)
	return matter, content, err
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

func parseFront(path string, in io.Reader, out io.Writer) error {
	var (
		err  error
		op   Output
		m, c []byte
	)

	op.Path = path

	if m, c, err = splitMatter(in); err != nil {
		return fmt.Errorf("Error splitting frontmatter: %w", err)
	} else if op.Meta, err = parseYAML(m); err != nil {
		return fmt.Errorf("Error parsing YAML: %w", err)
	} else if err := writeOutput(out, op.SetContent(c)); err != nil {
		return err
	} else {
		return nil
	}
}

func parseYAML(b []byte) (map[string]any, error) {
	var meta map[string]any

	if j, err := yaml.YAMLToJSON(b); err != nil {
		return nil, fmt.Errorf("Error YAML to JSON: %w", err)
	} else if err := json.Unmarshal(j, &meta); err != nil {
		return nil, err
	} else {
		return meta, nil
	}
}

func writeOutput(out io.Writer, output Output) error {
	if data, err := json.Marshal(output); err != nil {
		return fmt.Errorf("Error marshal output: %w", err)
	} else if _, err := out.Write(data); err != nil {
		return err
	} else {
		return nil
	}
}
