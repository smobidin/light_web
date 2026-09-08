package main

import (
	"fmt"
	"strings"
)

// generatePDFHTML creates an HTML page for viewing PDF files
func generatePDFHTML(filepath string, themeManager *ThemeManager) string {
	// Extract just the filename for display
	filename := filepath[strings.LastIndex(filepath, "/")+1:]
	if filename == "" {
		filename = filepath
	}
	
	fn := escapeHTML(filename)
	viewerJS := `(function () {
  'use strict';
  var RAW_URL = location.pathname + '?raw=1';
  var viewer = document.getElementById('viewer');
  var pagesEl = document.getElementById('pages');
  var loadingEl = document.getElementById('loading');
  var zoomEl = document.getElementById('zoom-level');
  var pageInput = document.getElementById('page-num');
  var pageCountEl = document.getElementById('page-count');
  var prevBtn = document.getElementById('prev');
  var nextBtn = document.getElementById('next');
  var pageCount = 0;
  var pageNum = 1;
  var zoom = 1.0;
  var slots = [];
  var hasSlots = false;

  document.getElementById('open-new').href = RAW_URL;
  pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.11.174/pdf.worker.min.js';

  function setZoom(z) {
    zoom = z;
    zoomEl.textContent = Math.round(zoom * 100) + '%';
    for (var i = 0; i < slots.length; i++) {
      var s = slots[i];
      if (!s) continue;
      s.canvas.style.width = Math.floor(s.w * zoom) + 'px';
      s.canvas.style.height = Math.floor(s.h * zoom) + 'px';
      if (s.task) { s.task.cancel(); s.task = null; }
      s.rendered = false;
      s.scale = null;
    }
    renderVisible();
  }

  function renderSlot(idx) {
    var s = slots[idx];
    if (!s || (s.rendered && s.scale === zoom)) return;
    var scale = zoom;
    var dpr = window.devicePixelRatio || 1;
    var canvas = s.canvas;
    s.scale = scale;
    s.rendered = false;
    if (s.task) { s.task.cancel(); s.task = null; }
    canvas.width = Math.max(1, Math.floor(s.w * scale * dpr));
    canvas.height = Math.max(1, Math.floor(s.h * scale * dpr));
    canvas.style.width = Math.floor(s.w * scale) + 'px';
    canvas.style.height = Math.floor(s.h * scale) + 'px';
    var renderTask = s.page.render({
      canvasContext: canvas.getContext('2d'),
      viewport: s.page.getViewport({ scale: scale * dpr })
    });
    s.task = renderTask;
    renderTask.promise.then(function () {
      s.rendered = true;
    })['catch'](function () {});
  }

  function clearSlot(idx) {
    var s = slots[idx];
    if (!s) return;
    if (s.task) { s.task.cancel(); s.task = null; }
    s.canvas.width = 1;
    s.canvas.height = 1;
    s.rendered = false;
    s.scale = null;
  }

  function isNear(idx) {
    var s = slots[idx];
    if (!s) return false;
    var cr = s.canvas.getBoundingClientRect();
    var vr = viewer.getBoundingClientRect();
    return cr.bottom >= vr.top - 400 && cr.top <= vr.bottom + 400;
  }

  function renderVisible() {
    for (var i = 0; i < slots.length; i++) {
      if (!slots[i]) continue;
      if (isNear(i)) { renderSlot(i); } else { clearSlot(i); }
    }
  }

  function fitWidth() {
    if (!hasSlots) return;
    var maxW = 0;
    for (var i = 0; i < slots.length; i++) {
      if (slots[i] && slots[i].w > maxW) maxW = slots[i].w;
    }
    var w = viewer.clientWidth - 44;
    setZoom(Math.max(0.25, w / maxW));
  }

  function scrollToTop(idx) {
    var s = slots[idx];
    if (!s) return;
    var rel = s.canvas.getBoundingClientRect().top - viewer.getBoundingClientRect().top + viewer.scrollTop;
    viewer.scrollTop = rel;
  }

  function go(n) {
    n = Math.min(Math.max(1, n), pageCount);
    scrollToTop(n - 1);
  }

  function currentPage() {
    var vt = viewer.getBoundingClientRect().top;
    for (var i = 0; i < slots.length; i++) {
      var s = slots[i];
      if (!s) continue;
      if (s.canvas.getBoundingClientRect().bottom > vt + 10) return i + 1;
    }
    return pageNum;
  }

  function updateStatus() {
    pageNum = currentPage();
    pageInput.value = pageNum;
    pageCountEl.textContent = pageCount;
    prevBtn.disabled = pageNum <= 1;
    nextBtn.disabled = pageNum >= pageCount;
  }

  prevBtn.addEventListener('click', function () { go(pageNum - 1); });
  nextBtn.addEventListener('click', function () { go(pageNum + 1); });
  pageInput.addEventListener('change', function () {
    go(parseInt(pageInput.value, 10) || 1);
    updateStatus();
  });
  document.getElementById('zoomin').addEventListener('click', function () {
    setZoom(Math.min(6, zoom * 1.2));
  });
  document.getElementById('zoomout').addEventListener('click', function () {
    setZoom(Math.max(0.25, zoom / 1.2));
  });
  document.getElementById('fitw').addEventListener('click', fitWidth);

  viewer.addEventListener('scroll', updateStatus);

  window.addEventListener('resize', function () {
    if (hasSlots) fitWidth();
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'ArrowRight' || e.key === 'PageDown') { go(pageNum + 1); }
    else if (e.key === 'ArrowLeft' || e.key === 'PageUp') { go(pageNum - 1); }
  });

  var io = new IntersectionObserver(function (entries) {
    entries.forEach(function (en) {
      var idx = parseInt(en.target.getAttribute('data-index'), 10);
      if (en.isIntersecting) { renderSlot(idx); } else { clearSlot(idx); }
    });
  }, { root: viewer, rootMargin: '400px 0px 400px 0px' });

  pdfjsLib.getDocument(RAW_URL).promise.then(function (pdf) {
    pageCount = pdf.numPages;
    pageInput.max = pageCount;
    pageCountEl.textContent = pageCount;

    function loadOne(i) {
      return pdf.getPage(i).then(function (page) {
        var vp = page.getViewport({ scale: 1 });
        var canvas = document.createElement('canvas');
        canvas.className = 'pdf-slide';
        canvas.setAttribute('data-index', i - 1);
        canvas.style.width = Math.floor(vp.width * zoom) + 'px';
        canvas.style.height = Math.floor(vp.height * zoom) + 'px';
        pagesEl.appendChild(canvas);
        slots[i - 1] = { page: page, w: vp.width, h: vp.height, canvas: canvas, rendered: false, scale: null, task: null };
        io.observe(canvas);
      });
    }

    var chain = Promise.resolve();
    for (var i = 1; i <= pageCount; i++) chain = chain.then(loadOne.bind(null, i));
    return chain;
  }).then(function () {
    hasSlots = true;
    loadingEl.style.display = 'none';
    fitWidth();
    updateStatus();
  })['catch'](function (err) {
    loadingEl.style.display = 'none';
    pagesEl.innerHTML = '<p style="padding:3rem;color:#c0392b;">PDF load error: ' + err.message + '</p>';
  });
})();`;
	
	// Create navigation with theme toggle
	navHTML := themeManager.CreateNavWithThemeToggle(filepath)
	
	// Get toggle script without <script> tags
	toggleScript := `
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
				const defaultTheme = '` + string(themeManager.GetCurrentTheme()) + `';
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
	`
	
	// Combine JavaScript
	// combinedJS := viewerJS + ";\n" + toggleScript
	
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s — Light Web</title>
<script src="https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.11.174/pdf.min.js"></script>
<style>
:root {
%s
}

* { margin:0; padding:0; box-sizing:border-box; }
html,body { height:100%%; }
body { font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:var(--bg-color); color:var(--text-color); line-height:1.5; display:flex; flex-direction:column; }
.pdf-toolbar { flex:0 0 auto; background:var(--header-bg); border-bottom:1px solid var(--header-border);
               padding:0.6rem 1rem; display:flex; align-items:center; gap:0.75rem; flex-wrap:wrap; }
.fn { color:var(--text-color); font-size:0.85rem; margin-right:auto; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; max-width:45%%; opacity:0.7; }
.btn { display:inline-block; border:1px solid var(--header-border); background:var(--header-bg); color:var(--text-color); padding:0.25rem 0.7rem;
       border-radius:5px; cursor:pointer; text-decoration:none; font-size:0.85rem; line-height:1.5; }
.btn:hover { background:var(--table-row-hover); }
.theme-toggle { background:var(--header-bg); color:var(--text-color); border:1px solid var(--header-border); 
                border-radius:4px; padding:0.25rem 0.5rem; cursor:pointer; font-size:0.9rem; }
.theme-toggle:hover { background:var(--table-row-hover); }
.toolbar-group { display:flex; align-items:center; gap:0.35rem; }
.page-num { width:3.6rem; padding:0.2rem 0.4rem; border:1px solid var(--header-border); border-radius:5px; text-align:center; font-size:0.85rem; }
#zoom-level { min-width:3.2rem; text-align:center; font-size:0.85rem; color:var(--text-color); opacity:0.7; }
#page-count { font-size:0.85rem; color:var(--text-color); opacity:0.7; }
.pdf-viewer { flex:1 1 auto; overflow:auto; background:var(--table-row-hover); }
#pages { padding:16px 20px; }
.pdf-slide { display:block; margin:0 auto 14px; background:var(--nav-bg); box-shadow:0 1px 8px rgba(0,0,0,0.25); }
.pdf-slide:last-child { margin-bottom:24px; }
.pdf-loading { padding:3rem; text-align:center; color:var(--text-color); opacity:0.7; }
@media (max-width:640px) { .fn { display:none; } }
</style>
</head>
<body>
<div class="pdf-toolbar">
%s
<span class="toolbar-group">
<button id="prev" class="btn" title="Previous page" disabled>◀</button>
<input id="page-num" class="page-num" type="number" min="1" value="1">
<span>/</span>
<span id="page-count">–</span>
<button id="next" class="btn" title="Next page" disabled>▶</button>
</span>
<span class="toolbar-group">
<button id="zoomout" class="btn" title="Zoom out">−</button>
<span id="zoom-level">100%%</span>
<button id="zoomin" class="btn" title="Zoom in">+</button>
<button id="fitw" class="btn" title="Fit width">⤢</button>
</span>
<span class="toolbar-group">
<a id="open-new" class="btn" target="_blank" rel="noopener" href="#">Open in new tab ↗</a>
</span>
</div>
<div id="viewer" class="pdf-viewer">
<div id="loading" class="pdf-loading">Loading PDF…</div>
<div id="pages"></div>
</div>
<script>
%s
</script>
<script>
%s
</script>
</body>
</html>`, 
fn, themeManager.GetCSSVariables(), navHTML, viewerJS, toggleScript)
}