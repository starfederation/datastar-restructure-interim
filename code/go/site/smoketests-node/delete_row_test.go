package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func deleteRowExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		deleteRowInitialState,
		deleteRowDeleteRow,
	}

	errs := []error{}

	for _, test := range suite {
		page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Reset"}).Click()
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

func deleteRowInitialState(page playwright.Page) error {
	l := page.Locator("#contact_0")
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Joe Smith"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "joe@smith.org"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Active"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"})).ToBeAttached()

	l = page.Locator("#contact_1")
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Angie MacDowell"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "angie@macdowell.org"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Active"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"})).ToBeAttached()

	l = page.Locator("#contact_2")
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Fuqua Tarkenton"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "fuqua@tarkenton.org"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Active"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"})).ToBeAttached()

	l = page.Locator("#contact_3")
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Kim Yee"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "kim@yee.org"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Active"})).ToBeAttached()
	expect.Locator(l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"})).ToBeAttached()

	return nil
}

func deleteRowDeleteRow(page playwright.Page) error {
	// 1st Row.
	page.Once("dialog", func(dialog playwright.Dialog) {
		// fmt.Printf("Dialog message: %s\n", dialog.Message())
		dialog.Accept()
	})
	l := page.Locator("#contact_0")
	l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"}).Click()
	expect.Locator(l).Not().ToBeAttached()

	// 2nd Row.
	page.Once("dialog", func(dialog playwright.Dialog) {
		// fmt.Printf("Dialog message: %s\n", dialog.Message())
		dialog.Accept()
	})
	l = page.Locator("#contact_1")
	l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"}).Click()
	expect.Locator(l).Not().ToBeAttached()

	// 3rd Row.
	page.Once("dialog", func(dialog playwright.Dialog) {
		// fmt.Printf("Dialog message: %s\n", dialog.Message())
		dialog.Accept()
	})
	l = page.Locator("#contact_2")
	l.GetByRole("cell", playwright.LocatorGetByRoleOptions{Name: "Delete"}).Click()
	expect.Locator(l).Not().ToBeAttached()

	return nil
}
