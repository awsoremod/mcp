VERSION := "dev"

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=$$(git rev-parse HEAD) -X main.date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
	-o mcp ./cmd/mcp

run:
	./mcp stdio
