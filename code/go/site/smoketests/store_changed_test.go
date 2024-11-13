package smoketests

import (
	"context"
	"errors"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

func storeChangedExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		storeChangedLocalUpdate,
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

func storeChangedLocalUpdate(page playwright.Page) error {
	increment := page.Locator("#increment")
	clear := page.Locator("#clear")
	local_clicks := page.Locator("#local_clicks")
	serverChanged := page.Locator("#from_server")

	clear.Click()
	expect.Locator(local_clicks).ToContainText("0")

	for i := 0; i < 3; i++ {
		expect.Locator(local_clicks).ToContainText(fmt.Sprintf("%d", i))
		increment.Click()
		next := fmt.Sprintf("%d", i+1)
		expect.Locator(local_clicks).ToContainText(next)
		expect.Locator(serverChanged).ToContainText(next)
	}
	clear.Click()
	expect.Locator(local_clicks).ToContainText("0")

	return nil
}
