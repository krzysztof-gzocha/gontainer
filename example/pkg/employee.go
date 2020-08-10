package pkg

import (
	"fmt"
)

type Employee struct {
	*Person
	Salary   int
	Position string
}

func NewEmployee(person *Person, salary int, position string, aaa ...interface{}) *Employee {
	return &Employee{Person: person, Salary: salary, Position: position}
}

func NewEmplWithErr() (*Employee, error) {
	return nil, fmt.Errorf("my error")
}
