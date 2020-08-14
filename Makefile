clean:
	rm -rf pkg/template2/tpl_*.go

generate-templates: clean
	go generate pkg/template2/template.go

tests: generate-templates
	go test -coverprofile=coverage.out ./cmd/... ./pkg/...
