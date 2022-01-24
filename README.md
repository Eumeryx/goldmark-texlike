# goldmark-texlike

**goldmark-texlike** is an extension for the [goldmark](http://github.com/yuin/goldmark) to provide support for **TeX-like** grammar, inline math uses the `${` and `}$` delimiters by default, Block math uses the `$${` and `}$$` delimiters by default.

**Note: This software does not handle TeX code, you need to use MathJax or KaTeX to render the math.**

This software reference [goldmark-mathjax](http://github.com/litao91/goldmark-mathjax).

## Installation

```
go get github.com/Eumeryx/goldmark-texlike
```

## Usage

```go
import (
	"bytes"
	"fmt"

	texlike "github.com/Eumeryx/goldmark-texlike"
	"github.com/yuin/goldmark"
)

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(texlike.NewTexlike(texlike.DefaultConfig)),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(buf.String())
	}
}
```

## Config

```go
var DefaultConfig = Config{
	Block: TexBase{					// Block mathe configuration
		Delimit:       [2]string{`$${`, `}$$`},	// Set delimiters
		TagName:       `span`,			// HTML tag name
		ClassName:     []string{},		// HTML class name
		OutputDelimit: [2]string{},		// Replaces the delimiters in the output HTML
	},
	Inline: TexBase{				// Inline math configuration, see Block math
		Delimit:       [2]string{`${`, `}$`},
		TagName:       `span`,
		ClassName:     []string{},
		OutputDelimit: [2]string{},
	},
	RemoveOutputDelimit: false,			// Remove the delimiters from the output HTML
}
```

## Why is the delimiters `${`, `}$` and `$${`, `}$$`

See this [talk](https://talk.commonmark.org/t/mathjax-extension-for-latex-equations/698/17).

**You can customize math delimiters, but I can't vouch for the correctness of custom delimiters.**
 
## License
MIT
