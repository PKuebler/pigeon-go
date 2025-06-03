#  Pigeon-Go [![codecov](https://codecov.io/gh/PKuebler/pigeon-go/graph/badge.svg?token=YM26YKAWUJ)](https://codecov.io/gh/PKuebler/pigeon-go)

A JSON patch module to exchange data compatible only with the forked JS Package [Pigeon](https://github.com/spiegeltechlab/pigeon).

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

    err := doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Phil"`),
				Prev:  rawMessage(`"Philipp"`),
			},
		},
		TimestampMillis:  2,
		ClientID: "50reifj9hyt",
		ChangeID: "dva96nqsdd",
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

- With Changes it is possible to use a `msg_id` in addition to the `change_id`. For example, it is also possible to transport a Kafka, Redis or network protocol ID.
- If a command fails, the document remains in its initial state. A broken state is possible in PigeonJS.