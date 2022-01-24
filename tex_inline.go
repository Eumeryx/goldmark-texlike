package texlike

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// AST
type TexInlineNode struct {
	ast.BaseInline
}

func (n *TexInlineNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindTexInline = ast.NewNodeKind("TexInline")

func (n *TexInlineNode) Kind() ast.NodeKind {
	return KindTexInline
}

func NewTexInlineNode() *TexInlineNode {
	return &TexInlineNode{
		BaseInline: ast.BaseInline{},
	}
}

// Parser
type TexInlineParser struct {
	prefix    []byte
	suffix    []byte
	prefixLen int
	suffixLen int
}

func NewTexInlineParser(delimit [2]string) parser.InlineParser {
	prefix := []byte(delimit[0])
	suffix := []byte(delimit[1])
	return &TexInlineParser{prefix, suffix, len(prefix), len(suffix)}
}

func (s *TexInlineParser) Trigger() []byte {
	return s.prefix[:1]
}

func (s *TexInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if !bytes.HasPrefix(line, s.prefix) {
		return nil
	}
	l, pos := block.Position()
	block.Advance(s.prefixLen)
	node := NewTexInlineNode()
	for {
		line, segment := block.PeekLine()
		if line == nil {
			block.SetPosition(l, pos)
			return nil
		}
		if i := bytes.Index(line, s.suffix); -1 != i {
			segment = segment.WithStop(segment.Start + i)
			if !segment.IsEmpty() {
				node.AppendChild(node, ast.NewRawTextSegment(segment))
			}
			block.Advance(i + s.suffixLen)
			return node
		}
		node.AppendChild(node, ast.NewRawTextSegment(segment))
		block.AdvanceLine()
	}
}

// Renderer
type TexInlineRenderer struct {
	openTag  string
	closeTag string
}

func NewTexInlineRenderer(htmlTag [2]string) renderer.NodeRenderer {
	return &TexInlineRenderer{htmlTag[0], htmlTag[1]}
}

func (r *TexInlineRenderer) renderTexInline(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(r.openTag)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := util.EscapeHTML(segment.Value(source))
			w.Write(value)
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString(r.closeTag)
	return ast.WalkContinue, nil
}

func (r *TexInlineRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTexInline, r.renderTexInline)
}
