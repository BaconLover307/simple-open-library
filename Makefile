run:
	cd cmd; go run main.go

build:
	cd cmd; go build

build_run:
	cd cmd; go build; ./cmd

run_test:
	go test -v -cover ./repository/...; go test -v -cover ./service/...; go test -v -cover ./controller/...