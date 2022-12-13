
.PHONY: lint run

lint:
	golangci-lint run --out-format=tab

run:
	go run cmd/main.go