package smoketests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/delaneyj/toolbelt"
	playwright "github.com/playwright-community/playwright-go"
	"github.com/starfederation/datastar/code/go/site"
)

const port = 8080

var expect = playwright.NewPlaywrightAssertions()

func TestSite(t *testing.T) {
	ctx := context.Background()
	start := time.Now()
	if err := playwright.Install(); err != nil {
		t.Fatal(fmt.Errorf("error installing playwright: %w", err))
	}
	log.Printf("playwright installed in %v", time.Since(start))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("running site on port %d", port)
	defer log.Print("closing site")

	readyCh := make(chan struct{})
	go site.RunBlocking(port, readyCh)(ctx)
	<-readyCh

	start = time.Now()
	log.Printf("running smoke tests")

	tests := map[string]pageFunc{
		"":                           todoExampleTest,
		"examples/active_search":     activeSearchExampleTest,
		"examples/animations":        animationsExampleTest,
		"examples/bulk_update":       bulkUpdateExampleTest,
		"examples/click_to_edit":     clickToEditExampleTest,
		"examples/click_to_load":     clickToLoadExampleTest,
		"examples/delete_row":        deleteRowExampleTest,
		"examples/dialogs_browser":   dialogsBrowserExampleTest,
		"examples/edit_row":          editRowExampleTest,
		"examples/fetch_indicator":   fetchIndicatorExampleTest,
		"examples/infinite_scroll":   infiniteScrollExampleTest,
		"examples/inline_validation": inlineValidationExampleTest,
		"examples/lazy_load":         lazyLoadExampleTest,
		"examples/lazy_tabs":         lazyTabsExampleTest,
		"examples/progress_bar":      progressbarExampleTest,
		"examples/raf_update":        rafUpdateExampleTest,
		"examples/store_changed":     storeChangedExampleTest,
		"examples/value_select":      valueSelectExampleTest,
	}

	if err := setupSmokeTests(ctx, tests); err != nil {
		t.Fatal(fmt.Errorf("error running smoke tests: %w", err))
	}
	log.Printf("entire smoke tests process took %v", time.Since(start))

}

type pageFunc func(ctx context.Context, page playwright.Page) error

func setupSmokeTests(ctx context.Context, fns map[string]pageFunc) error {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	eg := toolbelt.NewErrGroupSharedCtx(ctx)

	browsers := []playwright.BrowserType{
		pw.Chromium,
		// pw.Firefox,
		// pw.WebKit,
	}
	for _, browser := range browsers {
		eg.Go(func(ctx context.Context) error {
			browserName := browser.Name()
			log.Printf("running smoke tests in %s", browserName)

			browser, err := browser.Launch()
			if err != nil {
				log.Fatalf("could not launch browser: %v", err)
			}
			defer browser.Close()

			page, err := browser.NewPage()
			if err != nil {
				log.Fatalf("could not create page: %v", err)
			}

			errs := make([]error, 0, len(fns))

			testStart := time.Now()
			for subPage, fn := range fns {
				url := fmt.Sprintf("http://localhost:%d/%s", port, subPage)
				log.Printf("running smoke tests in %s on %s", browserName, url)
				if _, err = page.Goto(url); err != nil {
					return fmt.Errorf("could not go to %s: %w", url, err)
				}

				errs = append(errs, fn(ctx, page))
			}

			if err := errors.Join(errs...); err != nil {
				return fmt.Errorf("error running smoke tests in %s: %w", browserName, err)
			}

			log.Printf("smoke tests ran in %s in %v", browserName, time.Since(testStart))

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error running smoke tests: %w", err)
	}

	log.Print("ALL SMOKE TESTS PASSED")

	return nil
}
