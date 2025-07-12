package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	_ "embed"
)

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
