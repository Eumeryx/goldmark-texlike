package texlike

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TexBase struct {
	Delimit       [2]string
	TagName       string
	ClassName     []string
	OutputDelimit [2]string
}

type Config struct {
	Block               TexBase
	Inline              TexBase
	RemoveOutputDelimit bool
}

var DefaultConfig = Config{
	Block: TexBase{
		Delimit:       [2]string{`$${`, `}$$`},
		TagName:       `span`,
		ClassName:     []string{},
		OutputDelimit: [2]string{},
	},
	Inline: TexBase{
		Delimit:       [2]string{`${`, `}$`},
		TagName:       `span`,
		ClassName:     []string{},
		OutputDelimit: [2]string{},
	},
	RemoveOutputDelimit: false,
}

func creatHtmlTag(b *TexBase, isBlock, rmDelim bool) [2]string {
	openTag, closeTag := "", ""
	if tagName := strings.TrimSpace(b.TagName); tagName != "" {
		classList := strings.TrimSpace(strings.Join(b.ClassName, ` `))
		if classList != "" {
			classList = ` class="` + classList + `"`
		}
		openTag = `<` + tagName + classList + `>`
		closeTag = `</` + tagName + `>`
	}
	if !rmDelim {
		Delimit := b.Delimit
		if b.OutputDelimit != [2]string{} {
			Delimit = b.OutputDelimit
		}
		openTag += Delimit[0]
		closeTag = Delimit[1] + closeTag
		if isBlock {
			openTag += "\n"
			closeTag = "\n" + closeTag
		}
	}
	return [2]string{openTag, closeTag}
}

type Texlike struct {
	blockDelimit  [2]string
	blockHtmlTag  [2]string
	inlineDelimit [2]string
	inlineHtmlTag [2]string
}

func NewTexlike(r Config) *Texlike {
	return &Texlike{
		blockDelimit:  r.Block.Delimit,
		blockHtmlTag:  creatHtmlTag(&r.Block, true, r.RemoveOutputDelimit),
		inlineDelimit: r.Inline.Delimit,
		inlineHtmlTag: creatHtmlTag(&r.Inline, false, r.RemoveOutputDelimit),
	}
}

func (e *Texlike) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewTexBlockParser(e.blockDelimit), 701),
	))
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTexInlineParser(e.inlineDelimit), 501),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTexBlockRenderer(e.blockHtmlTag), 501),
		util.Prioritized(NewTexInlineRenderer(e.inlineHtmlTag), 502),
	))
}
