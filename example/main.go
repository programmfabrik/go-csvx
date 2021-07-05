package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/programmfabrik/go-csvx"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file, err := os.Open("untyped.csv")
	panicIfError(err)
	defer file.Close()

	bts, err := ioutil.ReadAll(file)
	panicIfError(err)

	data, err := csvx.NewCSV(',', '#', true).ToMap(bts)
	panicIfError(err)

	fmt.Printf("untyped data:\n\t%+v\n", data)

	file2, err := os.Open("typed.csv")
	panicIfError(err)
	defer file.Close()

	bts2, err := ioutil.ReadAll(file2)
	panicIfError(err)

	data2, err := csvx.NewCSV(',', '#', true).ToTypedMap(bts2)
	panicIfError(err)

	fmt.Printf("typed data2:\n\t%+v\n", data2)

	fmt.Printf("check reflect type:\n\tfirst: %v\n\tsecond: %v\n\tthird: %v\n\tfourth: %v\n",
		reflect.TypeOf(data2[0]["foo"]),
		reflect.TypeOf(data2[0]["bar"]),
		reflect.TypeOf(data2[0]["counter"]),
		reflect.TypeOf(data2[0]["names"]))
}
