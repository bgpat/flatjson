# flatjson

[![Go Reference](https://pkg.go.dev/badge/github.com/bgpat/flatjson.svg)](https://pkg.go.dev/github.com/bgpat/flatjson)

`flatjson` is a Go package that provides utilities for flattening nested JSON structures into a sequence of path-value pairs.
This makes it easier to compare differences between JSON objects, which is useful for various purposes such as data synchronization, version control, and data analysis.

## Installation

```sh
go install github.com/bgpat/flatjson/cmd/flatjson
```

## Example

### Flatten

```console
$ cat << EOF > input.json
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
EOF
$ flatjson < input.json
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

### Diff

```console
$ jq 'del(.a) | .d[1] = "buzz" | .e = null' input.json | flatjson -diff input.json -
[
  {
    "type": "replace",
    "path": [],
    "value": {"b": {"c": 2}, "d": ["foo", "buzz"], "e": null},
    "old_value": {"a": true, "b": {"c": 2}, "d": ["foo", "bar"]}
  },
  {
    "type": "remove",
    "path": ["a"]
  },
  {
    "type": "replace",
    "path": ["d"],
    "value": ["foo", "buzz"],
    "old_value": ["foo", "bar"]
  },
  {
    "type": "replace",
    "path": ["d", 1],
    "value": "buzz",
    "old_value": "bar"
  },
  {
    "type": "add",
    "path": ["e"],
    "value": null
  }
]
```
