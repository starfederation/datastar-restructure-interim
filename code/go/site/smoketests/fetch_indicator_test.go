package smoketests

import (
	"context"
	"errors"

	"github.com/playwright-community/playwright-go"
)

func fetchIndicatorExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		fetchIndicatorFirstRun,
		fetchIndicatorSecondRun,
		fetchIndicatorIsFetching,
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

func fetchIndicatorFirstRun(page playwright.Page) error {
	ind := page.Locator("#ind")
	greeting := page.Locator("#greeting")

	page.GetByTestId("greeting_button").Click()
	expect.Locator(ind).ToBeVisible()
	expect.Locator(greeting).ToBeVisible()
	expect.Locator(greeting).ToContainText("Hello, the time is")
	return nil
}

func fetchIndicatorSecondRun(page playwright.Page) error {
	ind := page.Locator("#ind")
	greeting := page.Locator("#greeting")
	greetingButton := page.GetByTestId("greeting_button")

	greetingButton.Click()
	expect.Locator(ind).ToBeVisible()
	expect.Locator(greeting).ToBeVisible()
	expect.Locator(greeting).ToContainText("Hello, the time is")
	greetingButton.Click()
	expect.Locator(ind).ToBeVisible()
	expect.Locator(greeting).ToBeVisible()
	expect.Locator(greeting).ToContainText("Hello, the time is")
	return nil
}

func fetchIndicatorIsFetching(page playwright.Page) error {
	greeting := page.Locator("#greeting")
	greetingButton := page.GetByTestId("greeting_button")

	greetingButton.Click()
	expect.Locator(greetingButton).ToBeDisabled()
	expect.Locator(greeting).ToBeVisible()
	expect.Locator(greeting).ToContainText("Hello, the time is")
	expect.Locator(greetingButton).ToBeEnabled()
	return nil
}
