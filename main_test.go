package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	_ "embed"
)

func TestSplitMatter(t *testing.T) {
	tests := []struct {
		input   string
		matter  string
		content string
	}{
		// invalid magic
		{"", "", ""},
		{"-", "", "-"},
		{"--", "", "--"},
		{"---", "", "---"},

		// valid magic, empty
		{"---\n", "", ""},
		{"---\n---", "", ""},
		{"---\n---\n", "", ""},

		// no frontmatter end
		{"---\nL1\nL2", "L1\nL2", ""},
		{"---\nL1\nL2\n", "L1\nL2\n", ""},

		// frontmatter, no content
		{"---\nL1\nL2\nL3\n---", "L1\nL2\nL3\n", ""},
		{"---\nL1\nL2\nL3\n---\n", "L1\nL2\nL3\n", ""},

		// frontmatter, with content
		{"---\nL1\nL2\nL3\n---\nDATA", "L1\nL2\nL3\n", "DATA"},
		{"---\nL1\nL2\nL3\n---\nDATA\n", "L1\nL2\nL3\n", "DATA\n"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test_%d", i), func(t *testing.T) {
			r := strings.NewReader(tt.input)
			m, c, err := splitMatter(r)

			if err != nil {
				t.Errorf("err != nil: %#v", err)
			}

			if sm := string(m); sm != tt.matter {
				t.Errorf("matter != %#v: %#v", tt.matter, sm)
			}

			if sm := string(c); sm != tt.content {
				t.Errorf("content != %#v: %#v", tt.content, sm)
			}
		})
	}
}

//go:embed example/hello.md
var testInput []byte

//go:embed example/hello.json
var testOutput []byte

func TestParseFront(t *testing.T) {
	var (
		in  = bytes.NewReader(testInput)
		out bytes.Buffer
	)

	err := parseFront("example/hello.md", in, &out)

	if err != nil {
		t.Errorf("err != nil: %s", err)
	}

	var expect, got map[string]any

	if err := json.Unmarshal(testOutput, &expect); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(expect, got) {
		t.Errorf(
			"output not equal: %#v != %#v",
			expect,
			got,
		)
	}
}
