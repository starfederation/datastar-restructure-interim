package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	start := time.Now()
	log.Print("Datastar built in TS compiler!")
	defer func() {
		log.Printf("Datastar built in %s", time.Since(start))
	}()

	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

}

func run(ctx context.Context) error {
	_ = ctx
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
		// MangleQuoted:      api.MangleQuotedTrue,
		// External:          []string{"@starfederation/datastar"},
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
