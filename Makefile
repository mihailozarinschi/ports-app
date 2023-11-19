test:
	go test -mod=vendor -count=1 -race -cover -short ./...

lint:
	go run -mod=vendor github.com/golangci/golangci-lint/cmd/golangci-Lint run --timeout 2m

run:
	docker-compose up portsd
