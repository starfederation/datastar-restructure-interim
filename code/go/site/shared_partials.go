package site

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

func markdownRenders(staticMdPath string) (mdElementRenderers map[string]string, mdAnchors map[string][]string, err error) {
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
