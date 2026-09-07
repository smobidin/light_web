# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- PDF viewer built on PDF.js (cdnjs): continuous page scrolling, lazy page
  rendering via `IntersectionObserver`, zoom (−/+), fit-width mode, page
  navigation (prev/next, page number, keyboard arrows), and "Open in new tab"
  link to the raw PDF.
  - `/dir/file.pdf` renders the viewer page.
  - `/dir/file.pdf?raw=1` serves the raw PDF (`application/pdf`).
  - PDF files are marked with a book icon (📕) in the directory listing.

### Fixed

- Source code rendering: removed redundant `<pre>` wrapper around Pygments
  output (double-wrapping produced stray empty lines).
- Source code rendering: applied `.strip()` to Pygments output to drop leading
  and trailing newlines.
- Source code rendering: aligned the Pygments linenos table (same line-height
  for code and line-number columns, `vertical-align: top`) so the code column
  no longer floats vertically with empty space at the top and bottom.