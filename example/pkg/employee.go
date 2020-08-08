package pkg

type Employee struct {
	*Person
	Salary   int
	Position string
}

func NewEmployee(person *Person, salary int, position string) *Employee {
	return &Employee{Person: person, Salary: salary, Position: position}
}
