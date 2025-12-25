package markdown

import (
	"bytes"
	"regexp"
	"strings"
	"time"

	"site/internal/build/markdown/extensions"
	"site/internal/models"

	embed "github.com/13rac1/goldmark-embed"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	alertcallouts "github.com/zmtcreative/gm-alert-callouts"
	"gopkg.in/yaml.v3"
)

var frontmatterRegex = regexp.MustCompile(`(?s)^---\n(.+?)\n---\n(.*)$`)

// ParseFrontmatter parses YAML frontmatter from markdown content
func ParseFrontmatter(content []byte) (*models.PostFrontmatter, string, error) {
	matches := frontmatterRegex.FindSubmatch(content)
	if matches == nil {
		// No frontmatter, treat entire content as markdown
		return &models.PostFrontmatter{}, string(content), nil
	}

	var fm models.PostFrontmatter
	if err := yaml.Unmarshal(matches[1], &fm); err != nil {
		return nil, "", err
	}

	return &fm, string(matches[2]), nil
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// Renderer wraps goldmark with all custom extensions
type Renderer struct {
	md goldmark.Markdown
}

// NewRenderer creates a new markdown renderer with all extensions
func NewRenderer() *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			alertcallouts.NewAlertCallouts(
				alertcallouts.UseHybridIcons(),
				alertcallouts.WithFolding(true),
			),
			embed.New(),
			extensions.NewPDFEmbedExtension(),
			extensions.NewAsideExtension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // Allow raw HTML in markdown
		),
	)

	return &Renderer{md: md}
}

// Render converts markdown to HTML
func (r *Renderer) Render(source string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ExtractTOC extracts table of contents from markdown
func ExtractTOC(markdown string) []models.TOCItem {
	var items []models.TOCItem

	// Match markdown headings (## Heading)
	headingRegex := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	matches := headingRegex.FindAllStringSubmatch(markdown, -1)

	for _, match := range matches {
		level := len(match[1])
		text := strings.TrimSpace(match[2])
		id := slugify(text)

		items = append(items, models.TOCItem{
			Level: level,
			ID:    id,
			Text:  text,
		})
	}

	return items
}

// slugify converts text to a URL-friendly slug
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from ends
	s = strings.Trim(s, "-")

	return s
}

// StripHTML removes HTML tags from a string
func StripHTML(s string) string {
	// Remove all HTML tags
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	s = tagRegex.ReplaceAllString(s, " ")

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	s = spaceRegex.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}
