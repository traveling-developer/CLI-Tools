# web2md

A CLI tool to convert a web page to Markdown. Built on top of [defuddle](https://github.com/kepano/defuddle), whose platform-specific binaries are embedded directly into the Go binary — no runtime dependency.

## Usage

```
web2md --url <url> [options]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--url` | `-u` | — | URL of the web page to convert **(required)** |
| `--output` | `-o` | stdout | Write output to file |
| `--help` | `-h` | | Show help message |

### Examples

```bash
# Print Markdown to stdout
web2md --url https://example.com

# Write Markdown to a file
web2md --url https://example.com --output article.md
web2md -u https://example.com -o article.md
```

## How it works

The defuddle binaries for linux, darwin, and windows (amd64 + arm64) are bundled via `//go:embed`. At runtime, `web2md` picks the binary matching `runtime.GOOS`/`runtime.GOARCH`, extracts it to a temp file (e.g. `/var/folders/.../T/defuddle-*` on macOS), and invokes it as a subprocess. The temp file is removed on exit.

On macOS, the extracted binary is ad-hoc re-signed with `codesign` so Apple Silicon accepts it. If `codesign` is unavailable or fails, `web2md` falls back to bun's pre-applied signature.

## Building the defuddle binaries

The defuddle binaries are not checked into the repo beyond what `//go:embed` needs. Rebuild them via the provided Makefile, which uses the Docker image defined in `defuddle/dockerfile`:

```bash
# From tools/web2md/
make defuddle     # Build all platform binaries into defuddle/
make build        # Build web2md for the current host
make build-all    # Cross-compile web2md for all supported platforms
make clean        # Remove built artifacts
```

Supported target matrix:

| OS | Arch |
|----|------|
| linux | amd64, arm64 |
| darwin | amd64, arm64 |
| windows | amd64 |
