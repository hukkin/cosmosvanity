## Running tests
```bash
go test
```

## Installing pre-commit hooks
```bash
pre-commit install
```

## Running linters
```bash
pre-commit run -a
```

## Releasing
Run the following to tag a release and publish binaries on Github:
```bash
bump2version (major | minor | patch)
git push --follow-tags origin master
GITHUB_TOKEN=<secret-token-from-github> goreleaser
```
