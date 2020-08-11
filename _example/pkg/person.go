package pkg

type Person struct {
	Fullname string
	Age      uint
}

func NewPerson(fullname string, age uint) *Person {
	return &Person{Fullname: fullname, Age: age}
}
