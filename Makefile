.PHONY: test
test:
	go test ./...


.PHONY: fmt
fmt:
	go fmt ./...


.PHONY: tidy
tidy:
	go mod tidy
