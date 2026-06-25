build: vet
  go build -o `go env GOPATH`/bin/bark

vet:
  go vet ./...
