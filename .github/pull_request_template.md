## Summary

Describe the change and why it is needed.

## Testing

- [ ] `gofmt -l $(git ls-files '*.go' ':!:testdata/**')`
- [ ] `go vet ./...`
- [ ] `go test ./...`
- [ ] `go build ./cmd/goreview`

## Checklist

- [ ] Documentation updated if behavior changed
- [ ] Tests added or updated for changed behavior
- [ ] Security impact considered
