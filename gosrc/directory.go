package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// generateDirectoryHTML creates an HTML page for browsing a directory
func generateDirectoryHTML(dirPath, baseDir string, entries []fs.DirEntry) string {
	// Calculate the relative path for display
	relativePath := "/"
	if dirPath != baseDir {
		rel, err := filepath.Rel(baseDir, dirPath)
		if err == nil && !strings.HasPrefix(rel, "..") {
			relativePath = "/" + rel
		}
	}
	
	breadcrumb := escapeHTML(relativePath)
	
	// Generate the list of items
	var rows strings.Builder
	
	// Add parent directory link if not at the root
	if dirPath != baseDir {
		parent := filepath.Dir(dirPath)
		rel, err := getRelativePath(baseDir, parent)
		if err == nil {
			ep := escapeHTML("/" + rel)
			if ep == "//" {
				ep = "/"
			}
			rows.WriteString(fmt.Sprintf(`<a href="%s" class="fi"><span class="fi-icon">📁</span><span class="fi-name">..</span></a>`, ep))
		}
	}
	
	// Add directory entries
	itemCount := 0
	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		itemCount++
		name := entry.Name()
		relPath := filepath.Join(relativePath, name)
		if relativePath == "/" {
			relPath = "/" + name
		}
		
		ep := escapeHTML(relPath)
		en := escapeHTML(name)
		
		if entry.IsDir() {
			rows.WriteString(fmt.Sprintf(`<a href="%s" class="fi"><span class="fi-icon">📁</span><span class="fi-name">%s</span></a>`, ep, en))
		} else {
			icon := "📄"
			ext := strings.ToLower(filepath.Ext(name))
			
			if ext == ".md" || ext == ".markdown" {
				icon = "📝"
			} else if ext == ".pdf" {
				icon = "📕"
			} else if isSourceFile(ext, strings.ToLower(name)) {
				icon = "💻"
			}
			
			info, err := entry.Info()
			var sizeStr string
			if err == nil {
				sizeStr = fmt.Sprintf("%.1fKB", float64(info.Size())/1024.0)
			}
			
			rows.WriteString(fmt.Sprintf(`<a href="%s" class="fi"><span class="fi-icon">%s</span><span class="fi-name">%s</span><span class="fi-size">%s</span></a>`, ep, icon, en, sizeStr))
		}
	}
	
	// Generate the full HTML page
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s — Light Web</title>
<style>
* { margin:0;padding:0;box-sizing:border-box; }
body { font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#f8f9fa;color:#212529;line-height:1.6; }
.header { background:#fff;padding:1rem 2rem;border-bottom:1px solid #dee2e6;box-shadow:0 1px 3px rgba(0,0,0,0.06); }
.logo { font-size:1.4rem;font-weight:700;color:#007bff; }
.wrap { max-width:1000px;margin:2rem auto;padding:0 2rem; }
.browser { background:#fff;border-radius:8px;box-shadow:0 1px 6px rgba(0,0,0,0.08);overflow:hidden; }
.bc { background:#f8f9fa;padding:0.75rem 1.5rem;border-bottom:1px solid #dee2e6;font-size:0.9rem;color:#6c757d; }
.fi { display:flex;align-items:center;padding:0.6rem 1.5rem;border-bottom:1px solid #f0f0f0;
     text-decoration:none;color:#212529;transition:background 0.15s; }
.fi:hover { background:#f8f9fa; }
.fi-icon { font-size:1.1rem;margin-right:0.8rem;flex-shrink:0; }
.fi-name { flex:1;font-weight:500; }
.fi-size { font-size:0.85rem;color:#868e96;flex-shrink:0; }
.stats { padding:0.6rem 1.5rem;font-size:0.85rem;color:#868e96;background:#f8f9fa;border-top:1px solid #dee2e6; }
@media (max-width:768px) { .wrap { padding:0 1rem;margin:1rem auto; } .fi { padding:0.5rem 1rem; } }
</style>
</head>
<body>
<div class="header"><div class="logo">Light Web</div></div>
<div class="wrap"><div class="browser"><div class="bc">📁 %s</div>
%s
<div class="stats">%d item%s</div>
</div></div></body></html>`, 
breadcrumb, breadcrumb, rows.String(), itemCount, pluralize(itemCount))
}

// pluralize returns "s" if count is not 1
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}