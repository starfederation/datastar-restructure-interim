package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func bulkUpdateExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		bulkUpdateInitialState,
		bulkUpdateIndividualActivate,
		bulkUpdateIndividualDeactivate,
		bulkUpdateActivate,
		bulkUpdateDeactivate,
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

func checkboxLocator(page playwright.Page, name string) playwright.Locator {
	return page.GetByRole("row", playwright.PageGetByRoleOptions{Name: name}).Locator(".checkbox")
}

func bulkUpdateInitialState(page playwright.Page) error {
	expect.Locator(page.Locator("#contact_0")).ToContainText("Joe Smith")
	expect.Locator(page.Locator("#contact_0")).ToContainText("joe@smith.org")
	expect.Locator(page.Locator("#contact_0")).ToBeAttached()
	expect.Locator(page.Locator("#contact_1")).ToContainText("Angie MacDowell")
	expect.Locator(page.Locator("#contact_1")).ToContainText("angie@macdowell.org")
	expect.Locator(page.Locator("#contact_1")).ToBeAttached()
	expect.Locator(page.Locator("#contact_2")).ToContainText("Fuqua Tarkenton")
	expect.Locator(page.Locator("#contact_2")).ToContainText("fuqua@tarkenton.org")
	expect.Locator(page.Locator("#contact_2")).ToBeAttached()
	expect.Locator(page.Locator("#contact_3")).ToContainText("Kim Yee")
	expect.Locator(page.Locator("#contact_3")).ToContainText("kim@yee.org")
	expect.Locator(page.Locator("#contact_3")).ToBeAttached()
	return nil
}

func bulkUpdateIndividualActivate(page playwright.Page) error {
	joe := checkboxLocator(page, "Joe Smith joe@smith.org")
	joe.Click()
	expect.Locator(joe).ToBeChecked()
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Activate", Exact: playwright.Bool(true)}).Click()
	expect.Locator(page.Locator("#contact_0")).ToContainText("Active")
	angie := checkboxLocator(page, "Angie MacDowell angie@macdowell.org")
	angie.Click()
	expect.Locator(angie).ToBeChecked()
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Activate", Exact: playwright.Bool(true)}).Click()
	expect.Locator(page.Locator("#contact_1")).ToContainText("Active")
	return nil
}

func bulkUpdateIndividualDeactivate(page playwright.Page) error {

	joe := page.
		GetByRole("row", playwright.PageGetByRoleOptions{Name: "Joe Smith joe@smith.org"}).
		Locator(".checkbox")
	joe.Click()
	expect.Locator(joe).ToBeChecked()
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Deactivate"}).Click()
	expect.Locator(page.Locator("#contact_0")).ToContainText("Inactive")

	angie := checkboxLocator(page, "Angie MacDowell angie@macdowell.org")
	angie.Click()
	expect.Locator(angie).ToBeChecked()
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Deactivate", Exact: playwright.Bool(true)}).Click()
	expect.Locator(page.Locator("#contact_1")).ToContainText("Inactive")
	return nil
}

func bulkUpdateActivate(page playwright.Page) error {
	checkAll := page.GetByRole("row", playwright.PageGetByRoleOptions{Name: "Name Email Status"}).Locator(".checkbox")
	checkAll.Click()
	expect.Locator(checkAll).ToBeChecked()
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Activate", Exact: playwright.Bool(true)}).Click()
	expect.Locator(page.Locator("#contact_0")).ToContainText("Active")
	expect.Locator(page.Locator("#contact_1")).ToContainText("Active")
	expect.Locator(page.Locator("#contact_2")).ToContainText("Active")
	expect.Locator(page.Locator("#contact_3")).ToContainText("Active")
	return nil
}

func bulkUpdateDeactivate(page playwright.Page) error {
	checkAll := page.GetByRole("row", playwright.PageGetByRoleOptions{Name: "Name Email Status"}).Locator(".checkbox")
	checkAll.Click()
	expect.Locator(checkAll).ToBeChecked()
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Deactivate", Exact: playwright.Bool(true)}).Click()
	expect.Locator(page.Locator("#contact_0")).ToContainText("Inactive")
	expect.Locator(page.Locator("#contact_1")).ToContainText("Inactive")
	expect.Locator(page.Locator("#contact_2")).ToContainText("Inactive")
	expect.Locator(page.Locator("#contact_3")).ToContainText("Inactive")
	return nil
}
