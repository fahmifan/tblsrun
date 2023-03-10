dbdoc:
	@rm -rf example/dbdoc
	@echo "generate schema example"
	go run cmd/tblsrun/tblsrun.go postgres docker