# Light Web Server - Agent Guide

One-file Python web server for rendering Markdown files with syntax highlighting, LaTeX math, and a directory browser.

## Commands

```bash
# Install dependencies (venv recommended due to PEP 668)
python3 -m venv .venv && .venv/bin/pip install -r requirements.txt

# Run server (via launcher script — shows all accessible IPs)
~/bin/light_web_run.sh
~/bin/light_web_run.sh --ip 127.0.0.1 --port 9090 --directory /path/to/docs

# Or directly via venv
.venv/bin/python3 ultra_simple.py --host 0.0.0.0 --port 8080 --directory .
```

## Architecture

- **Single file**: `ultra_simple.py` (~17KB, ~430 lines)
- **Main class**: `UltraSimpleHandler(BaseHTTPRequestHandler)` handles everything
- HTML is generated inline — no template engine, no separate HTML files

## Dependencies

- `markdown >= 3.5` — Markdown→HTML conversion with extensions
- `Pygments >= 2.16` — Syntax highlighting (500+ languages)

## Key Methods

- `render_markdown()` — renders Markdown with syntax highlighting + LaTeX
- `render_source_code()` — highlights source files (`.py`, `.c`, `.cpp`, `.js`, etc.)
- `_math_on_html()` — replaces `$...$` / `$$...$$` with KaTeX-compatible markup
- `generate_directory_html()` — builds directory listing page
- `generate_markdown_html()` — wraps rendered Markdown in HTML page with KaTeX CDN
- `generate_source_html()` — wraps highlighted source code in HTML page

## Supported Source Files

`SOURCE_EXTS` in the code defines recognized file extensions. Covers: Python, C/C++, JS/TS, Java/Kotlin, Rust, Go, Ruby, PHP, Swift, Elixir, Haskell, Lua, Shell, SQL, CSS, YAML, TOML, JSON, and many more.

Special filenames: `Makefile`, `Dockerfile`, `CMakeLists.txt`, `.gitignore`.

## Conventions

- Server binds to `0.0.0.0` by default (accessible on network)
- Hidden files (`.` prefix) excluded from directory listing
- Path traversal protection via `Path.relative_to()` check
- All user content HTML-escaped via `html.escape()`
- Binary files fall back to `read_bytes()` (served as-is)
- Markdown files with `UnicodeDecodeError` fall back to binary serving

## LaTeX Math

- Uses **KaTeX** from CDN (client-side rendering, no Python dep)
- `$...$` for inline math, `$$...$$` for display math
- Math inside code blocks (`` `...` `` / ` ``` `) is protected and not processed
