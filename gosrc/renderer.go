package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	markdown "github.com/gomarkdown/markdown"
	mhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// renderMarkdown converts Markdown content to HTML with syntax highlighting
func renderMarkdown(content, filename string) (string, error) {
	// Check if content starts with Hugo front matter (+++ ... +++)
	frontMatterRegex := regexp.MustCompile(`(?s)^(\+{3}\n.*?\n\+{3}\n*)(.*)$`)
	matches := frontMatterRegex.FindStringSubmatch(content)
	
	var bodyContent string
	if len(matches) == 3 {
		// If we found front matter, process only the content part
		bodyContent = matches[2]
	} else {
		// Otherwise process the entire content
		bodyContent = content
	}
	
	// Set up the parser with common extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(bodyContent))

	// Set up the renderer
	htmlFlags := mhtml.CommonFlags | mhtml.HrefTargetBlank
	opts := mhtml.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := mhtml.NewRenderer(opts)

	// Render to HTML
	body := string(markdown.Render(doc, renderer))
	
	// Process math expressions
	body = processMath(body)

	// Generate complete HTML page
	return generateMarkdownHTML(filename, body, ""), nil
}

// renderSourceCode converts source code to HTML with syntax highlighting
func renderSourceCode(content, filename string) (string, error) {
	// Determine lexer based on filename
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(content)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	
	// Get the GitHub style
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}
	
	// Create HTML formatter with line numbers and classes
	formatter := html.New(html.WithLineNumbers(true), html.WithClasses(true))
	
	// Tokenize the content
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		// Fallback to plain text if tokenization fails
		return generateSourceHTML(filename, "<pre><code>"+escapeHTML(content)+"</code></pre>"), nil
	}
	
	// Format the tokens to HTML
	var buf strings.Builder
	if err := formatter.Format(&buf, style, iterator); err != nil {
		// Fallback to plain text if formatting fails
		return generateSourceHTML(filename, "<pre><code>"+escapeHTML(content)+"</code></pre>"), nil
	}
	
	body := buf.String()
	
	// Generate complete HTML page
	return generateSourceHTML(filename, body), nil
}

