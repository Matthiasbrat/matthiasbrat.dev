package markdown

import (
	"bytes"
	"regexp"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Highlighter handles syntax highlighting for code blocks
type Highlighter struct {
	formatter *html.Formatter
	style     *chroma.Style
}

// NewHighlighter creates a new syntax highlighter
func NewHighlighter() *Highlighter {
	// Use a minimal, black & white friendly style
	formatter := html.New(
		html.WithClasses(true),
		html.WithLineNumbers(false),
		html.TabWidth(4),
	)

	// Use a neutral style that works with black & white
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	return &Highlighter{
		formatter: formatter,
		style:     style,
	}
}

// Highlight applies syntax highlighting to a code block
func (h *Highlighter) Highlight(code, language string) (string, error) {
	// Get the lexer for the language
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	// Format to HTML
	var buf bytes.Buffer
	if err := h.formatter.Format(&buf, h.style, iterator); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// CSS returns the CSS for the syntax highlighting style
func (h *Highlighter) CSS() string {
	var buf bytes.Buffer
	h.formatter.WriteCSS(&buf, h.style)
	return buf.String()
}

// ProcessCodeBlocks finds fenced code blocks and applies syntax highlighting
// This should be called AFTER directive processing but BEFORE markdown rendering
func ProcessCodeBlocks(content string, highlighter *Highlighter) string {
	// Match fenced code blocks: ```language\ncode\n```
	codeBlockRegex := regexp.MustCompile("(?ms)^```([a-zA-Z0-9_+-]*)\n(.*?)\n```$")

	return codeBlockRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := codeBlockRegex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		language := parts[1]
		code := parts[2]

		if language == "" {
			language = "text"
		}

		highlighted, err := highlighter.Highlight(code, language)
		if err != nil {
			// Fallback to plain code block
			return "<pre><code>" + escapeHTML(code) + "</code></pre>"
		}

		return "<div class=\"code-block\" data-language=\"" + language + "\">" + highlighted + "</div>"
	})
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	replacer := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&#39;",
	}
	for old, new := range replacer {
		s = regexp.MustCompile(regexp.QuoteMeta(old)).ReplaceAllString(s, new)
	}
	return s
}
