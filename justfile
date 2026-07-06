build: vet
  go build -buildvcs=false -o `go env GOPATH`/bin/bark

vet:
  go vet ./...
