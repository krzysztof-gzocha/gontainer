# Gontainer

Depenendency Injection container for GO inspired by Symfony.

## Example

```
meta:
  pkg: container

parameters:
  first_name: "Jane"
  last_name: "Doe"
  age: 30
  salary: 30000
  position: "CTO"

services:
  person:
    constructor: "github.com/gomponents/gontainer/example/pkg.NewPerson"
    args: ["%first_name% %last_name%", "%age%"]

  employee:
    constructor: "github.com/gomponents/gontainer/example/pkg.NewEmployee"
    args:
      - "@person.(*github.com/gomponents/gontainer/example/pkg.Person)"
      - "%salary%"
      - "%position%"
```
