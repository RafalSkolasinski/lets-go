run:
	@go run ./cmd/web

open:
	xdg-open http://localhost:4000

test:
	go test -v ./...
