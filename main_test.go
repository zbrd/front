package main

import (
	"fmt"
	"strings"
	"testing"
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
			s, err := splitMatter(r)

			if err != nil {
				t.Errorf("err != nil: %#v", err)
			}

			if sm := string(s.matter); sm != tt.matter {
				t.Errorf("matter != %#v: %#v", tt.matter, sm)
			}

			if sm := string(s.content); sm != tt.content {
				t.Errorf("content != %#v: %#v", tt.content, sm)
			}
		})
	}
}
