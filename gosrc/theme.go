package main

import (
	"strings"
)

// Theme represents the color theme
type Theme string

const (
	LightTheme Theme = "light"
	DarkTheme  Theme = "dark"
)

// ThemeManager manages the theme settings
type ThemeManager struct {
	currentTheme Theme
}

// NewThemeManager creates a new theme manager with default light theme
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		currentTheme: LightTheme,
	}
}

// SetTheme sets the current theme
func (tm *ThemeManager) SetTheme(theme Theme) {
	tm.currentTheme = theme
}

// GetCurrentTheme returns the current theme
func (tm *ThemeManager) GetCurrentTheme() Theme {
	return tm.currentTheme
}

// GetCSSVariables returns CSS variables for the current theme
func (tm *ThemeManager) GetCSSVariables() string {
	return `
		/* Light theme variables */
		--bg-color: #fff;
		--text-color: #212529;
		--header-bg: #f8f9fa;
		--header-border: #dee2e6;
		--link-color: #007bff;
		--link-hover: #0056b3;
		--code-bg: #f8f9fa;
		--code-border: #e9ecef;
		--table-header-bg: #f8f9fa;
		--table-row-hover: #f8f9fa;
		--table-border: #dee2e6;
		--blockquote-bg: #f8f9fa;
		--blockquote-border: #007bff;
		--blockquote-text: #495057;
		--nav-bg: #fff;
		--nav-text: #212529;
		--nav-link: #007bff;
		--highlight-bg: #f8f9fa;
		--highlight-border: #e9ecef;
		
		/* Dark theme variables */
		--dark-bg-color: #1e1e1e;
		--dark-text-color: #e0e0e0;
		--dark-header-bg: #2d2d2d;
		--dark-header-border: #444;
		--dark-link-color: #4da6ff;
		--dark-link-hover: #66b3ff;
		--dark-code-bg: #2d2d2d;
		--dark-code-border: #444;
		--dark-table-header-bg: #2d2d2d;
		--dark-table-row-hover: #333;
		--dark-table-border: #444;
		--dark-blockquote-bg: #2d2d2d;
		--dark-blockquote-border: #4da6ff;
		--dark-blockquote-text: #cccccc;
		--dark-nav-bg: #252526;
		--dark-nav-text: #e0e0e0;
		--dark-nav-link: #4da6ff;
		--dark-highlight-bg: #2d2d2d;
		--dark-highlight-border: #444;
`
}

