#!/usr/bin/env python3
"""
Light Web Server — Markdown server with syntax highlighting, LaTeX math,
directory browsing, and source code rendering.

Usage:
  python3 ultra_simple.py [--host ADDR] [--port PORT] [--directory DIR]
"""

import os
import sys
import re
import html
import mimetypes
import argparse
from pathlib import Path
from http.server import HTTPServer, BaseHTTPRequestHandler

try:
    import markdown
    from markdown.extensions import codehilite, fenced_code, tables, toc, sane_lists
    from pygments import highlight
    from pygments.lexers import get_lexer_by_name, guess_lexer, get_lexer_for_filename
    from pygments.formatters import HtmlFormatter
except ImportError as e:
    print(f"Missing dependencies: {e}")
    print("Install: pip install markdown Pygments")
    sys.exit(1)


SOURCE_EXTS = {
    '.py', '.pyw', '.pyx', '.pxd', '.pxi',
    '.c', '.h', '.cpp', '.hpp', '.cc', '.hh', '.cxx', '.hxx',
    '.js', '.ts', '.jsx', '.tsx', '.mjs', '.cjs',
    '.java', '.kt', '.kts', '.scala', '.clj', '.cljs', '.cljc',
    '.rs', '.go', '.rb', '.php', '.phtml', '.swift',
    '.pl', '.pm', '.t', '.r', '.m', '.mm',
    '.sh', '.bash', '.zsh', '.fish', '.lua', '.vim',
    '.sql', '.css', '.scss', '.less', '.sass',
    '.xml', '.yaml', '.yml', '.toml', '.json',
    '.proto', '.gradle', '.cmake', '.tex',
    '.hs', '.erl', '.hrl', '.ex', '.exs', '.elm',
    '.asm', '.s', '.S', '.inc', '.d', '.ml',
    '.zig', '.nim', '.crystal', '.racket', '.rkt',
}

SOURCE_FILENAMES = {
    'makefile', 'makefile.in', 'gnumakefile',
    'dockerfile', 'containerfile',
    'cmakelists.txt', '.env.example', '.gitignore',
}


