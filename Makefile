clean:
	rm -rf pkg/template/tmpl_*.go
	rm -rf app.bin

templates: clean
	go generate pkg/template/template.go

tests-unit: templates
	go test -coverprofile=coverage.out ./cmd/... ./pkg/...

code-coverage:
	go tool cover -func=coverage.out

build: clean templates
	go build -v -o app.bin main.go

run-example-library: build
	./app.bin build -i examples/library/container/gontainer.yaml -o examples/library/container/container.go
	cd examples/library && go run main.go
