# go-csvx

[![unit-tests](https://github.com/programmfabrik/go-csvx/actions/workflows/tests.yml/badge.svg)](https://github.com/programmfabrik/go-csvx/actions/workflows/tests.yml)

go-csvx implements the csv encoder and extends it with various functionalities. This library supports a typed csv format. The second column defines the data type.

## Getting started

Since this repository serves as a library for typed and untyped csv parsing, you need to define it as a Go dependency to work with it. To declare it as a dependency in your project, use the **go get** command from below to attach it to your project:

```bash
go get github.com/programmfabrik/go-csvx
```

## Defining a csv

**Typed:**

```csv
foo,bar,counter,names
string,string,int,"string,array"
test1,test2,10,"hello,world,how,is,it,going"
```

**Untyped:**

```csv
foo,bar,counter,names
hello,world,10,"hello,world,how,is,it,going"
```

## Supported datatypes

This library supports the following list of data types:

- string
- int64
- int
- float64
- bool
- "string,array"
- "int64,array"
- "float64,array"
- "bool,array"
- json

## Example

**Untyped example:**

```go
package main

import (
    "fmt"

    "github.com/programmfabrik/go-csvx"
)

func main() {
    data, _ := csvx.NewCSV(',', '#', true).ToMap([]byte(
        `foo,bar,counter,names\n
        hello,world,10,"hello,world,how,is,it,going"`))

    fmt.Printf("untyped data:\n\t%+#v\n", data)
}

```

Result:

```txt
untyped data:
        []map[string]interface {}{map[string]interface {}{"bar":"world", "counter":"10", "foo":"hello", "names":"hello,world,how,is,it,going"}}
```

**Typed example:**

```go
package main

import (
    "fmt"

    "github.com/programmfabrik/go-csvx"
)

func main() {
    data2, _ := csvx.NewCSV(',', '#', true).ToTypedMap([]byte(
        `foo,bar,counter,names
        string,string,int,"string,array"
        hello,world,10,"hello,world,how,is,it,going"`))

    fmt.Printf("typed data:\n\t%+#v\n", data2)
}

```

Result:

```txt
typed data:
        []map[string]interface {}{map[string]interface {}{"bar":"world", "counter":10, "foo":"hello", "names":[]string{"hello", "world", "how", "is", "it", "going"}}}
```
