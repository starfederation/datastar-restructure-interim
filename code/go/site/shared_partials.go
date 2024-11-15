package site

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	datastar "github.com/starfederation/datastar/code/go/sdk"
	"github.com/valyala/bytebufferpool"
)

var mdRenderer func() *mdhtml.Renderer

func markdownRenders(staticMdPath string) (mdElementRenderers map[string]string, mdAnchors map[string][]string, err error) {
	if mdRenderer == nil {
		htmlFormatter := html.New(html.WithClasses(true), html.TabWidth(2))
		if htmlFormatter == nil {
			return nil, nil, fmt.Errorf("couldn't create html formatter")
		}
		styleName := "nord"
		highlightStyle := styles.Get(styleName)
		if highlightStyle == nil {
			return nil, nil, fmt.Errorf("couldn't find style %s", styleName)
		}
		highlightCSSBuffer := &bytes.Buffer{}
		if err := htmlFormatter.WriteCSS(highlightCSSBuffer, highlightStyle); err != nil {
			return nil, nil, fmt.Errorf("error writing highlight css: %w", err)
		}
		highlightCSS = templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, fmt.Sprintf(`<style>%s</style>`, highlightCSSBuffer.String()))
			return err
		})
		// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
		htmlHighlight := func(w io.Writer, source, lang, defaultLang string) error {
			if lang == "" {
				lang = defaultLang
			}
			l := lexers.Get(lang)
			if l == nil {
				l = lexers.Analyse(source)
			}
			if l == nil {
				l = lexers.Fallback
			}
			l = chroma.Coalesce(l)

			it, err := l.Tokenise(nil, source)
			if err != nil {
				return err
			}
			return htmlFormatter.Format(w, highlightStyle, it)
		}

		mdRenderer = func() *mdhtml.Renderer {
			return mdhtml.NewRenderer(mdhtml.RendererOptions{
				Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank,
				RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
					skipDefaultRenderer := false
					switch n := node.(type) {
					case *ast.CodeBlock:
						defaultLang := ""
						lang := string(n.Info)
						htmlHighlight(w, string(n.Literal), lang, defaultLang)
						skipDefaultRenderer = true
					case *ast.Heading:
						if entering {
							break
						}
						buf := bytebufferpool.Get()
						defer bytebufferpool.Put(buf)
						level := strconv.Itoa(n.Level)
						if level != "1" {
                            buf.WriteString(`<a href="#`)
                            buf.WriteString(n.HeadingID)
                            buf.WriteString(`">#</a>`)
                        }
						buf.WriteString(`</h`)
						buf.WriteString(level)
						buf.WriteString(`>`)
						buf.WriteTo(w)
						skipDefaultRenderer = true
					}
					return ast.GoToNext, skipDefaultRenderer
				},
			})
		}
	}

	mdDir := "static/md/" + staticMdPath
	docs, err := staticFS.ReadDir(mdDir)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading docs dir: %w", err)
	}

	// regExpImg := regexp.MustCompile(`(?P<whole>!\[[^\]]+]\((?P<path>[^)]+)\))`)
	// prefix := []byte("/static/")

	mdElementRenderers = map[string]string{}
	mdAnchors = map[string][]string{}
	for _, de := range docs {
		fullPath := mdDir + "/" + de.Name()

		b, err := staticFS.ReadFile(fullPath)
		if err != nil {
			return nil, nil, fmt.Errorf("error reading doc %s: %w", de.Name(), err)
		}

		// Package version
		b = bytes.ReplaceAll(b, []byte("PACKAGE_VERSION"), []byte(datastar.Version))

		// Get all anchors
		anchors := []string{}
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				parts := strings.Split(line, " ")
				anchor := strings.Join(parts[1:], " ")
				anchors = append(anchors, anchor)
			}
		}

		mdParser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes)
		doc := mdParser.Parse(b)
		renderedHTML := string(markdown.Render(doc, mdRenderer()))

		name := de.Name()[0 : len(de.Name())-3]
		mdElementRenderers[name] = renderedHTML
		mdAnchors[name] = anchors
	}

	return mdElementRenderers, mdAnchors, nil
}

func KVPairsAttrs(kvPairs ...string) templ.Attributes {
	if len(kvPairs)%2 != 0 {
		panic("kvPairs must be a multiple of 2")
	}
	attrs := templ.Attributes{}
	for i := 0; i < len(kvPairs); i += 2 {
		attrs[kvPairs[i]] = kvPairs[i+1]
	}
	return attrs
}
