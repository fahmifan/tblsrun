.PHONY: dbdoc-foo
dbdoc-foo:
	@rm -rf dbdoc/foo
	@echo "generate schema example"
	go run ../cmd/tblsrun/tblsrun.go postgres docker --env-file .env.example.foo

dbdoc-bar:
	@rm -rf dbdoc/bar
	@echo "generate schema example"
	go run ../cmd/tblsrun/tblsrun.go postgres docker --env-file .env.example.bar

.PHONY: dbdoc
dbdoc: dbdoc-foo dbdoc-bar
