.Phony: run

run: go.sum
	go run ./cmd/compress/main.go

go.sum: go.mod
	@echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy