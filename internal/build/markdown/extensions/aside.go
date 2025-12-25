package extensions

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Aside represents an aside node in the AST
type Aside struct {
	ast.BaseBlock
}

// Dump implements ast.Node.Dump
func (n *Aside) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindAside is the kind for Aside nodes
var KindAside = ast.NewNodeKind("Aside")

// Kind implements ast.Node.Kind
func (n *Aside) Kind() ast.NodeKind {
	return KindAside
}

// NewAside creates a new Aside node
func NewAside() *Aside {
	return &Aside{}
}

// asideParser parses aside syntax
type asideParser struct{}

// NewAsideParser creates a new aside parser
func NewAsideParser() parser.BlockParser {
	return &asideParser{}
}

// Trigger returns the characters that trigger this parser
func (p *asideParser) Trigger() []byte {
	return []byte{':'}
}

// Open checks if the line starts an aside block
func (p *asideParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	lineStr := strings.TrimSpace(string(line))

	// Match "::: aside" or ":::aside" at start of line
	if lineStr == ":::aside" || lineStr == "::: aside" {
		reader.Advance(len(line))
		return NewAside(), parser.HasChildren
	}

	return nil, parser.NoChildren
}

// Continue is called when the parser should continue parsing
func (p *asideParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	lineStr := strings.TrimSpace(string(line))

	// Check for closing :::
	if lineStr == ":::" {
		reader.Advance(segment.Len())
		return parser.Close
	}

	node.Lines().Append(segment)
	return parser.Continue | parser.HasChildren
}

// Close is called when the parser is done
func (p *asideParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// Nothing to do
}

// CanInterruptParagraph returns true if this parser can interrupt a paragraph
func (p *asideParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine returns false
func (p *asideParser) CanAcceptIndentedLine() bool {
	return false
}

// asideHTMLRenderer renders Aside nodes to HTML
type asideHTMLRenderer struct {
	html.Config
}

// NewAsideHTMLRenderer creates a new aside HTML renderer
func NewAsideHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &asideHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers rendering functions
func (r *asideHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindAside, r.renderAside)
}

// renderAside renders an Aside node to HTML
func (r *asideHTMLRenderer) renderAside(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("<aside class=\"aside\">\n")
		w.WriteString("<div class=\"aside-content\">\n")
	} else {
		w.WriteString("</div>\n")
		w.WriteString("</aside>\n")
	}
	return ast.WalkContinue, nil
}

// AsideExtension is a goldmark extension for asides
type AsideExtension struct{}

// Extend extends the goldmark parser with aside support
func (e *AsideExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewAsideParser(), 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewAsideHTMLRenderer(), 500),
		),
	)
}

// NewAsideExtension creates a new aside extension
func NewAsideExtension() goldmark.Extender {
	return &AsideExtension{}
}
