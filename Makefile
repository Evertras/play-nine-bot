.PHONY: run
run:
	@go run ./cmd/playnine

.PHONY: test
test:
	go test ./...
