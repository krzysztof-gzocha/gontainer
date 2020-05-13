package main

import (
	"fmt"

	"github.com/gomponents/gontainer/example/container"
)

func main() {
	c := container.CreateContainer()
	person, _ := c.Get("person")
	fmt.Printf("%#v\n", person)
}

