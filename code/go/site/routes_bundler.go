package site

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/delaneyj/toolbelt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/code/go/sdk"
	"github.com/valyala/bytebufferpool"
	"github.com/zeebo/xxh3"
)

var datastarBundlerRegexp = regexp.MustCompile(`import { (?P<name>[^"]*) } from "(?P<path>[^"]*)";`)

type BundlerStore struct {
	IncludedPlugines map[string]bool `json:"includedPlugins"`
}

type PluginDetails struct {
	Path        string `json:"path"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Authors     string `json:"author"`
	Description string `json:"description"`
	Contents    string `json:"contents"`
}

type PluginManifest struct {
	Version string           `json:"version"`
	Plugins []*PluginDetails `json:"plugins"`
}

type BundlerContentData struct {
	Keys  []string
	Paths []string
	Names []string
}

func setupBundler(router chi.Router) error {

	tmpDir, err := os.MkdirTemp("", "datastar-bundler")
	if err != nil {
		return fmt.Errorf("error creating temp dir: %w", err)
	}

	manifest := &PluginManifest{
		Version: datastar.Version,
	}

	allIncludedBundler, err := staticFS.ReadFile("static/librarySource/bundles/datastar.ts")
	if err != nil {
		return fmt.Errorf("error reading all included bundler: %w", err)
	}

	allBundleContent := &BundlerContentData{}
	allIncludedBundlerContents := string(allIncludedBundler)
	matches := datastarBundlerRegexp.FindAllStringSubmatch(allIncludedBundlerContents, -1)
	for _, match := range matches {
		allBundleContent.Names = append(allBundleContent.Names, match[1])
		allBundleContent.Paths = append(allBundleContent.Paths, match[2])

	}

	const embeddedBaseDir = "static/librarySource"
	if err := fs.WalkDir(staticFS, embeddedBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking library source: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		relpath := path[len(embeddedBaseDir)+1:]
		targetPath := filepath.Join(tmpDir, relpath)

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("error creating dir: %w", err)
		}

		f, err := staticFS.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer f.Close()

		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)

		if _, err := io.Copy(buf, f); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}

		if err := os.WriteFile(targetPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		// Add to manifest
		relPathParts := strings.Split(relpath, string(filepath.Separator))
		if relPathParts[0] != "plugins" {
			return nil // skip non-plugins
		}

		for _, part := range relPathParts {
			if part == "core" {
				return nil // skip core
			}
		}

		baseName := filepath.Base(relpath)
		ext := filepath.Ext(baseName)
		name := strings.TrimSuffix(baseName, ext)

		key := toolbelt.Snake(strings.ReplaceAll(relpath, string(filepath.Separator), "_"))
		contents := buf.String()

		details := &PluginDetails{
			Path:     relpath,
			Key:      key,
			Name:     name,
			Contents: contents,
		}

		lines := strings.Split(contents, "\n")
		for _, line := range lines {
			if !strings.HasPrefix(line, "//") {
				break
			}

			line = strings.TrimPrefix(line, "//")
			lineParts := strings.SplitN(line, ":", 2)

			if len(lineParts) != 2 {
				continue
			}

			key := strings.TrimSpace(lineParts[0])
			value := strings.TrimSpace(lineParts[1])

			switch key {
			case "Description":
				details.Description = value
			case "Authors":
				details.Authors = value
			}
		}

		manifest.Plugins = append(manifest.Plugins, details)

		return nil
	}); err != nil {
		return fmt.Errorf("error walking library source: %w", err)
	}

	router.Route("/bundler", func(bundlerRouter chi.Router) {
		bundlerRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			store := &BundlerStore{
				IncludedPlugines: map[string]bool{},
			}
			for _, plugin := range manifest.Plugins {
				store.IncludedPlugines[plugin.Key] = true
			}
			PageBundler(manifest, store).Render(r.Context(), w)
		})

		bundlerRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
			store := &BundlerStore{}
			if err := datastar.ParseIncoming(r, store); err != nil {
				http.Error(w, "error parsing request: "+err.Error(), http.StatusBadRequest)
				return
			}

			bundleContents, err := bundlePlugins(manifest, tmpDir, store)
			if err != nil {
				http.Error(w, "error bundling plugins: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(bundleContents))

		})
	})

	return nil
}

func bundlePlugins(manifest *PluginManifest, tmpDir string, store *BundlerStore) (bundledContents string, err error) {
	data := &BundlerContentData{}
	h := xxh3.New()
	// Create a new bundle based on the store
	for _, plugin := range manifest.Plugins {
		if !store.IncludedPlugines[plugin.Key] {
			continue
		}

		data.Names = append(data.Names, toolbelt.Pascal(plugin.Name)+"Plugin")
		data.Paths = append(data.Paths, plugin.Path)
		h.WriteString(plugin.Contents)
	}
	hash := h.Sum64()
	bundleName := fmt.Sprintf("datastar-%x.js", hash)
	bundleRelDir := filepath.Join("bundles", bundleName)
	log.Printf("Creating bundle %s", bundleRelDir)

	bundleOutFile := filepath.Join(tmpDir, bundleRelDir)
	bundleFileContents := bundlerContent(data)

	if err = os.WriteFile(bundleOutFile, []byte(bundleFileContents), 0644); err != nil {
		return "", fmt.Errorf("error writing bundle file: %w", err)
	}

	outDir := filepath.Join(tmpDir, "dist")
	if err = os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("error creating out dir: %w", err)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{bundleOutFile},
		Outdir:            outDir,
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Sourcemap:         api.SourceMapLinked,
		Target:            api.ES2023,
	})

	if len(result.Errors) > 0 {
		errs := make([]error, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = errors.New(err.Text)
		}
		if err = errors.Join(errs...); err != nil {
			return "", fmt.Errorf("error joining errors: %w", err)
		}
	}

	finalBundlePath := filepath.Join(outDir, bundleName)

	b, err := os.ReadFile(finalBundlePath)
	if err != nil {
		return "", fmt.Errorf("error reading bundle file: %w", err)
	}

	return string(b), nil
}
