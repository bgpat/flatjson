# flatjson

[![Go Reference](https://pkg.go.dev/badge/github.com/bgpat/flatjson.svg)](https://pkg.go.dev/github.com/bgpat/flatjson)

`flatjson` is a Go package that provides utilities for flattening nested JSON structures into a sequence of path-value pairs.
This makes it easier to compare differences between JSON objects, which is useful for various purposes such as data synchronization, version control, and data analysis.

## Installation

```sh
go install github.com/bgpat/flatjson/cmd/flatjson
```

## Example

input:

```json
{
  "a": true,
  "b": {
    "c": 2
  },
  "d": [
    "foo",
    "bar"
  ]
}
```

output:

```json
[
  {
    "path": [],
    "value": {"a": true, "b": {"c": 2}, "d": ["foo", "bar"]}
  },
  {
    "path": ["a"],
    "value": true
  },
  {
    "path": ["b"],
    "value": {"c": 2}
  },
  {
    "path": ["b", "c"],
    "value": 2
  },
  {
    "path": ["d"],
    "value": ["foo", "bar"]
  },
  {
    "path": ["d", 0],
    "value": "foo"
  },
  {
    "path": ["d", 1],
    "value": "bar"
  }
]
```
