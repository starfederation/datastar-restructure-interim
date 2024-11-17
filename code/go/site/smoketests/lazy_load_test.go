package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func lazyLoadExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		lazyLoadInitialState,
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

func lazyLoadInitialState(page playwright.Page) error {
	page.Locator("#lazy_load").GetByText("Loading...").Click()
	expect.Locator(page.Locator("img#lazy_load")).ToBeAttached()
	return nil
}