// GetToggleScript returns JavaScript for theme toggling
func (tm *ThemeManager) GetToggleScript() string {
	return `
	<script>
		function toggleTheme() {
			const root = document.documentElement;
			const isDark = root.classList.contains('dark');
			
			if (isDark) {
				// Switch to light theme
				root.classList.remove('dark');
				// Update CSS variables for light theme
				root.style.setProperty('--bg-color', '#fff');
				root.style.setProperty('--text-color', '#212529');
				root.style.setProperty('--header-bg', '#f8f9fa');
				root.style.setProperty('--header-border', '#dee2e6');
				root.style.setProperty('--link-color', '#007bff');
				root.style.setProperty('--link-hover', '#0056b3');
				root.style.setProperty('--code-bg', '#f8f9fa');
				root.style.setProperty('--code-border', '#e9ecef');
				root.style.setProperty('--table-header-bg', '#f8f9fa');
				root.style.setProperty('--table-row-hover', '#f8f9fa');
				root.style.setProperty('--table-border', '#dee2e6');
				root.style.setProperty('--blockquote-bg', '#f8f9fa');
				root.style.setProperty('--blockquote-border', '#007bff');
				root.style.setProperty('--blockquote-text', '#495057');
				root.style.setProperty('--nav-bg', '#fff');
				root.style.setProperty('--nav-text', '#212529');
				root.style.setProperty('--nav-link', '#007bff');
				root.style.setProperty('--highlight-bg', '#f8f9fa');
				root.style.setProperty('--highlight-border', '#e9ecef');
				
				// Update theme toggle button text
				const toggleButton = document.getElementById('theme-toggle');
				if (toggleButton) {
					toggleButton.textContent = '🌙 Dark';
				}
				// Save preference to localStorage
				localStorage.setItem('theme', 'light');
			} else {
				// Switch to dark theme
				root.classList.add('dark');
				// Update CSS variables for dark theme
				root.style.setProperty('--bg-color', '#1e1e1e');
				root.style.setProperty('--text-color', '#e0e0e0');
				root.style.setProperty('--header-bg', '#2d2d2d');
				root.style.setProperty('--header-border', '#444');
				root.style.setProperty('--link-color', '#4da6ff');
				root.style.setProperty('--link-hover', '#66b3ff');
				root.style.setProperty('--code-bg', '#2d2d2d');
				root.style.setProperty('--code-border', '#444');
				root.style.setProperty('--table-header-bg', '#2d2d2d');
				root.style.setProperty('--table-row-hover', '#333');
				root.style.setProperty('--table-border', '#444');
				root.style.setProperty('--blockquote-bg', '#2d2d2d');
				root.style.setProperty('--blockquote-border', '#4da6ff');
				root.style.setProperty('--blockquote-text', '#cccccc');
				root.style.setProperty('--nav-bg', '#252526');
				root.style.setProperty('--nav-text', '#e0e0e0');
				root.style.setProperty('--nav-link', '#4da6ff');
				root.style.setProperty('--highlight-bg', '#2d2d2d');
				root.style.setProperty('--highlight-border', '#444');
				
				// Update theme toggle button text
				const toggleButton = document.getElementById('theme-toggle');
				if (toggleButton) {
					toggleButton.textContent = '☀️ Light';
				}
				// Save preference to localStorage
				localStorage.setItem('theme', 'dark');
			}
		}
		
		// Apply saved theme on page load
		document.addEventListener('DOMContentLoaded', function() {
			const savedTheme = localStorage.getItem('theme');
			const root = document.documentElement;
			
			if (savedTheme === 'dark') {
				root.classList.add('dark');
				// Update CSS variables for dark theme
				root.style.setProperty('--bg-color', '#1e1e1e');
				root.style.setProperty('--text-color', '#e0e0e0');
				root.style.setProperty('--header-bg', '#2d2d2d');
				root.style.setProperty('--header-border', '#444');
				root.style.setProperty('--link-color', '#4da6ff');
				root.style.setProperty('--link-hover', '#66b3ff');
				root.style.setProperty('--code-bg', '#2d2d2d');
				root.style.setProperty('--code-border', '#444');
				root.style.setProperty('--table-header-bg', '#2d2d2d');
				root.style.setProperty('--table-row-hover', '#333');
				root.style.setProperty('--table-border', '#444');
				root.style.setProperty('--blockquote-bg', '#2d2d2d');
				root.style.setProperty('--blockquote-border', '#4da6ff');
				root.style.setProperty('--blockquote-text', '#cccccc');
				root.style.setProperty('--nav-bg', '#252526');
				root.style.setProperty('--nav-text', '#e0e0e0');
				root.style.setProperty('--nav-link', '#4da6ff');
				root.style.setProperty('--highlight-bg', '#2d2d2d');
				root.style.setProperty('--highlight-border', '#444');
				
				// Update theme toggle button text
				const toggleButton = document.getElementById('theme-toggle');
				if (toggleButton) {
					toggleButton.textContent = '☀️ Light';
				}
			} else if (savedTheme === 'light') {
				root.classList.remove('dark');
				// Update CSS variables for light theme
				root.style.setProperty('--bg-color', '#fff');
				root.style.setProperty('--text-color', '#212529');
				root.style.setProperty('--header-bg', '#f8f9fa');
				root.style.setProperty('--header-border', '#dee2e6');
				root.style.setProperty('--link-color', '#007bff');
				root.style.setProperty('--link-hover', '#0056b3');
				root.style.setProperty('--code-bg', '#f8f9fa');
				root.style.setProperty('--code-border', '#e9ecef');
				root.style.setProperty('--table-header-bg', '#f8f9fa');
				root.style.setProperty('--table-row-hover', '#f8f9fa');
				root.style.setProperty('--table-border', '#dee2e6');
				root.style.setProperty('--blockquote-bg', '#f8f9fa');
				root.style.setProperty('--blockquote-border', '#007bff');
				root.style.setProperty('--blockquote-text', '#495057');
				root.style.setProperty('--nav-bg', '#fff');
				root.style.setProperty('--nav-text', '#212529');
				root.style.setProperty('--nav-link', '#007bff');
				root.style.setProperty('--highlight-bg', '#f8f9fa');
				root.style.setProperty('--highlight-border', '#e9ecef');
				
				// Update theme toggle button text
				const toggleButton = document.getElementById('theme-toggle');
				if (toggleButton) {
					toggleButton.textContent = '🌙 Dark';
				}
			} else {
				// Use default theme from server
				const defaultTheme = '` + string(tm.currentTheme) + `';
				if (defaultTheme === 'dark') {
					root.classList.add('dark');
					// Update CSS variables for dark theme
					root.style.setProperty('--bg-color', '#1e1e1e');
					root.style.setProperty('--text-color', '#e0e0e0');
					root.style.setProperty('--header-bg', '#2d2d2d');
					root.style.setProperty('--header-border', '#444');
					root.style.setProperty('--link-color', '#4da6ff');
					root.style.setProperty('--link-hover', '#66b3ff');
					root.style.setProperty('--code-bg', '#2d2d2d');
					root.style.setProperty('--code-border', '#444');
					root.style.setProperty('--table-header-bg', '#2d2d2d');
					root.style.setProperty('--table-row-hover', '#333');
					root.style.setProperty('--table-border', '#444');
					root.style.setProperty('--blockquote-bg', '#2d2d2d');
					root.style.setProperty('--blockquote-border', '#4da6ff');
					root.style.setProperty('--blockquote-text', '#cccccc');
					root.style.setProperty('--nav-bg', '#252526');
					root.style.setProperty('--nav-text', '#e0e0e0');
					root.style.setProperty('--nav-link', '#4da6ff');
					root.style.setProperty('--highlight-bg', '#2d2d2d');
					root.style.setProperty('--highlight-border', '#444');
					
					// Update theme toggle button text
					const toggleButton = document.getElementById('theme-toggle');
					if (toggleButton) {
						toggleButton.textContent = '☀️ Light';
					}
				} else {
					root.classList.remove('dark');
					// Update CSS variables for light theme
					root.style.setProperty('--bg-color', '#fff');
					root.style.setProperty('--text-color', '#212529');
					root.style.setProperty('--header-bg', '#f8f9fa');
					root.style.setProperty('--header-border', '#dee2e6');
					root.style.setProperty('--link-color', '#007bff');
					root.style.setProperty('--link-hover', '#0056b3');
					root.style.setProperty('--code-bg', '#f8f9fa');
					root.style.setProperty('--code-border', '#e9ecef');
					root.style.setProperty('--table-header-bg', '#f8f9fa');
					root.style.setProperty('--table-row-hover', '#f8f9fa');
					root.style.setProperty('--table-border', '#dee2e6');
					root.style.setProperty('--blockquote-bg', '#f8f9fa');
					root.style.setProperty('--blockquote-border', '#007bff');
					root.style.setProperty('--blockquote-text', '#495057');
					root.style.setProperty('--nav-bg', '#fff');
					root.style.setProperty('--nav-text', '#212529');
					root.style.setProperty('--nav-link', '#007bff');
					root.style.setProperty('--highlight-bg', '#f8f9fa');
					root.style.setProperty('--highlight-border', '#e9ecef');
					
					// Update theme toggle button text
					const toggleButton = document.getElementById('theme-toggle');
					if (toggleButton) {
						toggleButton.textContent = '🌙 Dark';
					}
				}
			}
		});
	</script>
	`
}

