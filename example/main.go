package main

import (
	"fmt"
	"github.com/gomponents/gontainer/example/pkg"

	"github.com/gomponents/gontainer/example/container"
)

func main() {
	c := container.CreateContainer()
	employee := c.MustGet("employee").(*pkg.Employee)

	fmt.Printf("Full name: %s\n", employee.Fullname)
	fmt.Printf("Position: %s\n", employee.Position)
	fmt.Printf("Salary: %d\n", employee.Salary)
}
