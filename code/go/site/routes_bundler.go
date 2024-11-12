package site

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	Name        string `json:"name"`
	Path        string `json:"path"`
	Authors     string `json:"author"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Key         string `json:"key"`
	Contents    string `json:"contents"`
}

type PluginManifest struct {
	Version string          `json:"version"`
	Plugins []PluginDetails `json:"plugins"`
}

func setupBundler(router chi.Router) error {

	tmpDir, err := os.MkdirTemp("", "datastar-bundler")
	if err != nil {
		return fmt.Errorf("error creating temp dir: %w", err)
	}

	distDir := filepath.Join(tmpDir, "dist")
	if err = os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("error creating out dir: %w", err)
	}

	// Copy the static files to the temp dir.
	if err := fs.WalkDir(staticFS, "static/librarySource", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relDir := strings.TrimPrefix(path, "static/librarySource")

		if d.IsDir() {
			if err := os.MkdirAll(filepath.Join(tmpDir, relDir), 0755); err != nil {
				return err
			}
		} else {
			src, err := staticFS.Open(path)
			if err != nil {
				return fmt.Errorf("error opening static file: %w", err)
			}
			defer src.Close()

			dst, err := os.Create(filepath.Join(tmpDir, relDir))
			if err != nil {
				return fmt.Errorf("error creating file: %w", err)
			}

			if _, err := io.Copy(dst, src); err != nil {
				return fmt.Errorf("error copying static file: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error copying static files: %w", err)
	}

	manifest := PluginManifest{
		Version: datastar.Version,
	}

	allIncludedBundler, err := staticFS.ReadFile("static/librarySource/bundles/datastar.ts")
	if err != nil {
		return fmt.Errorf("error reading all included bundler: %w", err)
	}

	allIncludedBundlerContents := string(allIncludedBundler)
	matches := datastarBundlerRegexp.FindAllStringSubmatch(allIncludedBundlerContents, -1)
	for _, match := range matches {
		name := match[1]
		path := match[2]

		if !strings.HasPrefix(path, "../plugins") {
			continue
		}

		tsRelpath := path[3:] + ".ts"
		tsSrcFilepath := filepath.Join("static", "librarySource", tsRelpath)
		b, err := staticFS.ReadFile(tsSrcFilepath)
		if err != nil {
			return fmt.Errorf("error reading plugin file: %w", err)
		}
		contents := string(b)

		key := toolbelt.Snake(strings.ReplaceAll(tsSrcFilepath, string(filepath.Separator), "_"))

		details := PluginDetails{
			Name:     name,
			Path:     path,
			Key:      key,
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
			case "Icon":
				details.Icon = value
			}
		}

		manifest.Plugins = append(manifest.Plugins, details)
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

			sse := datastar.NewSSE(w, r)

			revisedManifest := PluginManifest{
				Version: manifest.Version,
			}
			for _, plugin := range manifest.Plugins {
				if !store.IncludedPlugines[plugin.Key] {
					continue
				}

				revisedManifest.Plugins = append(revisedManifest.Plugins, plugin)
			}

			bundleContents, err := bundlePlugins(tmpDir, revisedManifest)
			if err != nil {
				http.Error(w, "error bundling plugins: "+err.Error(), http.StatusInternalServerError)
				return
			}

			c := bundlerResultsFragment(*bundleContents)
			sse.RenderFragmentTempl(c)
		})
	})

	return nil
}

type BundleResults struct {
	Hash              string        `json:"hash"`
	SourceSize        uint64        `json:"sourceSize"`
	SourceSizeGzipped uint64        `json:"sourceSizeGzipped"`
	CompileTime       time.Duration `json:"compileTime"`
}

func bundlePlugins(tmpDir string, manifest PluginManifest) (results *BundleResults, err error) {
	start := time.Now()
	h := xxh3.New()
	h.WriteString(manifest.Version)
	for _, plugin := range manifest.Plugins {
		h.WriteString(plugin.Contents)
	}
	hash := h.Sum64()
	hashedName := fmt.Sprintf("datastar-%x", hash)
	bundleRelDir := filepath.Join("bundles", hashedName+".ts")

	distDir := filepath.Join(tmpDir, "dist")
	// distFile := filepath.Join(distDir, hashedName+".js")

	bundleOutFile := filepath.Join(tmpDir, bundleRelDir)
	bundleFileContents := bundlerContent(manifest)
	if err = os.WriteFile(bundleOutFile, []byte(bundleFileContents), 0644); err != nil {
		return nil, fmt.Errorf("error writing bundle file: %w", err)
	}

	buildResult := api.Build(api.BuildOptions{
		EntryPoints:       []string{bundleOutFile},
		Outdir:            distDir,
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelSilent,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Sourcemap:         api.SourceMapLinked,
		Target:            api.ES2023,
	})

	if len(buildResult.Errors) > 0 {
		errs := make([]error, len(buildResult.Errors))
		for i, err := range buildResult.Errors {
			errs[i] = errors.New(err.Text)
		}
		if err = errors.Join(errs...); err != nil {
			return nil, fmt.Errorf("error joining errors: %w", err)
		}
	}

	var contents []byte
	for _, file := range buildResult.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			contents = file.Contents
			break
		}
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	w := gzip.NewWriter(buf)
	if _, err := w.Write(contents); err != nil {
		return nil, fmt.Errorf("error writing gzipped contents: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}
	gzipContents := buf.Bytes()

	results = &BundleResults{
		Hash:              fmt.Sprintf("%x", hash),
		SourceSize:        uint64(len(contents)),
		SourceSizeGzipped: uint64(len(gzipContents)),
		CompileTime:       time.Since(start),
	}

	return results, nil

}
