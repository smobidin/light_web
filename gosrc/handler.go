package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler serves files with various rendering options
type Handler struct {
	Directory    string
	ThemeManager *ThemeManager
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	rawPDF := query.Get("raw") == "1"

	// Clean the path
	path := filepath.Clean(r.URL.Path)
	if path == "/" {
		path = "."
	} else {
		path = strings.TrimPrefix(path, "/")
	}

	// Resolve the full path
	fullPath := filepath.Join(h.Directory, path)

	// Security check: ensure the path is within the serving directory
	relPath, err := filepath.Rel(h.Directory, fullPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Resolve symlinks and get the real path
	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if it's a directory
	info, err := os.Stat(resolvedPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		h.serveDirectory(w, r, resolvedPath)
		return
	}

	// Handle file serving
	h.serveFile(w, r, resolvedPath, rawPDF)
}

// serveDirectory generates and serves a directory listing
func (h *Handler) serveDirectory(w http.ResponseWriter, r *http.Request, dirPath string) {
	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	// Read directory entries
	entries, err := dir.ReadDir(-1)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate HTML
	htmlContent := generateDirectoryHTML(dirPath, h.Directory, entries, h.ThemeManager)
	
	// Send response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, htmlContent)
}

// serveFile serves a file with appropriate rendering
func (h *Handler) serveFile(w http.ResponseWriter, r *http.Request, filePath string, rawPDF bool) {
	ext := strings.ToLower(filepath.Ext(filePath))
	name := filepath.Base(filePath)
	nameLower := strings.ToLower(name)

	// Handle PDF files
	if (ext == ".pdf" || nameLower == "pdf") && !rawPDF {
		htmlContent := generatePDFHTML(name, h.ThemeManager)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, htmlContent)
		return
	}

	// Handle Markdown files
	if ext == ".md" || ext == ".markdown" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		htmlContent, err := renderMarkdown(string(content), name, h.ThemeManager)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, htmlContent)
		return
	}

	// Handle source code files
	if isSourceFile(ext, nameLower) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		htmlContent, err := renderSourceCode(string(content), name, h.ThemeManager)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, htmlContent)
		return
	}

	// Serve other files as-is
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}