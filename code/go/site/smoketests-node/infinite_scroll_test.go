package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func infiniteScrollExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		infiniteScrollInitial,
		infiniteScrollRun,
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

func infiniteScrollInitial(page playwright.Page) error {
	expect.Locator(page.Locator("#click_to_load_rows tr")).ToHaveCount(10)
	expect.Locator(page.Locator("#agent_0 > td:first-child")).ToContainText("Agent Smith 0")
	expect.Locator(page.Locator("#agent_0 > td:nth-child(2)")).ToContainText("void1@null.org")
	expect.Locator(page.Locator("#agent_0 > td:nth-child(3)")).ToBeAttached()
	expect.Locator(page.Locator("#agent_9 > td:first-child")).ToContainText("Agent Smith 9")
	expect.Locator(page.Locator("#agent_9 > td:nth-child(2)")).ToContainText("void10@null.org")
	expect.Locator(page.Locator("#agent_9 > td:nth-child(3)")).ToBeAttached()
	return nil
}

func infiniteScrollRun(page playwright.Page) error {
	lm := page.Locator("#loading_message")
	expect.Locator(lm).ToContainText("Loading...")

	rows := page.Locator("#click_to_load_rows tr")
	for i := 10; i <= 50; i += 10 {
		expect.Locator(rows).ToHaveCount(i)
		lm.ScrollIntoViewIfNeeded()
	}
	return nil
}
