#  Pigeon-Go [![codecov](https://codecov.io/gh/PKuebler/pigeon-go/graph/badge.svg?token=YM26YKAWUJ)](https://codecov.io/gh/PKuebler/pigeon-go)

A JSON patch module to exchange data compatible with the JS Package [Pigeon](https://github.com/frameable/pigeon).

## Production ready?

No, this is a very early unfinished version. Please use with caution.

## Example

```golang
package main

import (
    "fmt"
    pigeongo "github.com/pkuebler/pigeon-go"
)

func main() {
    doc := pigeongo.NewDocument([]byte(`{ "name": "Philipp" }`))

    err := doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Phil"`),
				Prev:  rawMessage(`"Philipp"`),
			},
		},
		Ts:  2,
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
	})

    // Print JSON
    fmt.Println(doc.JSON)
    // Print Warnings
    fmt.Println(err)
}
```

## Custom Identifier

```golang
pigeongo.NewDocument([]byte(`[{ "attrs": { "id": 123 }, "name": "Philipp" }]`), pigeongo.WithCustomIdentifier([][]string{{"id"},{"attrs", "id"}}))
```

## Differences to the Javascript version

- With Changes it is possible to use a `mid` in addition to the `gid`. For example, it is also possible to transport a Kafka, Redis or network protocol ID.
- If a command fails, the document remains in its initial state. A broken state is possible in PigeonJS.