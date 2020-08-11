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