* build using multiple of different yaml files (build -i container/*.yml - container/prod/*.yml ...)
* tests
* getters: c.MustGetEmployee() instead of c.MustGet("employee").(*pkg.Employee)
* go generate for templates
* fix the following error: cannot create service person: cannot create person due to: cannot create service wallet: cannot create wallet due to: service serviceContainer does not exist
* rename constructor to provider, allow for syntax `@serviceContainer.CreateSth` ???
* allow injecting custom functions for parameters, e.g.:

```yaml
meta:
    functions:
      name: "env"
      callee: "os.Getenv"

params:
    name: '%env("NAME")%'
```

* real time params

```yaml
parameters:
  name: Jane
  first_name: "%name%"
```

**Now**
```go
params["first_name"] = "Jane"
```
**TODO**
```go
params["first_name"] = func () interface{} {
    return container.GetParam("name")
}
```

* `ValidateAllServices` 

replace

```
	for _, id := range []string{
		"anotherTest",
		"employee",
```

by

```
    for _, id := range.GetAllServices
```

* flag `todo` in services
