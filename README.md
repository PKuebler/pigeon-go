# Pigeon-Go [![codecov](https://codecov.io/gh/PKuebler/pigeon-go/graph/badge.svg?token=YM26YKAWUJ)](https://codecov.io/gh/PKuebler/pigeon-go)

Pigeon-Go is a Go library for synchronizing JSON documents across multiple machines, designed for collaborative editing with robust conflict resolution. Unlike CRDTs, Pigeon-Go uses operational transforms and identifier-based paths for changes, making debugging easier and change tracking more transparent. All modifications are applied to entire fields rather than supporting live cursor positions, which simplifies the architecture and improves reliability, especially in environments with network delays or out-of-order updates. The system extends standard JSON Patch with enhanced array operations and maintains a complete change history, ensuring consistent document state and easier conflict resolution.

It's working only with the Forked version of [Pigeon JS](https://github.com/spiegeltechlab/pigeon).

## Overview

Pigeon-Go is similar to JSON Patch but provides additional conflict resolution mechanisms. It works with the operations add, replace, remove, and move. When applying changes, the system uses operational transform: it takes the new change with its diff, reverts all changes in the history that occurred after the new change's timestamp, applies the new change, and then reapplies all subsequent changes.

The key innovation is the use of identifier-based paths for array operations instead of fragile index-based references, significantly reducing conflicts in collaborative editing scenarios.

## What is JSON Patch?

JSON Patch is a format for describing changes to a JSON document. It uses operations like:

- `add`: Insert new values
- `replace`: Change existing values  
- `remove`: Delete values
- `move`: Relocate values

Standard JSON Patch uses index-based paths like `/items/0` which can cause conflicts when multiple users edit arrays simultaneously.

## Key Features

### Identifier-Based Array Operations

Instead of using fragile array indices, Pigeon-Go uses object identifiers in square brackets:

```
// Standard JSON Patch (fragile)
"/items/0/name"

// Pigeon-Go (robust)
"/items/[uuid-123]/name"
```

### Working Copy Architecture

Unlike some implementations, Pigeon-Go uses a working copy approach. It creates a clone of the document, applies changes to the clone, and only commits the changes if successful. This prevents document corruption during operations.

### Timebased conflict resolution

The system maintains both the current state and a complete history of changes, enabling proper conflict resolution through operational transform.

Applying a change implies automatically the following steps:

1. **Receive Change**: A new change arrives with a timestamp
2. **Revert History**: All changes after the new change's timestamp are temporarily reverted
3. **Apply Change**: The new change is applied to the document
4. **Reapply History**: All reverted changes are reapplied in chronological order
5. **Conflict Resolution**: The system automatically resolves conflicts using identifier-based paths

## Basic Usage

