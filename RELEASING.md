# Release Process

This document describes how to create a new release of slka.

## Prerequisites

1. Install GoReleaser:
   ```bash
   # macOS
   brew install goreleaser

   # Linux
   go install github.com/goreleaser/goreleaser/v2@latest

   # Or download from https://github.com/goreleaser/goreleaser/releases
   ```

2. Ensure you have push access to the GitHub repository

## Testing the Release Process

Before creating an actual release, test the configuration locally:

```bash
# Test the release process without publishing
goreleaser release --snapshot --clean --skip=publish

# This will:
# - Run tests
# - Build binaries for all platforms
# - Create archives
# - Generate checksums
# - Create a local dist/ directory with all artifacts
```

Check the `dist/` directory to verify:
- All binaries are built correctly
- Archives contain the right files
- Checksums are generated

## Creating a Release

### 1. Prepare the Release

1. Ensure all changes are committed and pushed to main
2. Update version numbers if needed
3. Run tests: `go test ./...`

### 2. Create and Push a Tag

```bash
# Create a new tag (use semantic versioning)
git tag -a v0.3.0 -m "Release v0.3.0: Add reaction commands"

# Push the tag to GitHub
git push origin v0.3.0
```

### 3. Automated Release

Once you push the tag, GitHub Actions will automatically:
1. Checkout the code
2. Set up Go
3. Run GoReleaser
4. Build binaries for all platforms (Linux, macOS, Windows; AMD64, ARM64)
5. Create archives with documentation
6. Generate checksums
7. Create a GitHub release with all artifacts
8. Generate a changelog

### 4. Verify the Release

1. Go to https://github.com/ulfschnabel/slka/releases
2. Verify the new release is published
3. Check that all binaries are attached
4. Review the generated changelog
5. Test download and run a binary

## Release Artifacts

Each release includes:

### Binaries
- `slka-{version}-linux-amd64.tar.gz`
- `slka-{version}-linux-arm64.tar.gz`
- `slka-{version}-darwin-amd64.tar.gz` (macOS Intel)
- `slka-{version}-darwin-arm64.tar.gz` (macOS Apple Silicon)
- `slka-{version}-windows-amd64.zip`
- `slka-{version}-windows-arm64.zip`

### Documentation (included in each archive)
- README.md
- QUICKSTART.md
- USER_TOKEN_SETUP.md
- MANIFEST_SETUP.md
- slack-manifest-user-token.yaml
- slack-manifest-bot-token.yaml

### Checksums
- checksums.txt - SHA256 checksums for all artifacts

## Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backward-compatible manner
- **PATCH** version for backward-compatible bug fixes

Examples:
- `v0.1.0` - Initial release with basic functionality
- `v0.2.0` - Added unified binary (breaking change from split binaries)
- `v0.3.0` - Added reaction commands (new feature)
- `v0.3.1` - Fixed bug in reaction handling (bug fix)

## Rollback

If you need to remove a release:

1. Delete the release from GitHub UI (Releases → Edit → Delete release)
2. Delete the tag locally: `git tag -d v0.3.0`
3. Delete the tag remotely: `git push origin :refs/tags/v0.3.0`

## Troubleshooting

### GoReleaser fails locally

```bash
# Check configuration is valid
goreleaser check

# Show more details during build
goreleaser release --snapshot --clean --skip=publish --verbose
```

### GitHub Actions fails

1. Check the Actions tab: https://github.com/ulfschnabel/slka/actions
2. View the workflow logs for errors
3. Common issues:
   - Tests failing: Fix tests and push a new commit
   - Build failures: Check Go version compatibility
   - Permission issues: Verify GITHUB_TOKEN permissions

### Missing artifacts in release

- Check `.goreleaser.yaml` archives section
- Ensure files exist at build time
- Verify file patterns match correctly

## Local Development

For local development without publishing:

```bash
# Build for current platform only
go build -ldflags "-X main.Version=dev" ./cmd/slka

# Build for all platforms manually (without GoReleaser)
GOOS=linux GOARCH=amd64 go build -o dist/slka-linux-amd64 ./cmd/slka
GOOS=darwin GOARCH=arm64 go build -o dist/slka-darwin-arm64 ./cmd/slka
GOOS=windows GOARCH=amd64 go build -o dist/slka-windows-amd64.exe ./cmd/slka
```

## Notes

- GoReleaser automatically generates a changelog from commit messages
- Use conventional commit messages for better changelog organization:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `enhance:` for improvements
- The version is injected at build time via ldflags
- Archives include all essential documentation files
