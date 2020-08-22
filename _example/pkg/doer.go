package pkg

import (
	"log"
)

type Doer struct {
	logger *log.Logger
	name   string
	age    int
}

func NewDoer(logger *log.Logger) *Doer {
	return &Doer{logger: logger}
}

func (d *Doer) SetAge(age int) {
	d.age = age
}

func (d Doer) Do() {
	d.logger.Println("Start")
	d.logger.Printf("My name is %s\n", d.name)
	d.logger.Printf("I'm %d years old", d.age)
	d.logger.Println("End")
}
