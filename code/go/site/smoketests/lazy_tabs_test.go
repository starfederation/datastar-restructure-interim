package smoketests

import (
	"context"
	"errors"
	"fmt"

	playwright "github.com/playwright-community/playwright-go"
)

func lazyTabsExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		lazyTabsTest,
	}

	errs := []error{}

	for _, test := range suite {
		err := test(page)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func lazyTabsTest(page playwright.Page) error {
	selector := "#tab_content"

	initialText, err := page.Locator(selector).TextContent()
	if err != nil {
		return fmt.Errorf("err while attempting to read TextContent: %w", err)
	}
	page.GetByTestId("tab_1").Click()

	condition := fmt.Sprintf(`
		() => {
			const element = document.querySelector('%s');
			return element && element.textContent !== '%s';
		}
	`, selector, initialText)

	_, err = page.WaitForFunction(condition, nil)

	newText, err := page.Locator(selector).TextContent()
	if err != nil {
		return fmt.Errorf("err while attempting to read TextContent: %w", err)
	}

	if newText == initialText {
		return errors.New("expected initialText to not equal newText")
	}

	return nil
}
