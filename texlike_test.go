package texlike

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/yuin/goldmark"
)

func TestNewTexlike(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(NewTexlike(DefaultConfig)),
	)

	mdContent := []byte(`
$${
e^{i\pi} + 1 = 0
}$$

Inline math ${ E=mc^2 }$
`)

	expected := `<p><span>$${
e^{i\pi} + 1 = 0
}$$</span></p>
<p>Inline math <span>${ E=mc^2 }$</span></p>
`

	var buf bytes.Buffer
	if err := md.Convert(mdContent, &buf); err == nil {
		if actual := buf.String(); actual != expected {
			t.Fatalf("\nA simple test:\nExpected:\n%#v\n\nActual:\n%#v\n", expected, actual)
		}
	} else {
		fmt.Println(err)
	}
}

func TestCreatHtmlTag(t *testing.T) {
	cases := []struct {
		name     string
		isBlock  bool
		rmDelim  bool
		texBase  TexBase
		expected [2]string
	}{
		{
			name:    "1. Inline",
			isBlock: false,
			rmDelim: false,
			texBase: TexBase{
				Delimit:       [2]string{"${", "}$"},
				TagName:       "span",
				ClassName:     []string{},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{"<span>${", "}$</span>"},
		},
		{
			name:    "2. Set OutputDelimit",
			isBlock: false,
			rmDelim: false,
			texBase: TexBase{
				Delimit:       [2]string{"${", "}$"},
				TagName:       "span",
				ClassName:     []string{},
				OutputDelimit: [2]string{"\\(", "\\)"},
			},
			expected: [2]string{`<span>\(`, `\)</span>`},
		},
		{
			name:    "3. Empty TagName",
			isBlock: false,
			rmDelim: false,
			texBase: TexBase{
				Delimit:       [2]string{"${", "}$"},
				TagName:       "",
				ClassName:     []string{},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{"${", "}$"},
		},
		{
			name:    "4. Set ClassName",
			isBlock: false,
			rmDelim: false,
			texBase: TexBase{
				Delimit:       [2]string{"${", "}$"},
				TagName:       "span",
				ClassName:     []string{"tex", "tex-inline"},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{`<span class="tex tex-inline">${`, `}$</span>`},
		},
		{
			name:    "5. Empty TagName && Set ClassName",
			isBlock: false,
			rmDelim: false,
			texBase: TexBase{
				Delimit:       [2]string{"${", "}$"},
				TagName:       "",
				ClassName:     []string{"tex", "tex-inline"},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{"${", "}$"},
		},
		{
			name:    "6. Inline && RemoveOutputDelimit",
			isBlock: false,
			rmDelim: true,
			texBase: TexBase{
				Delimit:       [2]string{"${", "}$"},
				TagName:       "span",
				ClassName:     []string{},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{"<span>", "</span>"},
		},
		{
			name:    "7. Block",
			isBlock: true,
			rmDelim: false,
			texBase: TexBase{
				Delimit:       [2]string{"$${", "}$$"},
				TagName:       "span",
				ClassName:     []string{},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{"<span>$${\n", "\n}$$</span>"},
		},
		{
			name:    "8. Block && RemoveOutputDelimit",
			isBlock: true,
			rmDelim: true,
			texBase: TexBase{
				Delimit:       [2]string{"$${", "}$$"},
				TagName:       "span",
				ClassName:     []string{},
				OutputDelimit: [2]string{},
			},
			expected: [2]string{"<span>", "</span>"},
		},
	}

	for _, c := range cases {
		actual := creatHtmlTag(&c.texBase, c.isBlock, c.rmDelim)
		if actual != c.expected {
			t.Fatalf("\n%v\nExpected:\n%#v\n\nActual:\n%#v\n", c.name, c.expected, actual)
		}
	}
}
