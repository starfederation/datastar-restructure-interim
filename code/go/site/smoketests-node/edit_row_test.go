package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func editRowExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		editRowInitialState,
		editRowEditAndSave,
		editRowEditAndCancel,
	}

	errs := []error{}

	for _, test := range suite {
		page.GetByTestId("reset").Click()
		if err := test(page); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func editRowInitialState(page playwright.Page) error {
	// 1
	expect.Locator(page.Locator("#contact_0")).ToContainText("Joe Smith")
	expect.Locator(page.Locator("#contact_0")).ToContainText("joe@smith.org")
	expect.Locator(page.GetByTestId("contact_0_edit")).ToBeAttached()
	// 2
	expect.Locator(page.Locator("#contact_1")).ToContainText("Angie MacDowell")
	expect.Locator(page.Locator("#contact_1")).ToContainText("angie@macdowell.org")
	expect.Locator(page.GetByTestId("contact_1_edit")).ToBeAttached()
	// 3
	expect.Locator(page.Locator("#contact_2")).ToContainText("Fuqua Tarkenton")
	expect.Locator(page.Locator("#contact_2")).ToContainText("fuqua@tarkenton.org")
	expect.Locator(page.GetByTestId("contact_2_edit")).ToBeAttached()
	// 4
	expect.Locator(page.Locator("#contact_3")).ToContainText("Kim Yee")
	expect.Locator(page.Locator("#contact_3")).ToContainText("kim@yee.org")
	expect.Locator(page.GetByTestId("contact_3_edit")).ToBeAttached()
	return nil
}

func editRowEditAndSave(page playwright.Page) error {
	page.GetByTestId("contact_0_edit").Click()
	page.GetByTestId("contact_0_name").Fill("Foo Fighter")
	page.GetByTestId("contact_0_email").Fill("foo@fighter.org")
	page.GetByTestId("contact_0_save").Click()
	expect.Locator(page.Locator("#contact_0")).ToContainText("Foo Fighter")
	expect.Locator(page.Locator("#contact_0")).ToContainText("foo@fighter.org")
	expect.Locator(page.GetByTestId("contact_0_edit")).ToBeAttached()
	return nil
}

func editRowEditAndCancel(page playwright.Page) error {
	page.GetByTestId("contact_0_edit").Click()
	page.GetByTestId("contact_0_name").Fill("Foo Fighter")
	page.GetByTestId("contact_0_email").Fill("foo@fighter.org")
	page.GetByTestId("contact_0_cancel").Click()
	expect.Locator(page.Locator("#contact_0")).ToContainText("Joe Smith")
	expect.Locator(page.Locator("#contact_0")).ToContainText("joe@smith.org")
	expect.Locator(page.GetByTestId("contact_0_edit")).ToBeAttached()

	return nil
}
