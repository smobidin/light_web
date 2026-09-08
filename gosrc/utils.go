package main

import (
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"strings"
)

// SOURCE_EXTS defines file extensions that should be treated as source code
var SOURCE_EXTS = map[string]bool{
	".py":     true,
	".pyw":    true,
	".pyx":    true,
	".pxd":    true,
	".pxi":    true,
	".c":      true,
	".h":      true,
	".cpp":    true,
	".hpp":    true,
	".cc":     true,
	".hh":     true,
	".cxx":    true,
	".hxx":    true,
	".js":     true,
	".ts":     true,
	".jsx":    true,
	".tsx":    true,
	".mjs":    true,
	".cjs":    true,
	".java":   true,
	".kt":     true,
	".kts":    true,
	".scala":  true,
	".clj":    true,
	".cljs":   true,
	".cljc":   true,
	".rs":     true,
	".go":     true,
	".rb":     true,
	".php":    true,
	".phtml":  true,
	".swift":  true,
	".pl":     true,
	".pm":     true,
	".t":      true,
	".r":      true,
	".m":      true,
	".mm":     true,
	".sh":     true,
	".bash":   true,
	".zsh":    true,
	".fish":   true,
	".lua":    true,
	".vim":    true,
	".sql":    true,
	".css":    true,
	".scss":   true,
	".less":   true,
	".sass":   true,
	".xml":    true,
	".yaml":   true,
	".yml":    true,
	".toml":   true,
	".json":   true,
	".proto":  true,
	".gradle": true,
	".cmake":  true,
	".tex":    true,
	".hs":     true,
	".erl":    true,
	".hrl":    true,
	".ex":     true,
	".exs":    true,
	".elm":    true,
	".asm":    true,
	".s":      true,
	".S":      true,
	".inc":    true,
	".d":      true,
	".ml":     true,
	".zig":    true,
	".nim":    true,
	".crystal": true,
	".racket": true,
	".rkt":    true,
}

// SOURCE_FILENAMES defines special filenames that should be treated as source code
var SOURCE_FILENAMES = map[string]bool{
	"makefile":      true,
	"makefile.in":   true,
	"gnumakefile":   true,
	"dockerfile":    true,
	"containerfile": true,
	"cmakelists.txt": true,
	".env.example":  true,
	".gitignore":    true,
}

// isSourceFile checks if a file should be treated as source code
func isSourceFile(ext, name string) bool {
	if SOURCE_EXTS[ext] {
		return true
	}
	
	// Check for special filenames
	return SOURCE_FILENAMES[strings.ToLower(name)]
}

// escapeHTML escapes HTML characters in a string
func escapeHTML(s string) string {
	return html.EscapeString(s)
}

// getRelativePath gets the relative path from base to target
func getRelativePath(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	
	// Clean the path to remove any .. elements
	rel = filepath.Clean(rel)
	
	// Check if the path is outside the base directory
	if strings.HasPrefix(rel, "..") {
		return "", filepath.ErrBadPattern
	}
	
	return rel, nil
}

// processMath processes LaTeX math expressions in HTML content
func processMath(htmlContent string) string {
	// Store code blocks to protect them from math processing
	codes := []string{}
	
	// Save code blocks
	savePre := func(m string) string {
		codes = append(codes, m)
		return fmt.Sprintf("\x01MP%d\x01", len(codes)-1)
	}
	
	// Protect pre blocks
	rePre := regexp.MustCompile(`<pre>.*?</pre>`)
	htmlContent = rePre.ReplaceAllStringFunc(htmlContent, func(m string) string {
		return savePre(m)
	})
	
	// Protect code blocks
	reCode := regexp.MustCompile(`<code>[^<]*</code>`)
	htmlContent = reCode.ReplaceAllStringFunc(htmlContent, func(m string) string {
		codes = append(codes, m)
		return fmt.Sprintf("\x01MP%d\x01", len(codes)-1)
	})
	
	// Process display math ($$...$$)
	reDisplayMath := regexp.MustCompile(`\$\$(.+?)\$\$`)
	htmlContent = reDisplayMath.ReplaceAllString(htmlContent, `<div class="math-display">\\[$1\\]</div>`)
	
	// Process inline math ($...$)
	// Using a more compatible regex without negative lookbehind/lookahead
	reInlineMath := regexp.MustCompile(`([^\\]|^)\$(.+?)\$([^\\]|$)`)
	htmlContent = reInlineMath.ReplaceAllString(htmlContent, `$1<span class="math-inline">\\($2\\)</span>$3`)
	
	// Restore code blocks
	for i, code := range codes {
		placeholder := fmt.Sprintf("\x01MP%d\x01", i)
		htmlContent = strings.Replace(htmlContent, placeholder, code, -1)
	}
	
	return htmlContent
}