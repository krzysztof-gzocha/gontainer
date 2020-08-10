# TEST REPOSITORY

[![Build Status](https://github.com/gomponents/gontainer/workflows/Tests/badge.svg?branch=master)](https://github.com/gomponents/gontainer/actions?query=workflow%3AGo)
[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer?branch=master)

# Gontainer

Depenendency Injection container for GO inspired by Symfony.

## Example

```yaml
meta:
  pkg: container
  container_type: MyContainer
  imports:
    gontainer: "github.com/gomponents/gontainer"

parameters:
  first_name: '%env("NAME")%'
  last_name: "Doe"
  age: '%envInt("AGE")%'
  salary: 30000
  position: "CTO"

services:
  personExample1:
    type: "*gontainer/example/pkg/Employee" # alias.Employee{}

  personExample2:
    type: "gontainer/example/pkg/Employee" # &alias.Employee{}

  person:
    constructor: "github.com/gomponents/gontainer/example/pkg.NewPerson"
    args: ["%first_name% %last_name%", "%age%"]

  employee:
    getter: "Employee"
    type: "*gontainer/example/pkg/Employee"
    constructor: "gontainer/example/pkg.NewEmployee"
    args:
      - "@person.(*gontainer/example/pkg.Person)"
      - "%salary%"
      - "%position%"
```
