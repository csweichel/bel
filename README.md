# [bel](https://en.wikipedia.org/wiki/Bel_(mythology))
Generate TypeScript interfaces from Go structs/interfaces - useful for JSON RPC

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#github.com/32leaves/bel)

## Getting started
`bel` is super easy to use. There are two steps involved: extract the Typescript information, and generate the Typescript code.
```Go
package main

import (
    "github.com/32leaves/bel"
)

type Demo struct {
    Foo string `json:"foo,omitempty"`
    Bar uint32
    Baz struct {
        FirstField  bool
        SecondField *string
    }
}

func main() {
    ts, err := Extract(Demo{})
    if err != nil {
        panic(err)
    }

    err = bel.Render(ts)
    if err != nil {
        panic(err)
    }
}
```

produces something akin to (sans formatting):

```TypeScript
export interface Demo {
    foo?: string
    Bar: number
    Baz: {
        FirstField: boolean
        SecondField: string
    }
}
```

### Converting interfaces
You can also convert Golang interfaces to TypeScript interfaces. This is particularly handy for JSON RPC:
```Go
package main

import (
    "os"
    "github.com/32leaves/bel"
)

type DemoService interface {
    SayHello(name, msg string) (string, error)
}

func main() {
    ts, err := Extract((*DemoService)(nil))
    if err != nil {
        panic(err)
    }

    err = bel.Render(ts)
    if err != nil {
        panic(err)
    }
}
```

produces something akin to (sans formatting):

```TypeScript
export interface DemoService {
    SayHello(arg0: string, arg1: string): string
}
```

## Advanced Usage
You can try all the examples mentioned below in [Gitpod](https://gitpod.io#github.com/32leaves/bel).

### FollowStructs
Follow structs enable the transitive generation of types. See See [examples/embed-structs.go](examples/follow-structs.go).

would produce code for the interface `UserService`, as well as the struct it refers to `AddUserRequest`, and `User` because it's referenced by `AddUserRequest`.
Without `FollowStructs` we'd simply refer to the types by name, but would not generate code for them.

### EmbedStructs
Embed structs is similar to `FollowStructs` except that it produces a single canonical type for each structure.
Whenever one struct references another, that reference is resolved and the definition of the other is embedded.
See [examples/embed-structs.go](examples/embed-structs.go).

### NameAnonStructs
`NameAnonStructs` is kind of the opposite of `EmbedStructs`. When we encounter a nested anonymous struct, we make give this previously anonymous structure a name and refer to it using this name.
See [examples/name-anon-structs.go](examples/name-anon-structs.go).

### CustomNamer
`CustomNamer` enables full control over the TypeScript type names. This is handy to enforce a custom coding guideline, or to add a prefix/suffix to the generated type names.
See [examples/custom-namer.go](examples/custom-namer.go).

### Enums
See [examples/enums.go](examples/enums.go).

Go famously does not have enums, but rather type aliases and consts. Using reflection alone there is no way to obtain a comprehensive list of type values, as the linker might optimize and remove some.
_bel_ supports the extraction of enums by parsing the Go source code. Note that this is merely a heuristic and may fail in your case. If it does not work, _bel_ falls back to the underlying type.

Enums can be generated as TypeScript `enum` or as sum types. Use the `bel.GenerateEnumsAsSumTypes` flag to change this behaviour.

### Code Generation
See [examples/code-generation.go](examples/code-generation.go).

When generating the TypeScript code you might want to wrap everything in a namespace. To that end `bel.GenerateNamespace("myNamespaceName")` can be used.

By default _bel_ adds a comment to the files it generates. You can influcence this comment (and any other code that comes before the generated code)
using `bel.GeneratePreamble` and `bel.GenerateAdditionalPreamble`.

You can configure the `io.Writer` that _bel_ uses using `bel.GenerateOutputTo`.

# Contributing
All contributions/PR/issue/beer are welcome ❤️.

It's easiest to work with _bel_ using Gitpod: [![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#github.com/32leaves/bel)