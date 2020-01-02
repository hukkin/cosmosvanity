## Running tests
```bash
go test
```

## Running linters
```bash
golangci-lint run
```

## Releasing
Run the following to tag a release and publish binaries on Github:
```bash
bump2version (major | minor | patch)
git push --follow-tags origin master
GITHUB_TOKEN=<secret-token-from-github> goreleaser
```
