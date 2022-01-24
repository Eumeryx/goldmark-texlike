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
type TexBlockNode struct {
	ast.BaseBlock
}

var KindTexBlock = ast.NewNodeKind("TexBlock")

func NewTexBlockNode() *TexBlockNode {
	return &TexBlockNode{}
}

func (n *TexBlockNode) Dump(source []byte, level int) {
	m := map[string]string{}
	ast.DumpHelper(n, source, level, m, nil)
}

func (n *TexBlockNode) Kind() ast.NodeKind {
	return KindTexBlock
}

func (n *TexBlockNode) IsRaw() bool {
	return true
}

// Parser
type TexBlockParser struct {
	prefix    []byte
	suffix    []byte
	prefixLen int
	suffixLen int
}

func NewTexBlockParser(delimit [2]string) parser.BlockParser {
	prefix := []byte(delimit[0])
	suffix := []byte(delimit[1])
	return &TexBlockParser{prefix, suffix, len(prefix), len(suffix)}
}

func (s *TexBlockParser) Trigger() []byte {
	return s.prefix[:1]
}

func (b *TexBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	start := pc.BlockOffset()
	if start < 0 || !bytes.HasPrefix(line[start:], b.prefix) {
		return nil, parser.NoChildren
	}
	node := NewTexBlockNode()
	l, pos := reader.Position()
	if util.IsBlank(line[start+b.prefixLen:]) {
		reader.AdvanceLine()
	} else {
		reader.Advance(start + b.prefixLen)
	}
	for {
		line, segment := reader.PeekLine()
		if util.IsBlank(line) {
			reader.SetPosition(l, pos)
			return nil, parser.NoChildren
		}
		end := bytes.Index(line, b.suffix)
		if end != -1 && util.IsBlank(line[end+b.suffixLen:]) {
			segment = segment.WithStop(segment.Start + end)
			if !segment.IsEmpty() {
				node.Lines().Append(segment)
			}
			return node, parser.NoChildren
		} else {
			node.Lines().Append(segment)
		}
		reader.AdvanceLine()
	}
}

func (b *TexBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (b *TexBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
}

func (b *TexBlockParser) CanInterruptParagraph() bool {
	return true
}

func (b *TexBlockParser) CanAcceptIndentedLine() bool {
	return false
}

// Renderer
type TexBlockRenderer struct {
	openTag  string
	closeTag string
}

func NewTexBlockRenderer(htmlTag [2]string) renderer.NodeRenderer {
	return &TexBlockRenderer{htmlTag[0], htmlTag[1]}
}

func (r *TexBlockRenderer) renderTexBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*TexBlockNode)
	if entering {
		_, _ = w.WriteString(`<p>` + r.openTag)
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString(r.closeTag + "</p>\n")
	}
	return ast.WalkContinue, nil
}

func (r *TexBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTexBlock, r.renderTexBlock)
}

func (r *TexBlockRenderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		value := util.EscapeHTML(line.Value(source))
		if tail := len(value) - 1; i == l-1 && value[tail] == '\n' {
			w.Write(value[:tail])
		} else {
			w.Write(value)
		}
	}
}
