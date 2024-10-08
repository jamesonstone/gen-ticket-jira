run:
	go run cmd/genticketjira/main.go
.PHONY: run

build:
	go build -o ./bin/genticketjira cmd/genticketjira/main.go
.PHONY: build
