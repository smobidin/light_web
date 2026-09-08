# Light Web Server - Go Version

This is a Go implementation of the Light Web Server, which provides the same functionality as the original Python version:

- File server with Markdown rendering
- Syntax highlighting for source code files
- Directory browsing
- PDF viewing with integrated PDF.js
- LaTeX math processing (server-side)

## Building

### Prerequisites

- [Go](https://go.dev/dl/) 1.21 or newer
- Internet access on first build (to fetch module dependencies)

### Build

The project is a proper Go module (`gosrc/go.mod`), so you can build it with:

```bash
cd gosrc
go build -o lightweb .
```

This builds the whole package (all `.go` files including `theme.go`) and creates a single binary `lightweb` with no external runtime dependencies. The binary is written to `gosrc/lightweb`.

Alternatively, to place the binary elsewhere:

```bash
cd gosrc
go build -o ../lightweb .
```

### Cross-compilation

You can cross-compile for other platforms (the binary is pure Go with no CGo):

```bash
GOOS=linux GOARCH=amd64 go build -o lightweb .   # Linux x86_64
GOOS=darwin GOARCH=arm64 go build -o lightweb .  # macOS ARM
GOOS=windows GOARCH=amd64 go build -o lightweb.exe .
```

### Verify

```bash
go vet ./...
```

## Running

The Go version supports the same command-line arguments as the Python version:

```bash
./lightweb --host localhost --port 8081 --directory /path/to/files
```

- `--host`: Host to bind to (default: all interfaces)
- `--port`: Port to listen on (default: 8080)
- `--directory`: Directory to serve (default: current directory)

## Features

- Single binary deployment with no runtime dependencies
- Improved performance compared to the Python version
- Full compatibility with existing functionality
- Server-side LaTeX math processing
- Same HTML output and styling as the Python version

## Implementation Details

The Go version consists of these main components:

1. `main.go` - Entry point and command-line argument parsing
2. `handler.go` - HTTP request handling and routing
3. `renderer.go` - Markdown and source code rendering
4. `directory.go` - Directory listing generation
5. `pdf.go` - PDF viewer HTML generation
6. `utils.go` - Utility functions including LaTeX math processing
7. `theme.go` - Light/dark theme management and toggle script

External dependencies are embedded in the binary during compilation:
- `github.com/gomarkdown/markdown` for Markdown processing
- `github.com/alecthomas/chroma/v2` for syntax highlighting