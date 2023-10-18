#  Pigeon-Go

A JSON patch module to exchange data compatible with the JS Package [Pigeon](https://github.com/frameable/pigeon).

## Production ready?

No, this is a very early unfinished version. Please use with caution.

```golang
package main

import (
    "fmt"
    "github.com/pkuebler/pigeon-go"
)

func main() {
    doc := NewDocument([]byte(`{ "name": "Philipp" }`))

    doc.ApplyChanges(Changes{
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
}
```