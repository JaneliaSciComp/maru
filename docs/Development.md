# Maru Development

## Building

To compile and install Maru into your standard Go bin directory:
```
go install
```

Any time the templates change, the serialization needs to be updated as follows:
```
go generate ./...
```

## Testing

TBD

## Releasing

New releases are built and deployed to GitHub using GoReleaser. You first need to install GoReleaser and configure your GitHub token as per the [quickstart instructions](https://goreleaser.com/quick-start/). Then you can tag and release a new version as follows:

1. Update the MaruVersion constant in utils/utils.go
2. Tag and release:
```
git tag -a 0.1.0 -m "Release 0.1.0"
git push origin 0.1.0
goreleaser --rm-dist
```

To test a SNAPSHOT release without tagging the code:
```
goreleaser --snapshot --skip-publish --rm-dist
```

