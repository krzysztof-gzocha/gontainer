package main

import (
	"fmt"
)

type Foo struct {
}

func (Foo) Bar(a string, b string) string {
	fmt.Printf("Input %s, %s\n", a, b)
	return "qwertyyyy"
}

func test(a string, b ...int) string {
	fmt.Printf("Input %v, %v\n", a, b)
	return "qwerty"
}

func main() {
	//defer func() {
	//	if v := recover(); v != nil {
	//		fmt.Printf("%v\n", v)
	//	}
	//}()
	//result := reflect.ValueOf((&Foo{}).Bar).Call([]reflect.Value{
	//	reflect.ValueOf("aaaa"),
	//	reflect.ValueOf("bbb"),
	//})
	//fmt.Printf("Function test returned `%v`\n", result[0].Interface())
	//a, b := pkg.SafeCall(Foo{}.Bar, "qqq", "www", "qqq")
	//fmt.Printf("%#v %s\n", a, b)
	var a []int
	fmt.Printf("%#v", a)
}
