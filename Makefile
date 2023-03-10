dbdoc:
	@rm -rf example/dbdoc
	@echo "generate schema example"
	go run cmd/tblsrun/tblsrun.go postgres docker
	@# @echo "generate schema example 2"
	@# go run cmd/tblsrun/tblsrun.go --env-file=.env.example_2 postgres docker