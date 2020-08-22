package main

import (
	"fmt"
	"mymodule/container"
)

func main() {
	c := container.CreateContainer()
	doer, err := c.GetDoer()
	fmt.Printf("Err %s", err)
	doer.Do()
}
