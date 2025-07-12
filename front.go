package front

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

const DefaultMarker = "---"

var DefaultSplitter = NewSplitter()

type Splitter struct {
	marker, delim []byte
}

func NewSplitter() Splitter {
	return NewSplitterMark(DefaultMarker)
}

func NewSplitterMark(m string) Splitter {
	if n := len(m) - 1; m[n] == '\n' {
		m = m[0:n]
	}
	return Splitter{
		[]byte(m),
		[]byte(m + "\n"),
	}
}

func (s Splitter) Split(r io.Reader) ([]byte, []byte, error) {
	var (
		matter, content []byte
		b               = bufio.NewReader(r)
		magic, front    = s.readMagic(b)
	)

	if front {
		// file contains valid frontmatter
		for line := range s.readLines(b) {
			matter = append(matter, line...)
		}
	} else {
		// file has no valid frontmatter,
		// re-consume assumed 'magic' bytes into split.content
		content = append(content, magic...)
	}

	// read the rest of the file
	all, err := io.ReadAll(b)
	content = append(content, all...)
	return matter, content, err
}

func Split(r io.Reader) ([]byte, []byte, error) {
	return DefaultSplitter.Split(r)
}

func (s Splitter) readMagic(r io.Reader) ([]byte, bool) {
	var (
		n   int
		err error
	)

	magic := make([]byte, len(s.delim))

	if n, err = io.ReadFull(r, magic); err != nil {
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
		return magic[0:n], bytes.Equal(fmm, s.marker)
	}
}

func (s Splitter) readLines(r io.Reader) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		b := bufio.NewReader(r)

		for {
			line, err := b.ReadBytes('\n')

			if bytes.Equal(line, s.marker) ||
				bytes.Equal(line, s.delim) {
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
