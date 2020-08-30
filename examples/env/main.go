package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	c := NewContainer()
	p, _ := c.GetPerson()
	fmt.Printf("%s is %d years old\n", p.Name, p.Age)
}