// CreateNavWithThemeToggle creates a navigation bar with a theme toggle button
// For directory pages, filepath should be empty
// For file pages, filepath should contain the full file path
func (tm *ThemeManager) CreateNavWithThemeToggle(filepath string) string {
	var navContent string
	if filepath == "" {
		// Directory page - show logo
		navContent = `<div class="logo">Light Web</div>`
	} else {
		// File page - show back link and filename
		// Extract just the filename for display
		filename := filepath[strings.LastIndex(filepath, "/")+1:]
		if filename == "" {
			filename = filepath
		}
		
		// Compute parent directory path
		parentPath := "/"
		if idx := strings.LastIndex(filepath, "/"); idx > 0 {
			parentPath = filepath[:idx]
		}
		// Если idx == 0 (файл в корне), parentPath остается "/"
		// Если idx == -1 (нет слэшей), parentPath остается "/"
		
		escapedFilename := escapeHTML(filename)
		escapedParentPath := escapeHTML(parentPath)
		navContent = `<a href="` + escapedParentPath + `">← Back</a><span class="fn">` + escapedFilename + `</span>`
	}
	
	// Add theme toggle button
	themeButtonText := "🌙 Dark"
	if tm.currentTheme == DarkTheme {
		themeButtonText = "☀️ Light"
	}
	return `<div class="nav">` + navContent + `<button id="theme-toggle" class="theme-toggle" onclick="toggleTheme()">` + themeButtonText + `</button></div>`
}