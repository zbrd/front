package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/zbrd/front"
)

type statefn func(state) (state, statefn, error)

type state struct {
	s       front.Splitter
	in, out string
	rd      io.ReadCloser
	wr      io.WriteCloser
	meta    []byte
	content []byte
	data    map[string]any
}

func runProgram(delim, in, out string) {
	var err error

	s := state{
		s:   front.Splitter{Delim: delim},
		in:  in,
		out: out,
	}

	for fn := openInput; ; {
		if s, fn, err = fn(s); err != nil || fn == nil {
			break
		}
	}

	if err != nil {
		die(err)
	}
}

func openInput(s state) (state, statefn, error) {
	var err error

	if s.in == "-" {
		s.rd = os.Stdin
	} else if s.rd, err = os.Open(s.in); err != nil {
		return s, nil, fmt.Errorf("failed to open input: %w", err)
	}

	return s, splitFrontmatter, nil
}

func splitFrontmatter(s state) (state, statefn, error) {
	if meta, content, err := s.s.SplitReader(s.rd); err != nil {
		return s, nil, fmt.Errorf("failed to split: %w", err)
	} else {
		s.meta = meta
		s.content = content
	}

	return s, parseMeta, nil
}

func parseMeta(s state) (state, statefn, error) {
	var data map[string]any

	if b, err := yaml.YAMLToJSON(s.meta); err != nil {
		return s, nil, fmt.Errorf("failed to parse YAML: %w", err)
	} else if err := json.Unmarshal(b, &data); err != nil {
		return s, nil, fmt.Errorf("failed to parse JSON: %w", err)
	} else {
		s.data = data
	}

	return s, openOutput, nil
}

func openOutput(s state) (state, statefn, error) {
	var err error

	f := os.O_WRONLY | os.O_CREATE
	m := os.FileMode(0644)

	if s.out == "-" {
		s.wr = os.Stdout
	} else if s.wr, err = os.OpenFile(s.out, f, m); err != nil {
		return s, nil, fmt.Errorf("failed to open output: %w", err)
	}

	return s, outputResult, nil
}

func outputResult(s state) (state, statefn, error) {
	out := map[string]any{
		"path":    s.in,
		"meta":    s.data,
		"content": string(s.content),
	}

	if b, err := json.Marshal(out); err != nil {
		return s, nil, fmt.Errorf("failed to marshal JSON: %w", err)
	} else if _, err := s.wr.Write(b); err != nil {
		return s, nil, fmt.Errorf("failed to write output: %w", err)
	}

	s.rd.Close()
	s.wr.Close()

	return s, nil, nil
}