class UltraSimpleHandler(BaseHTTPRequestHandler):
    def __init__(self, *args, directory='.', **kwargs):
        self.directory = Path(directory).resolve()
        super().__init__(*args, **kwargs)

    # ─── HTTP ────────────────────────────────────────────────────────

    def do_GET(self):
        try:
            query = self.path.split('?', 1)[1] if '?' in self.path else ''
            raw_pdf = 'raw=1' in query.split('&')
            if self.path == '/':
                fs_path = self.directory
            else:
                clean_path = self.path.split('?')[0].lstrip('/')
                fs_path = self.directory / clean_path
            fs_path = fs_path.resolve()
            try:
                fs_path.relative_to(self.directory)
            except ValueError:
                self.send_error(403, "Forbidden")
                return
            if fs_path.is_dir():
                self.serve_directory(fs_path)
            elif fs_path.is_file():
                self.serve_file(fs_path, raw_pdf=raw_pdf)
            else:
                self.send_error(404, "Not found")
        except Exception as e:
            self.send_error(500, f"Error: {e}")

    def serve_directory(self, dir_path):
        items = []
        if dir_path != self.directory:
            parent = dir_path.parent.relative_to(self.directory)
            items.append({'name': '..', 'path': str(parent), 'type': 'dir', 'icon': '📁'})
        for item in sorted(dir_path.iterdir()):
            if item.name.startswith('.'):
                continue
            rel = str(item.relative_to(self.directory))
            if item.is_dir():
                items.append({'name': item.name, 'path': rel, 'type': 'dir', 'icon': '📁'})
            else:
                icon = '📄'
                if item.suffix.lower() in ('.md', '.markdown'):
                    icon = '📝'
                elif item.suffix.lower() == '.pdf':
                    icon = '📕'
                elif item.suffix.lower() in SOURCE_EXTS:
                    icon = '💻'
                items.append({'name': item.name, 'path': rel, 'type': 'file', 'icon': icon, 'size': item.stat().st_size})
        self.send_html(self.generate_directory_html(dir_path, items))

    def serve_file(self, file_path, raw_pdf=False):
        ext = file_path.suffix.lower()
        name_lower = file_path.name.lower()

        if ext == '.pdf' and not raw_pdf:
            self.render_pdf(file_path)
            return

        if ext in ('.md', '.markdown'):
            try:
                content = file_path.read_text(encoding='utf-8')
                self.render_markdown(content, file_path)
                return
            except UnicodeDecodeError:
                pass

        if ext in SOURCE_EXTS or name_lower in SOURCE_FILENAMES:
            try:
                content = file_path.read_text(encoding='utf-8')
                self.render_source_code(content, file_path)
                return
            except (UnicodeDecodeError, Exception):
                pass

        content_type, _ = mimetypes.guess_type(str(file_path))
        if not content_type:
            content_type = 'application/octet-stream'
        data = file_path.read_bytes()
        self.send_response(200)
        self.send_header('Content-Type', content_type)
        self.send_header('Content-Length', str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    # ─── Rendering ──────────────────────────────────────────────────

    def render_markdown(self, content, file_path):
        extensions = [
            'codehilite', 'fenced_code', 'tables', 'toc', 'sane_lists',
            'attr_list', 'def_list', 'footnotes', 'admonition',
        ]
        ext_configs = {
            'codehilite': {
                'css_class': 'highlight',
                'use_pygments': True,
                'noclasses': True,
                'pygments_style': 'default',
                'linenos': True,
            },
            'toc': {
                'permalink': True,
            },
        }
        md = markdown.Markdown(extensions=extensions, extension_configs=ext_configs)
        html_body = md.convert(content)
        html_body = self._math_on_html(html_body)

        toc_data = getattr(md, 'toc', '')

        self.send_html(self.generate_markdown_html(file_path.name, html_body, toc_data))

    def render_source_code(self, content, file_path):
        try:
            lexer = get_lexer_for_filename(file_path)
        except Exception:
            try:
                lexer = guess_lexer(content)
            except Exception:
                lexer = get_lexer_by_name('text')
        formatter = HtmlFormatter(
            style='default',
            cssclass='highlight',
            noclasses=True,
            linenos=True,
        )
        body = highlight(content, lexer, formatter).strip()
        self.send_html(self.generate_source_html(file_path.name, body))

    def render_pdf(self, file_path):
        self.send_html(self.generate_pdf_html(file_path.name))

    def _math_on_html(self, html_content):
        codes = []

        def save(m):
            codes.append(m.group(0))
            return f'\x01MP{codes.__len__()-1}\x01'

        html_content = re.sub(r'<pre>.*?</pre>', save, html_content, flags=re.DOTALL)
        html_content = re.sub(r'<code>[^<]*</code>', save, html_content)
        html_content = re.sub(
            r'\$\$(.+?)\$\$',
            r'<div class="math-display">\\[\1\\]</div>',
            html_content,
            flags=re.DOTALL,
        )
        html_content = re.sub(
            r'(?<!\$)\$(.+?)\$(?!\$)',
            r'<span class="math-inline">\\(\1\\)</span>',
            html_content,
        )
        for i, c in enumerate(codes):
            html_content = html_content.replace(f'\x01MP{i}\x01', c)
        return html_content

    # ─── HTML pages ─────────────────────────────────────────────────

    def generate_markdown_html(self, filename, body, toc=''):
        fn = html.escape(filename)
        katex_css = (
            '<link rel="stylesheet" '
            'href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">'
        )
        katex_js = (
            '<script src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js">'
            '</script>'
            '<script src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/'
            'auto-render.min.js"></script>'
        )
        katex_init = (
            '<script>'
            'document.addEventListener("DOMContentLoaded",function(){'
            'renderMathInElement(document.body,{'
            'delimiters:['
            '{left:"\\\\[",right:"\\\\]",display:true},'
            '{left:"\\\\( ",right:" \\\\)",display:false}'
            ']'
            '})});'
            '</script>'
        )
        return f'''<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{fn} — Light Web</title>
{katex_css}
{katex_js}
<style>
* {{ margin:0; padding:0; box-sizing:border-box; }}
body {{ font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#fff;color:#212529;line-height:1.7; }}
.header {{ background:#f8f9fa;border-bottom:1px solid #dee2e6;padding:1rem 2rem; }}
.nav {{ max-width:1000px;margin:0 auto;display:flex;align-items:center;gap:1rem; }}
.nav a {{ color:#007bff;text-decoration:none;font-size:0.9rem; }}
.nav a:hover {{ text-decoration:underline; }}
.nav .fn {{ color:#6c757d;font-size:0.85rem;margin-left:auto; }}
.container {{ max-width:1000px;margin:0 auto;padding:1.5rem 2rem; }}
h1 {{ font-size:2.2rem;margin:0 0 1rem;border-bottom:2px solid #dee2e6;padding-bottom:0.5rem; }}
h2 {{ font-size:1.6rem;margin:2rem 0 0.75rem;border-bottom:1px solid #dee2e6;padding-bottom:0.3rem; }}
h3 {{ font-size:1.35rem;margin:1.5rem 0 0.5rem; }}
h4 {{ font-size:1.15rem;margin:1.2rem 0 0.5rem; }}
p {{ margin-bottom:1rem; }}
a {{ color:#007bff;text-decoration:none; }}
a:hover {{ text-decoration:underline; }}
blockquote {{ margin:1.5rem 0;padding:0.75rem 1.5rem;background:#f8f9fa;
            border-left:4px solid #007bff;color:#495057; }}
ul,ol {{ margin:0.75rem 0;padding-left:2rem; }}
li {{ margin-bottom:0.3rem; }}
table {{ width:100%;border-collapse:collapse;margin:1.5rem 0;border-radius:6px;overflow:hidden;box-shadow:0 1px 4px rgba(0,0,0,0.08); }}
th,td {{ padding:0.6rem 0.8rem;text-align:left;border-bottom:1px solid #dee2e6; }}
th {{ background:#f8f9fa;font-weight:600; }}
tr:hover {{ background:#f8f9fa; }}
img {{ max-width:100%;border-radius:6px;margin:1.5rem 0; }}
hr {{ border:none;height:1px;background:#dee2e6;margin:2rem 0; }}
pre {{ margin:1rem 0;border-radius:6px;overflow-x:auto; }}
pre code {{ background:transparent;padding:0;font-size:0.9rem; }}
code {{ font-family:'SF Mono',Consolas,'Fira Code',monospace;
       background:#f1f3f4;padding:0.15em 0.3em;border-radius:3px;font-size:0.85em; }}
pre {{ background:#f8f9fa;border:1px solid #e9ecef;padding:1.2rem; }}
.math-display {{ overflow-x:auto;margin:1rem 0;text-align:center; }}
.math-inline {{ }}
.highlight {{ background:transparent !important; }}
.highlight pre {{ background:transparent !important;border:none !important;padding:0 !important;margin:0 !important; }}
@media (max-width:768px) {{ .container {{ padding:1rem; }} h1 {{ font-size:1.8rem; }} }}
</style>
</head>
<body>
<div class="header">
<div class="nav">
<a href="/">← Back</a>
<span class="fn">{fn}</span>
</div>
</div>
<main class="container">
{body}
</main>
{katex_init}
</body>
</html>'''

    def generate_source_html(self, filename, body):
        fn = html.escape(filename)
        return f'''<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{fn} — Light Web</title>
<style>
* {{ margin:0; padding:0; box-sizing:border-box; }}
body {{ font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#fff;color:#212529;line-height:1.5; }}
.header {{ background:#f8f9fa;border-bottom:1px solid #dee2e6;padding:0.8rem 2rem; }}
.nav {{ max-width:1000px;margin:0 auto;display:flex;align-items:center;gap:1rem; }}
.nav a {{ color:#007bff;text-decoration:none;font-size:0.9rem; }}
.nav a:hover {{ text-decoration:underline; }}
.nav .fn {{ color:#6c757d;font-size:0.85rem;margin-left:auto; }}
.container {{ max-width:1000px;margin:0 auto;padding:1.5rem 2rem; }}
pre {{ margin:0;border-radius:6px;overflow-x:auto; }}
.highlight {{ background:#f8f9fa !important;border:1px solid #e9ecef;border-radius:6px;padding:1.2rem; }}
.highlight table.highlighttable {{ width:100%;border-collapse:collapse; }}
.highlight td {{ vertical-align:top;padding:0; }}
.highlight td.linenos {{ padding-right:1rem;white-space:nowrap; }}
.highlight td.code pre {{ line-height:inherit !important; }}
.highlight td.linenos pre {{ line-height:inherit !important; }}
@media (max-width:768px) {{ .container {{ padding:1rem; }} }}
</style>
</head>
<body>
<div class="header">
<div class="nav">
<a href="/">← Back</a>
<span class="fn">{fn}</span>
</div>
</div>
<main class="container">
{body}
</main>
</body>
</html>'''

    def generate_pdf_html(self, filename):
        fn = html.escape(filename)
        viewer_js = """(function () {
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
})();"""
        return f'''<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{fn} — Light Web</title>
<script src="https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.11.174/pdf.min.js"></script>
<style>
* {{ margin:0; padding:0; box-sizing:border-box; }}
html,body {{ height:100%; }}
body {{ font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#fff;color:#212529;line-height:1.5;display:flex;flex-direction:column; }}
.pdf-toolbar {{ flex:0 0 auto;background:#f8f9fa;border-bottom:1px solid #dee2e6;
               padding:0.6rem 1rem;display:flex;align-items:center;gap:0.75rem;flex-wrap:wrap; }}
.fn {{ color:#6c757d;font-size:0.85rem;margin-right:auto;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;max-width:45%; }}
.btn {{ display:inline-block;border:1px solid #dee2e6;background:#fff;color:#212529;padding:0.25rem 0.7rem;
       border-radius:5px;cursor:pointer;text-decoration:none;font-size:0.85rem;line-height:1.5; }}
.btn:hover {{ background:#e9ecef; }}
.btn:disabled {{ opacity:0.4;cursor:not-allowed; }}
.toolbar-group {{ display:flex;align-items:center;gap:0.35rem; }}
.page-num {{ width:3.6rem;padding:0.2rem 0.4rem;border:1px solid #dee2e6;border-radius:5px;text-align:center;font-size:0.85rem; }}
#zoom-level {{ min-width:3.2rem;text-align:center;font-size:0.85rem;color:#495057; }}
#page-count {{ font-size:0.85rem;color:#495057; }}
.pdf-viewer {{ flex:1 1 auto;overflow:auto;background:#f1f3f4; }}
#pages {{ padding:16px 20px; }}
.pdf-slide {{ display:block;margin:0 auto 14px;background:#fff;box-shadow:0 1px 8px rgba(0,0,0,0.25); }}
.pdf-slide:last-child {{ margin-bottom:24px; }}
.pdf-loading {{ padding:3rem;text-align:center;color:#6c757d; }}
@media (max-width:640px) {{ .fn {{ display:none; }} }}
</style>
</head>
<body>
<div class="pdf-toolbar">
  <a class="btn" href="/">← Back</a>
  <span class="fn">{fn}</span>
  <span class="toolbar-group">
    <button id="prev" class="btn" title="Previous page" disabled>◀</button>
    <input id="page-num" class="page-num" type="number" min="1" value="1">
    <span>/</span>
    <span id="page-count">–</span>
    <button id="next" class="btn" title="Next page" disabled>▶</button>
  </span>
  <span class="toolbar-group">
    <button id="zoomout" class="btn" title="Zoom out">−</button>
    <span id="zoom-level">100%</span>
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
{viewer_js}
</script>
</body>
</html>'''

    def generate_directory_html(self, dir_path, items):
        relative_path = dir_path.relative_to(self.directory)
        breadcrumb = html.escape(str(relative_path) if relative_path != Path('.') else '/')

        rows = ''
        for item in items:
            ep = html.escape(item['path'])
            en = html.escape(item['name'])
            if item['type'] == 'dir':
                rows += f'''<a href="/{ep}" class="fi"><span class="fi-icon">{item['icon']}</span><span class="fi-name">{en}</span></a>'''
            else:
                sz = f'{item["size"] / 1024:.1f}KB' if item.get('size') else ''
                rows += f'''<a href="/{ep}" class="fi"><span class="fi-icon">{item['icon']}</span><span class="fi-name">{en}</span><span class="fi-size">{sz}</span></a>'''

        return f'''<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{breadcrumb} — Light Web</title>
<style>
* {{ margin:0;padding:0;box-sizing:border-box; }}
body {{ font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#f8f9fa;color:#212529;line-height:1.6; }}
.header {{ background:#fff;padding:1rem 2rem;border-bottom:1px solid #dee2e6;box-shadow:0 1px 3px rgba(0,0,0,0.06); }}
.logo {{ font-size:1.4rem;font-weight:700;color:#007bff; }}
.wrap {{ max-width:1000px;margin:2rem auto;padding:0 2rem; }}
.browser {{ background:#fff;border-radius:8px;box-shadow:0 1px 6px rgba(0,0,0,0.08);overflow:hidden; }}
.bc {{ background:#f8f9fa;padding:0.75rem 1.5rem;border-bottom:1px solid #dee2e6;font-size:0.9rem;color:#6c757d; }}
.fi {{ display:flex;align-items:center;padding:0.6rem 1.5rem;border-bottom:1px solid #f0f0f0;
     text-decoration:none;color:#212529;transition:background 0.15s; }}
.fi:hover {{ background:#f8f9fa; }}
.fi-icon {{ font-size:1.1rem;margin-right:0.8rem;flex-shrink:0; }}
.fi-name {{ flex:1;font-weight:500; }}
.fi-size {{ font-size:0.85rem;color:#868e96;flex-shrink:0; }}
.stats {{ padding:0.6rem 1.5rem;font-size:0.85rem;color:#868e96;background:#f8f9fa;border-top:1px solid #dee2e6; }}
@media (max-width:768px) {{ .wrap {{ padding:0 1rem;margin:1rem auto; }} .fi {{ padding:0.5rem 1rem; }} }}
</style>
</head>
<body>
<div class="header"><div class="logo">Light Web</div></div>
<div class="wrap"><div class="browser"><div class="bc">📁 {breadcrumb}</div>
{rows}
<div class="stats">{len(items)} item{'s' if len(items)!=1 else ''}</div>
</div></div></body></html>'''

    def send_html(self, body):
        self.send_response(200)
        self.send_header('Content-Type', 'text/html; charset=utf-8')
        self.send_header('Cache-Control', 'no-cache')
        encoded = body.encode('utf-8')
        self.send_header('Content-Length', str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def log_message(self, format, *args):
        pass


def main():
    parser = argparse.ArgumentParser(description='Light Web Server')
    parser.add_argument('--host', default='', help='Bind address (default: all interfaces)')
    parser.add_argument('--port', type=int, default=8080, help='Port')
    parser.add_argument('--directory', default='.', help='Root directory')
    args = parser.parse_args()

    bind_host = args.host if args.host else '0.0.0.0'
    display_host = bind_host

    class HandlerFactory:
        def __init__(self, directory):
            self.directory = directory
        def __call__(self, *args, **kwargs):
            return UltraSimpleHandler(*args, directory=self.directory, **kwargs)

    os.chdir(args.directory)
    httpd = HTTPServer((args.host, args.port), HandlerFactory(args.directory))
    print(f'Light Web Server @ http://{display_host}:{args.port}', flush=True)
    print(f'Serving: {Path(args.directory).resolve()}', flush=True)
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print('\nStopped')


if __name__ == '__main__':
    main()
