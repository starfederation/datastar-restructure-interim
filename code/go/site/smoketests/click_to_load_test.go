package smoketests

import (
	"context"
	"errors"

	"github.com/playwright-community/playwright-go"
)

func clickToLoadExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		clickToLoadInitialState,
		clickToLoadClickAndLoad,
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

func clickToLoadInitialState(page playwright.Page) error {
	expect.Locator(page.Locator("#click_to_load_rows tr")).ToHaveCount(10)
	expect.Locator(page.Locator("#agent_0 > td:first-child")).ToContainText("Agent Smith")
	expect.Locator(page.Locator("#agent_0 > td:nth-child(2)")).ToContainText("void1@null.org")
	expect.Locator(page.Locator("#agent_0 > td:nth-child(3)")).ToBeAttached()
	expect.Locator(page.Locator("#agent_9 > td:first-child")).ToContainText("Agent Smith")
	expect.Locator(page.Locator("#agent_9 > td:nth-child(2)")).ToContainText("void10@null.org")
	expect.Locator(page.Locator("#agent_9 > td:nth-child(3)")).ToBeAttached()
	return nil
}

func clickToLoadClickAndLoad(page playwright.Page) error {
	expect.Locator(page.Locator("#click_to_load_rows tr")).ToHaveCount(10)
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Load More"}).Click()
	expect.Locator(page.Locator("#click_to_load_rows tr")).ToHaveCount(20)
	expect.Locator(page.Locator("#agent_0> td:first-child")).ToContainText("Agent Smith")
	expect.Locator(page.Locator("#agent_0 > td:nth-child(2)")).ToContainText("void1@null.org")
	expect.Locator(page.Locator("#agent_0 > td:nth-child(3)")).ToBeAttached()
	expect.Locator(page.Locator("#agent_10 > td:first-child")).ToContainText("Agent Smith")
	expect.Locator(page.Locator("#agent_10 > td:nth-child(2)")).ToContainText("void11@null.org")
	expect.Locator(page.Locator("#agent_10 > td:nth-child(3)")).ToBeAttached()
	expect.Locator(page.Locator("#agent_19 > td:first-child")).ToContainText("Agent Smith")
	expect.Locator(page.Locator("#agent_19 > td:nth-child(2)")).ToContainText("void20@null.or")
	expect.Locator(page.Locator("#agent_19 > td:nth-child(3)")).ToBeAttached()
	return nil
}
