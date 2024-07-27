package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bgpat/flatjson"
)

var diffMode bool

var stdin io.ReadCloser = os.Stdin

func init() {
	flag.BoolVar(&diffMode, "diff", false, "enable diff mode: -diff FILE1 FILE2")
}

func main() {
	flag.Parse()

	if diffMode {
		code, err := diff(flag.Args(), os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		os.Exit(code)
		return
	}

	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(r io.ReadCloser, w io.WriteCloser) error {
	f, err := flatten(r)
	if err != nil {
		return err
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(f)
}

func diff(args []string, w io.WriteCloser) (int, error) {
	if len(args) != 2 {
		return 2, errors.New("required 2 args in diff mode")
	}

	x, xc, err := load(args[0])
	if err != nil {
		return 2, err
	}
	defer xc.Close()

	y, yc, err := load(args[1])
	if err != nil {
		return 2, err
	}
	defer yc.Close()

	diff := x.Diff(y)

	if diff == nil {
		// no changes
		return 0, nil
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return 1, e.Encode(diff)
}

func flatten(r io.ReadCloser) (flatjson.FlatJSON, error) {
	d := json.NewDecoder(r)
	var v any
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	return flatjson.Flatten(v)
}

func load(name string) (flatjson.FlatJSON, io.Closer, error) {
	var input io.ReadCloser
	if name == "-" {
		var b bytes.Buffer
		input = io.NopCloser(io.TeeReader(stdin, &b))
		stdin = io.NopCloser(&b)
	} else {
		f, err := os.Open(name)
		if err != nil {
			return nil, nil, err
		}
		input = f
	}
	f, err := flatten(input)
	if err != nil {
		input.Close()
		return nil, nil, err
	}
	return f, input, nil
}
