package extensions

import (
	"fmt"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// PDFEmbed represents a PDF embed node in the AST
type PDFEmbed struct {
	ast.BaseBlock
	Src string
}

// Dump implements ast.Node.Dump
func (n *PDFEmbed) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Src": n.Src,
	}, nil)
}

// KindPDFEmbed is the kind for PDFEmbed nodes
var KindPDFEmbed = ast.NewNodeKind("PDFEmbed")

// Kind implements ast.Node.Kind
func (n *PDFEmbed) Kind() ast.NodeKind {
	return KindPDFEmbed
}

// NewPDFEmbed creates a new PDFEmbed node
func NewPDFEmbed(src string) *PDFEmbed {
	return &PDFEmbed{
		Src: src,
	}
}

// pdfEmbedParser parses PDF embed syntax
type pdfEmbedParser struct{}

// NewPDFEmbedParser creates a new PDF embed parser
func NewPDFEmbedParser() parser.BlockParser {
	return &pdfEmbedParser{}
}

// Trigger returns the characters that trigger this parser
func (p *pdfEmbedParser) Trigger() []byte {
	return []byte{':'}
}

// Open checks if the line starts a PDF embed block
func (p *pdfEmbedParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()

	// Match :::pdf{src="/path/to/file.pdf"}
	re := regexp.MustCompile(`^:::pdf\{src="([^"]+)"\}\s*$`)
	matches := re.FindSubmatch(line)

	if matches == nil {
		return nil, parser.NoChildren
	}

	src := string(matches[1])
	reader.Advance(len(line))

	return NewPDFEmbed(src), parser.NoChildren
}

// Continue is called when the parser should continue parsing
func (p *pdfEmbedParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

// Close is called when the parser is done
func (p *pdfEmbedParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// Nothing to do
}

// CanInterruptParagraph returns true if this parser can interrupt a paragraph
func (p *pdfEmbedParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine returns false
func (p *pdfEmbedParser) CanAcceptIndentedLine() bool {
	return false
}

// pdfEmbedHTMLRenderer renders PDFEmbed nodes to HTML
type pdfEmbedHTMLRenderer struct {
	html.Config
}

// NewPDFEmbedHTMLRenderer creates a new PDF embed HTML renderer
func NewPDFEmbedHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &pdfEmbedHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers rendering functions
func (r *pdfEmbedHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPDFEmbed, r.renderPDFEmbed)
}

// renderPDFEmbed renders a PDFEmbed node to HTML
func (r *pdfEmbedHTMLRenderer) renderPDFEmbed(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*PDFEmbed)

	html := fmt.Sprintf(`<div class="pdf-embed">
<iframe src="%s" width="100%%" height="600" type="application/pdf" title="PDF Document">
<p>Your browser does not support PDF embeds. <a href="%s">Download the PDF</a>.</p>
</iframe>
</div>`, n.Src, n.Src)

	w.WriteString(html)
	w.WriteString("\n")

	return ast.WalkContinue, nil
}

// PDFEmbedExtension is a goldmark extension for PDF embeds
type PDFEmbedExtension struct{}

// Extend extends the goldmark parser with PDF embed support
func (e *PDFEmbedExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewPDFEmbedParser(), 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewPDFEmbedHTMLRenderer(), 500),
		),
	)
}

// NewPDFEmbedExtension creates a new PDF embed extension
func NewPDFEmbedExtension() goldmark.Extender {
	return &PDFEmbedExtension{}
}
