package main

import (
	"fmt"
	"reflect"
)

func print_val(x any) {
	val := reflect.ValueOf(x)
	
	fmt.Printf("string: %v\n", val.String())
	fmt.Printf("type: %v\n", val.Type())
	fmt.Printf("val: %v\n", val)
	fmt.Printf("kind: %v\n", val.Kind())
	fmt.Printf("elem: %v\n", val.Elem())
	fmt.Printf("elem.type: %v\n", val.Elem().Type())
	fmt.Printf("elem.kind: %v\n", val.Elem().Kind())
	
	if(val.Elem().Kind() == reflect.Slice) {
		fmt.Println()
		fmt.Printf("elem.type.elem.kind: %v\n", val.Elem().Type().Elem().Kind())
	}
	
	fmt.Println()
	fmt.Printf("can set: %v\n", val.CanSet())
	fmt.Printf("elem can set: %v\n", val.Elem().CanSet())
	fmt.Printf("can int: %v\n", val.CanInt())
	fmt.Printf("elem can int: %v\n", val.Elem().CanInt())
	
	//val.Elem().SetInt(20)
}

func main() {
	y := []byte("abc")
	print_val(&y)
	
	fmt.Printf("\ny: %v\n", y)
}
