build:
	go build -o bin/gomodsearch cmd/gomodsearch/gomodsearch.go
install:
	go install cmd/gomodsearch/gomodsearch.go
run:
	go run cmd/gomodsearch/gomodsearch.go $@