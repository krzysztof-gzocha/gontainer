package main

import "fmt"

func main() {
	c := CreateContainer()
	p, e := c.GetPerson()
	fmt.Printf("%#v %#v\n", p, e)
}
