package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func clickToEditExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		clickToEdittInitialState,
		clickToEditClickEdit,
		clickToEdittCancel,
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

func clickToEdittInitialState(page playwright.Page) error {
	expect.Locator(page.Locator("#contact_1")).ToContainText("First Name: John")
	expect.Locator(page.GetByText("First Name: John", playwright.PageGetByTextOptions{Exact: playwright.Bool(true)})).ToBeVisible()
	expect.Locator(page.Locator("#contact_1")).ToContainText("Last Name: Doe")
	expect.Locator(page.GetByText("Last Name: Doe", playwright.PageGetByTextOptions{Exact: playwright.Bool(true)})).ToBeVisible()
	expect.Locator(page.Locator("#contact_1")).ToContainText("Email: joe@blow.com")
	expect.Locator(page.GetByText("Email: joe@blow.com", playwright.PageGetByTextOptions{Exact: playwright.Bool(true)})).ToBeVisible()
	expect.Locator(page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Edit"})).ToBeVisible()
	expect.Locator(page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Reset"})).ToBeVisible()
	return nil
}

func clickToEditClickEdit(page playwright.Page) error {
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Edit"}).Click()
	page.GetByLabel("First Name").Click()
	page.GetByLabel("First Name").Fill("Foo")
	page.GetByLabel("Last Name").Click()
	page.GetByLabel("Last Name").Fill("Fighter")
	page.GetByLabel("Email").Click()
	page.GetByLabel("Email").Fill("foo@fighter.com")
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Save"}).Click()
	expect.Locator(page.Locator("#contact_1")).ToContainText("First Name: Foo")
	expect.Locator(page.Locator("#contact_1")).ToContainText("Last Name: Fighter")
	expect.Locator(page.Locator("#contact_1")).ToContainText("Email: foo@fighter.com")
	return nil
}

func clickToEdittCancel(page playwright.Page) error {
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Edit"}).Click()
	page.GetByLabel("First Name").Click()
	page.GetByLabel("First Name").Fill("Foo")
	page.GetByLabel("Last Name").Click()
	page.GetByLabel("Last Name").Fill("Fighter")
	page.GetByLabel("Email").Click()
	page.GetByLabel("Email").Fill("foo@fighter.com")
	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Cancel"}).Click()
	expect.Locator(page.Locator("#contact_1")).ToContainText("First Name: John")
	page.GetByText("Last Name: Doe", playwright.PageGetByTextOptions{Exact: playwright.Bool(true)}).Click()
	expect.Locator(page.Locator("#contact_1")).ToContainText("Last Name: Doe")
	expect.Locator(page.Locator("#contact_1")).ToContainText("Email: joe@blow.com")
	return nil
}
