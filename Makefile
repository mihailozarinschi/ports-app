test:
	go test -mod=vendor -count=1 -race -cover -short ./...

run:
	go run ./cmd/portsd/main.go