// generateMarkdownHTML creates a complete HTML page for Markdown content
func generateMarkdownHTML(filename, body, toc string) string {
	fn := escapeHTML(filename)
	
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s — Light Web</title>
<style>
* { margin:0; padding:0; box-sizing:border-box; }
body { font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#fff;color:#212529;line-height:1.7; }
.header { background:#f8f9fa;border-bottom:1px solid #dee2e6;padding:1rem 2rem; }
.nav { max-width:1000px;margin:0 auto;display:flex;align-items:center;gap:1rem; }
.nav a { color:#007bff;text-decoration:none;font-size:0.9rem; }
.nav a:hover { text-decoration:underline; }
.nav .fn { color:#6c757d;font-size:0.85rem;margin-left:auto; }
.container { max-width:1000px;margin:0 auto;padding:1.5rem 2rem; }
h1 { font-size:2.2rem;margin:0 0 1rem;border-bottom:2px solid #dee2e6;padding-bottom:0.5rem; }
h2 { font-size:1.6rem;margin:2rem 0 0.75rem;border-bottom:1px solid #dee2e6;padding-bottom:0.3rem; }
h3 { font-size:1.35rem;margin:1.5rem 0 0.5rem; }
h4 { font-size:1.15rem;margin:1.2rem 0 0.5rem; }
p { margin-bottom:1rem; }
a { color:#007bff;text-decoration:none; }
a:hover { text-decoration:underline; }
blockquote { margin:1.5rem 0;padding:0.75rem 1.5rem;background:#f8f9fa;
            border-left:4px solid #007bff;color:#495057; }
ul,ol { margin:0.75rem 0;padding-left:2rem; }
li { margin-bottom:0.3rem; }
table { width:100%%;border-collapse:collapse;margin:1.5rem 0;border-radius:6px;overflow:hidden;box-shadow:0 1px 4px rgba(0,0,0,0.08); }
th,td { padding:0.6rem 0.8rem;text-align:left;border-bottom:1px solid #dee2e6; }
th { background:#f8f9fa;font-weight:600; }
tr:hover { background:#f8f9fa; }
img { max-width:100%%;border-radius:6px;margin:1.5rem 0; }
hr { border:none;height:1px;background:#dee2e6;margin:2rem 0; }
pre { margin:1rem 0;border-radius:6px;overflow-x:auto; }
pre code { background:transparent;padding:0;font-size:0.9rem; }
code { font-family:'SF Mono',Consolas,'Fira Code',monospace;
       background:#f1f3f4;padding:0.15em 0.3em;border-radius:3px;font-size:0.85em; }
pre { background:#f8f9fa;border:1px solid #e9ecef;padding:1.2rem; }
.math-display { overflow-x:auto;margin:1rem 0;text-align:center; }
.math-inline { }
.highlight { background:transparent !important; }
.highlight pre { background:transparent !important;border:none !important;padding:0 !important;margin:0 !important; }
@media (max-width:768px) { .container { padding:1rem; } h1 { font-size:1.8rem; } }
</style>
</head>
<body>
<div class="header">
<div class="nav">
<a href="/">← Back</a>
<span class="fn">%s</span>
</div>
</div>
<main class="container">
%s
</main>
</body>
</html>`, fn, fn, body)
}

// generateSourceHTML creates a complete HTML page for source code
func generateSourceHTML(filename, body string) string {
	fn := escapeHTML(filename)
	
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s — Light Web</title>
<style>
* { margin:0; padding:0; box-sizing:border-box; }
body { font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;
       background:#fff;color:#212529;line-height:1.5; }
.header { background:#f8f9fa;border-bottom:1px solid #dee2e6;padding:0.8rem 2rem; }
.nav { max-width:1000px;margin:0 auto;display:flex;align-items:center;gap:1rem; }
.nav a { color:#007bff;text-decoration:none;font-size:0.9rem; }
.nav a:hover { text-decoration:underline; }
.nav .fn { color:#6c757d;font-size:0.85rem;margin-left:auto; }
.container { max-width:1000px;margin:0 auto;padding:1.5rem 2rem; }
pre { margin:0;border-radius:6px;overflow-x:auto; }
.highlight { background:#f8f9fa !important;border:1px solid #e9ecef;border-radius:6px;padding:1.2rem; }
.highlight table.highlighttable { width:100%%;border-collapse:collapse; }
.highlight td { vertical-align:top;padding:0; }
.highlight td.linenos { padding-right:1rem;white-space:nowrap; }
.highlight td.code pre { line-height:inherit !important; }
.highlight td.linenos pre { line-height:inherit !important; }
@media (max-width:768px) { .container { padding:1rem; } }
/* Chroma GitHub style */
/* Background */ .bg { color: #333; background-color: #fff; }
/* PreWrapper */ .chroma { color: #333; background-color: #fff; }
/* Error */ .chroma .err { color: #a61717; background-color: #e3d2d2 }
/* LineTableTD */ .chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
/* LineTable */ .chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; width: auto; overflow: auto; display: block; }
/* LineHighlight */ .chroma .hl { background-color: #ffffcc }
/* LineNumbersTable */ .chroma .lnt { margin-right: 0.4em; padding: 0 0.4em 0 0.4em; color: #7f7f7f }
/* LineNumbers */ .chroma .ln { margin-right: 0.4em; padding: 0 0.4em 0 0.4em; color: #7f7f7f }
/* Keyword */ .chroma .k { color: #000000; font-weight: bold }
/* KeywordConstant */ .chroma .kc { color: #000000; font-weight: bold }
/* KeywordDeclaration */ .chroma .kd { color: #000000; font-weight: bold }
/* KeywordNamespace */ .chroma .kn { color: #000000; font-weight: bold }
/* KeywordPseudo */ .chroma .kp { color: #000000; font-weight: bold }
/* KeywordReserved */ .chroma .kr { color: #000000; font-weight: bold }
/* KeywordType */ .chroma .kt { color: #445588; font-weight: bold }
/* NameAttribute */ .chroma .na { color: #008080 }
/* NameBuiltin */ .chroma .nb { color: #0086b3 }
/* NameBuiltinPseudo */ .chroma .bp { color: #999999 }
/* NameClass */ .chroma .nc { color: #445588; font-weight: bold }
/* NameConstant */ .chroma .no { color: #008080 }
/* NameDecorator */ .chroma .nd { color: #3c5d5d; font-weight: bold }
/* NameEntity */ .chroma .ni { color: #800080 }
/* NameException */ .chroma .ne { color: #990000; font-weight: bold }
/* NameFunction */ .chroma .nf { color: #990000; font-weight: bold }
/* NameLabel */ .chroma .nl { color: #990000; font-weight: bold }
/* NameNamespace */ .chroma .nn { color: #555555 }
/* NameTag */ .chroma .nt { color: #000080 }
/* NameVariable */ .chroma .nv { color: #008080 }
/* NameVariableClass */ .chroma .vc { color: #008080 }
/* NameVariableGlobal */ .chroma .vg { color: #008080 }
/* NameVariableInstance */ .chroma .vi { color: #008080 }
/* LiteralString */ .chroma .s { color: #dd1144 }
/* LiteralStringAffix */ .chroma .sa { color: #dd1144 }
/* LiteralStringBacktick */ .chroma .sb { color: #dd1144 }
/* LiteralStringChar */ .chroma .sc { color: #dd1144 }
/* LiteralStringDelimiter */ .chroma .dl { color: #dd1144 }
/* LiteralStringDoc */ .chroma .sd { color: #dd1144 }
/* LiteralStringDouble */ .chroma .s2 { color: #dd1144 }
/* LiteralStringEscape */ .chroma .se { color: #dd1144 }
/* LiteralStringHeredoc */ .chroma .sh { color: #dd1144 }
/* LiteralStringInterpol */ .chroma .si { color: #dd1144 }
/* LiteralStringOther */ .chroma .sx { color: #dd1144 }
/* LiteralStringRegexp */ .chroma .sr { color: #009926 }
/* LiteralStringSingle */ .chroma .s1 { color: #dd1144 }
/* LiteralStringSymbol */ .chroma .ss { color: #990073 }
/* LiteralNumber */ .chroma .m { color: #009999 }
/* LiteralNumberBin */ .chroma .mb { color: #009999 }
/* LiteralNumberFloat */ .chroma .mf { color: #009999 }
/* LiteralNumberHex */ .chroma .mh { color: #009999 }
/* LiteralNumberInteger */ .chroma .mi { color: #009999 }
/* LiteralNumberIntegerLong */ .chroma .il { color: #009999 }
/* LiteralNumberOct */ .chroma .mo { color: #009999 }
/* Operator */ .chroma .o { color: #000000; font-weight: bold }
/* OperatorWord */ .chroma .ow { color: #000000; font-weight: bold }
/* Comment */ .chroma .c { color: #999988; font-style: italic }
/* CommentHashbang */ .chroma .ch { color: #999988; font-style: italic }
/* CommentMultiline */ .chroma .cm { color: #999988; font-style: italic }
/* CommentSingle */ .chroma .c1 { color: #999988; font-style: italic }
/* CommentSpecial */ .chroma .cs { color: #999999; font-weight: bold; font-style: italic }
/* CommentPreproc */ .chroma .cp { color: #999999; font-weight: bold }
/* CommentPreprocFile */ .chroma .cpf { color: #999999; font-weight: bold }
/* GenericDeleted */ .chroma .gd { color: #000000; background-color: #ffdddd }
/* GenericEmph */ .chroma .ge { color: #000000; font-style: italic }
/* GenericError */ .chroma .gr { color: #aa0000 }
/* GenericHeading */ .chroma .gh { color: #999999 }
/* GenericInserted */ .chroma .gi { color: #000000; background-color: #ddffdd }
/* GenericOutput */ .chroma .go { color: #888888 }
/* GenericPrompt */ .chroma .gp { color: #555555 }
/* GenericStrong */ .chroma .gs { font-weight: bold }
/* GenericSubheading */ .chroma .gu { color: #aaaaaa }
/* GenericTraceback */ .chroma .gt { color: #aa0000 }
/* TextWhitespace */ .chroma .w { color: #bbbbbb }
</style>
</head>
<body>
<div class="header">
<div class="nav">
<a href="/">← Back</a>
<span class="fn">%s</span>
</div>
</div>
<main class="container">
%s
</main>
</body>
</html>`, fn, fn, body)
}