package front

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type result struct {
	meta, content string
}

func TestSplitReader(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect result
	}{
		{
			name:   "empty",
			input:  "",
			expect: result{},
		},
		{
			name:   "delim_only",
			input:  "---",
			expect: result{},
		},
		{
			name:   "delim_only_eol",
			input:  "---\n",
			expect: result{},
		},
		{
			name:  "content_only",
			input: "content",
			expect: result{
				content: "content",
			},
		},
		{
			name:  "meta_only",
			input: "---\nmeta\n---",
			expect: result{
				meta: "meta\n",
			},
		},
		{
			name:  "meta_only_eol",
			input: "---\nmeta\n---\n",
			expect: result{
				meta: "meta\n",
			},
		},
		{
			name:  "meta_content",
			input: "---\nmeta\n---\ncontent",
			expect: result{
				meta:    "meta\n",
				content: "content",
			},
		},
		{
			name:  "meta_content_eol",
			input: "---\nmeta\n---\ncontent\n",
			expect: result{
				meta:    "meta\n",
				content: "content\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader([]byte(tt.input))
			m, c, err := Default.SplitReader(b)

			assert.Nil(t, err)
			assert.Equal(t, tt.expect.meta, string(m))
			assert.Equal(t, tt.expect.content, string(c))
		})
	}
}
