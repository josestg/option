# option

Option is a generic type for Go to represent a presence or absence of a value as an alternative for `nil` values.

## Install

```bash
go get github.com/josestg/option
```


## Usage

```go
package main

import (
    "fmt"

    "github.com/josestg/option"
)

func main() {
    none := div(1, 0)
    fmt.Println("String()       ", none)                             // "None"
    fmt.Println("ValueOr()      ", none.ValueOr(42))                 // 42
    fmt.Println("ValueOrBy()    ", none.ValueOrBy(fallbackSupplier)) // 21

    some := div(4, 2)
    fmt.Println("String()       ", some)                             // "Some(2)"
    fmt.Println("ValueOr()      ", some.ValueOr(42))                 // 2
    fmt.Println("ValueOrBy()    ", some.ValueOrBy(fallbackSupplier)) // 2
}

func fallbackSupplier() int {
    return 3 * 7
}

func div(a, b int) option.Option[int] {
    if b == 0 {
        return option.None[int]()
    }
    return option.Some(a / b)
}

```