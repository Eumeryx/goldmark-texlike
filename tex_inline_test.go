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

type texInline struct {
	delimit [2]string
	htmlTag [2]string
}

func (e *texInline) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTexInlineParser(e.delimit), 501),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTexInlineRenderer(e.htmlTag), 502),
	))
}

func TestTexInline(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(&texInline{
			delimit: [2]string{`${`, `}$`},
			htmlTag: [2]string{`<span>${`, `}$</span>`},
		}),
	)
	testutil.DoTestCaseFile(markdown, "_test/tex_inline.txt", t, testutil.ParseCliCaseArg()...)
}
