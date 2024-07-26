package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bgpat/flatjson"
)

func main() {
	d := json.NewDecoder(os.Stdin)
	var v any
	if err := d.Decode(&v); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f, err := flatjson.Flatten(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
