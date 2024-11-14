package main

import (
	"fmt"
	"reflect"
)

type Email = string
type Email1 string

func (e Email1) Validate() bool {
	return true
}
func Do() {
	var e Email = "email"
	var ee Email1 = "email"
	ee.Validate()
	fmt.Println(reflect.TypeOf(e))
	fmt.Println(reflect.TypeOf(ee))
}

func main() {
	Do()
}
