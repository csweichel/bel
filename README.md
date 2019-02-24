# [bel](https://en.wikipedia.org/wiki/Bel_(mythology))
Generate TypeScript interfaces from Go structs/interfaces - useful for JSON RPC

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#github.com/32leaves/bel)

## Getting started
`bel` is super easy to use. There are two steps involved: extract the Typescript information, and generate the Typescript code.
```Go
package main

import (
    "os"
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
