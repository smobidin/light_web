package main

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
	if tm.currentTheme == DarkTheme {
		return `
		/* Dark theme variables */
		--bg-color: #1e1e1e;
		--text-color: #e0e0e0;
		--header-bg: #2d2d2d;
		--header-border: #444;
		--link-color: #4da6ff;
		--link-hover: #66b3ff;
		--code-bg: #2d2d2d;
		--code-border: #444;
		--table-header-bg: #2d2d2d;
		--table-row-hover: #333;
		--table-border: #444;
		--blockquote-bg: #2d2d2d;
		--blockquote-border: #4da6ff;
		--blockquote-text: #cccccc;
		--nav-bg: #252526;
		--nav-text: #e0e0e0;
		--nav-link: #4da6ff;
		--highlight-bg: #2d2d2d;
		--highlight-border: #444;
	`
	}
	
	// Light theme variables (default)
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
	`
}

// GetToggleScript returns JavaScript for theme toggling
func (tm *ThemeManager) GetToggleScript() string {
	return `
	<script>
		function toggleTheme() {
			const root = document.documentElement;
			const currentTheme = root.className || 'light';
			const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
			
			// Update CSS classes
			root.className = newTheme;
			
			// Update theme toggle button text
			const toggleButton = document.getElementById('theme-toggle');
			if (toggleButton) {
				toggleButton.textContent = newTheme === 'dark' ? '☀️ Light' : '🌙 Dark';
			}
			
			// Save preference to localStorage
			localStorage.setItem('theme', newTheme);
		}
		
		// Apply saved theme on page load
		document.addEventListener('DOMContentLoaded', function() {
			const savedTheme = localStorage.getItem('theme') || '` + string(tm.currentTheme) + `';
			const root = document.documentElement;
			root.className = savedTheme;
			
			// Update theme toggle button text
			const toggleButton = document.getElementById('theme-toggle');
			if (toggleButton) {
				toggleButton.textContent = savedTheme === 'dark' ? '☀️ Light' : '🌙 Dark';
			}
		});
	</script>
	`
}

// CreateNavWithThemeToggle creates a navigation bar with a theme toggle button
// For directory pages, filename should be empty
// For file pages, filename should contain the file name
func (tm *ThemeManager) CreateNavWithThemeToggle(filename string) string {
	var navContent string
	if filename == "" {
		// Directory page - show logo
		navContent = `<div class="logo">Light Web</div>`
	} else {
		// File page - show back link and filename
		escapedFilename := escapeHTML(filename)
		navContent = `<a href="/">← Back</a><span class="fn">` + escapedFilename + `</span>`
	}
	
	// Add theme toggle button
	themeButtonText := "🌙 Dark"
	return `<div class="nav">` + navContent + `<button id="theme-toggle" class="theme-toggle" onclick="toggleTheme()">` + themeButtonText + `</button></div>`
}