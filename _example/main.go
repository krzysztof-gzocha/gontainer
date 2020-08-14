package main

import (
	"fmt"
	"mymodule/container"
)

func main() {
	c := container.CreateContainer()
	fmt.Println(c.GetParam("es_host_port"))
	c = container.CreateContainer()
	// does not work fix it
	c.OverrideParam("es_port", func() interface{} {
		return 9300
	})
	fmt.Println(c.GetParam("es_host_port"))
	return

	//c := container.CreateContainer()
	////employee := c.MustGet("employee").(*pkg.Employee)
	//employee, err := c.GetEmployee()
	//fmt.Printf("err %s\n", err)
	//return
	//p, _ := c.Get("testPerson")
	//fmt.Printf("%#v", p)
	//q, w := p.(*pkg.Employee)
	//fmt.Printf("%#v %#v\n", q, w)
	////return
	//
	//fmt.Printf("Full name: %s\n", employee.Fullname)
	//fmt.Printf("Position: %s\n", employee.Position)
	//fmt.Printf("Salary: %d\n", employee.Salary)
	//
	//a, b := c.GetMyErr()
	//fmt.Printf("%#v %#v\n", a, b)
	//
	//p1, _ := c.GetPerson()
	//p2, _ := c.GetPerson()
	//e1, _ := c.GetEmployee()
	//e2, _ := c.GetEmployee()
	//
	//fmt.Printf("person is disposable, p1 == p2 %v\n", p1 == p2)
	//fmt.Printf("employee isn't disposable, e1 == e2 %v\n", e1 == e2)
}
