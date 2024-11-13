package tsbuild

import (
	"compress/gzip"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/delaneyj/toolbelt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/goccy/go-json"
	"github.com/valyala/bytebufferpool"
)

func Build() error {
	version, err := extractVersion()
	if err != nil {
		return fmt.Errorf("error extracting version: %w", err)
	}

	if err := errors.Join(
		// createPluginManifest(),
		createBundles(),
		writeOutConsts(version),
	); err != nil {
		return fmt.Errorf("error creating bundles: %w", err)
	}

	return nil
}

func extractVersion() (string, error) {
	packageJSONPath := "code/ts/library/package.json"
	packageJSON, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return "", fmt.Errorf("error reading package.json: %w", err)
	}

	type PackageJSON struct {
		Version string `json:"version"`
	}
	pj := &PackageJSON{}
	if err := json.Unmarshal(packageJSON, pj); err != nil {
		return "", fmt.Errorf("error unmarshalling package.json: %w", err)
	}

	version := pj.Version

	// Write out the version to the version file.
	versionPath := "code/ts/library/src/engine/version.ts"
	versionContens := fmt.Sprintf("export const VERSION = '%s';\n", version)
	if err := os.WriteFile(versionPath, []byte(versionContens), 0644); err != nil {
		return "", fmt.Errorf("error writing version file: %w", err)
	}

	return version, nil
}

func createBundles() error {
	log.Print("Creating bundles...")
	defer log.Print("Bundles created!")

	outDir := "./bundles"
	os.RemoveAll(outDir)

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{
			"code/ts/library/src/bundles/datastar-core.ts",
			"code/ts/library/src/bundles/datastar.ts",
		},
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
		return errors.Join(errs...)
	}

	return nil
}

func writeOutConsts(version string) error {
	log.Print("Extracting version...")

	ConstsData.Version = version

	build, err := os.ReadFile("bundles/datastar.js")
	if err != nil {
		return fmt.Errorf("error reading datastar.js: %w", err)
	}
	ConstsData.VersionClientByteSize = len(build)

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	w, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(build); err != nil {
		panic(err)
	}
	w.Close()
	ConstsData.VersionClientByteSizeGzip = buf.Len()

	var zeroCased toolbelt.CasedString
	// Make sure all enum are stepu
	for _, enum := range ConstsData.Enums {
		for _, value := range enum.Values {
			if value.Name == zeroCased {
				value.Name = toolbelt.ToCasedString(value.Value)
			}
		}
		if enum.DefaultIndex >= 0 {
			enum.Default = enum.Values[enum.DefaultIndex]
		}
	}

	templates := map[string]func(data *ConstTemplateData) string{
		"code/go/sdk/consts.go": goConsts,
	}

	for path, tmplFn := range templates {
		contents := tmplFn(ConstsData)
		if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
			return fmt.Errorf("error writing version file: %w", err)
		}
	}

	return nil
}

// var templates = map[string]string{
// 	"code/go/sdk/consts.go": `
// package datastar

// import "time"

// const (
// 	Version                        = "{{version}}"
// 	VersionClientByteSize          = {{datastarSizeBytes}}
// 	VersionClientByteSizeGzip      = {{datastarGzipSizeBytes}}
// 	VersionClientByteSizeGzipHuman = "{{datastarGzipSizByteseHuman}}"

// 	DefaultSettleTime = {{defaultSettleTimeMs}} * time.Millisecond
// 	DefaultSseSendRetry = {{defaultSSESendRetryMs}} * time.Millisecond
// 	DefaultFragmentMergeMode = FragmentMergeMode("{{defaultFragmentMergeMode}}")
// )
// `,
// 	"code/php/sdk/src/Defaults.php": `
// <?php
// namespace starfederation\datastar;

// use starfederation\datastar\enums\FragmentMergeMode;

// class Defaults
// {
//     public const DEFAULT_SETTLE_DURATION = {{defaultSettleTimeMs}};
//     public const DEFAULT_SSE_SEND_RETRY = {{defaultSSESendRetryMs}};
//     public const DEFAULT_FRAGMENT_MERGE_MODE = FragmentMergeMode::Morph;
// }
// `,
// }

func durationToMs(d time.Duration) int {
	return int(d.Milliseconds())
}
