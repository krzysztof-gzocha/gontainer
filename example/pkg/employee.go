package pkg

type Employee struct {
	*Person
	Salary   uint
	Position string
}

func NewEmployee(person *Person, salary uint, position string) *Employee {
	return &Employee{Person: person, Salary: salary, Position: position}
}
