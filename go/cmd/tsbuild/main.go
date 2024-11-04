package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	log.Print("Datastar built in TS compiler!")

	outDir := "ts/library/dist"

	os.RemoveAll(outDir)

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{
			"ts/library/bundles/datastar-core.ts",
			"ts/library/bundles/datastar-allinone.ts",
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
		// MangleQuoted:      api.MangleQuotedTrue,
		// External:          []string{"@starfederation/datastar"},
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	// if err := filepath.Walk(outDir, func(path string, info fs.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if info.IsDir() {
	// 		return nil
	// 	}

	// 	buf := bytebufferpool.Get()
	// 	defer bytebufferpool.Put(buf)

	// 	b, err := os.ReadFile(path)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	w := brotli.NewWriterV2(buf, brotli.DefaultCompression)
	// 	if _, err := w.Write(b); err != nil {
	// 		log.Print(err)
	// 		return err
	// 	}
	// 	w.Close()

	// 	ratio := float64(buf.Len()) / float64(len(b))
	// 	log.Printf(
	// 		"Brotli compressed %s from %s to %s (%.2f%%)",
	// 		path,
	// 		humanize.IBytes(uint64(len(b))),
	// 		humanize.IBytes(uint64(buf.Len())),
	// 		100-ratio*100,
	// 	)

	// 	return nil
	// }); err != nil {
	// 	log.Fatal(err)
	// }
}
