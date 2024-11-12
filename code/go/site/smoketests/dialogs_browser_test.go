package smoketests

import (
	"context"
	"errors"

	"github.com/playwright-community/playwright-go"
)

func dialogsBrowserExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		dialogsBrowserAccept,
		dialogsBrowserCancel,
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

func dialogsBrowserAccept(page playwright.Page) error {
	page.On("dialog", func(dialog playwright.Dialog) {
		dialog.Accept()
	})
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Click Me"}).Click()
	expect.Locator(page.Locator("#dialogs")).ToContainText("You clicked the button and confirmed with prompt of")
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Reset"}).Click()
	expect.Locator(page.Locator("#dialogs")).Not().ToContainText("You clicked the button and confirmed with prompt of")
	return nil
}

func dialogsBrowserCancel(page playwright.Page) error {
	page.On("dialog", func(dialog playwright.Dialog) {
		dialog.Dismiss()
	})
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Click Me"}).Click()
	expect.Locator(page.Locator("#dialogs")).Not().ToContainText("You clicked the button and confirmed with prompt of")
	return nil
}