```go
func main() {
    doc := pigeongo.NewDocument([]byte(`{ "name\": "Philipp" }`))

    // Create a change to replace the name
    change := pigeongo.Change{
        Diff: []pigeongo.Operation{
            {
                Op:    "replace",
                Path:  "/name",
                Value: json.RawMessage("Phil"),
                Prev:  json.RawMessage("Philipp"),
            },
        },
        TimestampMillis: 2,
        ClientID:        "client-1",
        ChangeID:        "change-1",
    }

    err := doc.ApplyChange(change)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println(string(doc.JSON())) // {"name":"Phil"}
}
```

### Working with Arrays

```go
func main() {
    // Document with an array of objects
    doc := pigeongo.NewDocument([]byte(`{
        "users": [
            {"id": "user-1", "name": "Alice"},
            {"id": "user-2", "name": "Bob"}
        ]
    }`))

    // Add a new user at the end using "-"
    change := pigeongo.Change{
        Diff: []pigeongo.Operation{
            {
                Op:    "add",
                Path:  "/users/-",
                Value: json.RawMessage(`{"id": "user-3", "name": "Charlie"}`),
            },
        },
        TimestampMillis: 1,
        ClientID:        "client-1",
        ChangeID:        "change-1",
    }

    doc.ApplyChange(change)

    // Update existing user by ID
    updateChange := pigeongo.Change{
        Diff: []pigeongo.Operation{
            {
                Op:    "replace",
                Path:  "/users/[user-1]/name",
                Value: json.RawMessage("Alice Smith"),
                Prev:  json.RawMessage("Alice"),
            },
        },
        TimestampMillis: 2,
        ClientID:        "client-1",
        ChangeID:        "change-2",
    }

    doc.ApplyChange(updateChange)
}
```

### Custom Identifiers

You can configure custom identifier paths for complex nested structures:

```go
func main() {
    // JSON with nested identifier structure
    jsonData := []byte(`[
        {
            "attrs": {"id": 123},
            "name": "Philipp"
        },
        {
            "uuid": "abc-def",
            "name": "Alice"
        }
    ]`)

    // Configure multiple identifier paths
    doc := pigeongo.NewDocument(jsonData, 
        pigeongo.WithCustomIdentifier([][]string{
            {"id"},           // Look for "id" field first
            {"attrs", "id"},  // Then look for "attrs.id"
            {"uuid"},         // Finally look for "uuid"
        }),
    )

    // Now you can reference items by their identifiers
    change := pigeongo.Change{
        Diff: []pigeongo.Operation{
            {
                Op:    "replace",
                Path:  "/[123]/name",  // References the object with attrs.id = 123
                Value: json.RawMessage("Phil"),
                Prev:  json.RawMessage("Philipp"),
            },
        },
        TimestampMillis: 1,
        ClientID:        "client-1",
        ChangeID:        "change-1",
    }

    doc.ApplyChange(change)
}
```

### Generating Diffs

```go
func main() {
    doc1 := pigeongo.NewDocument([]byte(`{"name": "Alice"}`))
    doc2 := pigeongo.NewDocument([]byte(`{"name": "Bob"}`))

    change, err := doc1.Diff(doc2)
    if err != nil {
        fmt.Printf("Error generating diff: %v\n", err)
        return
    }

    fmt.Printf("Change: %+v\n", change)
}
```

### Cloning Documents

```go
func main() {
    original := pigeongo.NewDocument([]byte(`{"data": "value"}`))
    clone := original.Clone()
    
    // Changes to clone don't affect original
    change := pigeongo.Change{
        Diff: []pigeongo.Operation{
            {
                Op:    "replace",
                Path:  "/data",
                Value: json.RawMessage("new value"),
            },
        },
        TimestampMillis: 1,
        ClientID:        "client-1",
        ChangeID:        "change-1",
    }
    
    clone.ApplyChange(change)
    
    fmt.Println("Original:", string(original.JSON()))
    fmt.Println("Clone:", string(clone.JSON()))
}
```

## Change Structure

```go
type Change struct {
    Diff            []Operation // The actual JSON Patch operations
    ClientID        string      // Who made the change
    ChangeID        string      // Unique ID for this change
    Seq             int         // Sequence number
    TimestampMillis int64       // When change was made
    MessageID       string      // Optional message ID for tracking
}

type Operation struct {
    Op    string          // Operation type: "add", "remove", "replace", "move"
    Path  string          // Target path
    From  string          // Source path (for move operations)
    Value json.RawMessage // New value (for add/replace operations)
    Prev  json.RawMessage // Previous value (internal use for rewind/fast-forward)
}
```

## Internal Architecture

### Operational Transform Process

When applying a change, Pigeon-Go uses a three-phase operational transform process:

- `Rewind`: All changes in history after the new change's timestamp are reversed
- `Apply`: The new change is applied to the rewound state  
- `Fast-forward`: All previously reversed changes are reapplied

This ensures proper conflict resolution and maintains document consistency even with out-of-order change arrival.

### Operation Reversal

#### Add Operation Reversal

- **add** becomes **remove**
- The `Prev` field is cleared (set to `nil`)
- For objects with identifiers, the path is updated from index-based to identifier-based
- Example: `/items/0` becomes `/items/[uuid-123]` where `uuid-123` is the object's ID

#### Remove Operation Reversal

- **remove** becomes **add**
- The Value field is cleared (set to `nil`)
- Identifier-based paths are converted back to index `0` for reinsertion at the beginning
- Example: `/items/[uuid-123]` becomes `/items/0`

#### Value and Prev Field Handling

During reversal, the `Value` and `Prev` fields are swapped:

- The original `Value` becomes the new `Prev`
- The original `Prev` becomes the new `Value`
- `nil` values are preserved as `nil`

### Position Preservation

The reversal process is critical for maintaining correct positions during rewind/fast-forward:

1. **During Rewind**: Operations are reversed to undo their effects
2. **During Fast-forward**: The same operations are applied again in their original form
3. **Prev Value Management**: External Prev values are cleared during apply to ensure they reference the local state at the specific point in history, preventing conflicts from asynchronous operations

This architecture ensures that even if changes arrive out of order or are applied asynchronously, the document state remains consistent and conflicts are properly resolved.

# Limitations

- Identifier-based paths require objects in arrays to have identifiable fields
- The system maintains full history, which may consume memory for long-lived documents
- Complex nested structures may require careful identifier configuration

# Differences to the Javascript version

- With Changes it is possible to use a `msg_id` in addition to the `change_id`. For example, it is also possible to transport a Kafka, Redis or network protocol ID.
- If a command fails, the document remains in its initial state. A broken state is possible in PigeonJS.

# Contributing

Contributions are welcome. Please ensure all tests pass and follow Go conventions.

# License

See LICENSE file for details.