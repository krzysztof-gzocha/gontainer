clean:
	rm -rf pkg/template/tmpl_*.go

templates: clean
	go generate pkg/template/template.go

tests: templates
	go test -coverprofile=coverage.out ./cmd/... ./pkg/...
