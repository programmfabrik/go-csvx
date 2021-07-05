# go-csvx

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
