clean:
	rm -rf pkg/template2/tpl_*.go

templates: clean
	go generate pkg/template2/template.go

tests: templates
	go test -coverprofile=coverage.out ./cmd/... ./pkg/...
