package texlike

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
	"github.com/yuin/goldmark/util"
)

type texBlock struct {
	delimit [2]string
	htmlTag [2]string
}

func (e *texBlock) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewTexBlockParser(e.delimit), 501),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTexBlockRenderer(e.htmlTag), 502),
	))
}

func TestTexBlock(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(&texBlock{
			delimit: [2]string{`$${`, `}$$`},
			htmlTag: [2]string{"<span>$${\n", "\n}$$</span>"},
		}),
	)
	testutil.DoTestCaseFile(markdown, "_test/tex_block.txt", t, testutil.ParseCliCaseArg()...)
}
