package front

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type Splitter struct {
	Delim string
}

func (s Splitter) SplitFile(file string) ([]byte, []byte, error) {
	if r, err := os.Open(file); err != nil {
		return nil, nil, err
	} else {
		return splitFront(r, s.Delim)
	}
}

func (s Splitter) SplitBytes(b []byte) ([]byte, []byte, error) {
	return splitFront(bytes.NewReader(b), s.Delim)
}

func (s Splitter) SplitReader(r io.Reader) ([]byte, []byte, error) {
	return splitFront(r, s.Delim)
}

//

var Default = Splitter{"---"}

func SplitFile(file string) ([]byte, []byte, error) {
	return Default.SplitFile(file)
}

func SplitBytes(b []byte) ([]byte, []byte, error) {
	return Default.SplitBytes(b)
}

func SplitReader(r io.Reader) ([]byte, []byte, error) {
	return Default.SplitReader(r)
}

//

type statefn func(state) (state, statefn, error)

type state struct {
	rd      *bufio.Reader
	delim   []byte
	delimln []byte
	meta    []byte
	content []byte
}

func splitFront(r io.Reader, delim string) ([]byte, []byte, error) {
	var err error

	s := state{
		rd:      bufio.NewReader(r),
		delim:   []byte(removeNewLine(delim)),
		delimln: []byte(addNewLine(delim)),
	}

	for fn := startState; ; {
		if s, fn, err = fn(s); err != nil || fn == nil {
			break
		}
	}

	return s.meta, s.content, ignoreEOF(err)
}

func startState(s state) (state, statefn, error) {
	line, err := s.rd.ReadBytes('\n')

	if equalAny(line, s.delim, s.delimln) {
		return s, metaState, err
	} else {
		s.content = line
		return s, contentState, err
	}
}

func metaState(s state) (state, statefn, error) {
	line, err := s.rd.ReadBytes('\n')

	if equalAny(line, s.delim, s.delimln) {
		return s, contentState, err
	} else {
		s.meta = append(s.meta, line...)
		return s, metaState, err
	}
}

func contentState(s state) (state, statefn, error) {
	all, err := io.ReadAll(s.rd)
	s.content = append(s.content, all...)
	return s, nil, err
}

//

func addNewLine(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	} else {
		return s
	}
}

func removeNewLine(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s[1 : len(s)-1]
	} else {
		return s
	}
}

func ignoreEOF(err error) error {
	if err == io.EOF {
		return nil
	} else {
		return err
	}
}

func equalAny(b []byte, bs ...[]byte) bool {
	for _, s := range bs {
		if bytes.Equal(b, s) {
			return true
		}
	}
	return false
}
