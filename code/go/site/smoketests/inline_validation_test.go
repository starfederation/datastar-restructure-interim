package smoketests

import (
	"context"
	"errors"

	playwright "github.com/playwright-community/playwright-go"
)

func inlineValidationExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		inlineValidationInitial,
		inlineValidationEmail,
		inlineValidationFirstName,
		inlineValidationLastName,
		inlineValidationSuccess,
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

func inlineValidationInitial(page playwright.Page) error {
	expect.Locator(page.Locator("#inline_validation")).ToContainText("Email Address")
	expect.Locator(page.Locator("#inline_validation")).ToContainText("First Name")
	expect.Locator(page.Locator("#inline_validation")).ToContainText("Last Name")
	expect.Locator(page.GetByTestId("validation_email")).ToHaveText("Email '' is already taken or is invalid. Please enter another email.")
	expect.Locator(page.GetByTestId("validation_firstName")).ToHaveText("First name must be at least 2 characters.")
	expect.Locator(page.GetByTestId("validation_lastName")).ToHaveText("Last name must be at least 2 characters.")
	expect.Locator(page.GetByTestId("submit_button")).ToBeAttached()
	return nil
}

func inlineValidationEmail(page playwright.Page) error {
	page.GetByTestId("input_email").PressSequentially("test@test.com")
	expect.Locator(page.GetByTestId("input_email")).ToHaveValue("test@test.com")
	expect.Locator(page.Locator("div#inline_validation")).Not().ToContainText("Email '' is already taken or is invalid. Please enter another email.")
	page.GetByTestId("input_email").Fill("")
	page.GetByTestId("input_email").PressSequentially("test")
	expect.Locator(page.GetByTestId("input_email")).ToHaveValue("test")
	expect.Locator(page.Locator("div#inline_validation")).ToContainText("Email 'test' is already taken or is invalid. Please enter another email.")
	return nil
}

func inlineValidationFirstName(page playwright.Page) error {
	page.GetByTestId("input_firstName").PressSequentially("Alexander")
	expect.Locator(page.GetByTestId("input_firstName")).ToHaveValue("Alexander")
	expect.Locator(page.Locator("div#inline_validation")).Not().ToContainText("First name must be at least 2 characters.")
	page.GetByTestId("input_firstName").Fill("")
	page.GetByTestId("input_firstName").PressSequentially("t")
	expect.Locator(page.GetByTestId("input_firstName")).ToHaveValue("t")
	expect.Locator(page.Locator("div#inline_validation")).ToContainText("First name must be at least 2 characters.")
	return nil
}

func inlineValidationLastName(page playwright.Page) error {
	page.GetByTestId("input_lastName").PressSequentially("Alexander")
	expect.Locator(page.GetByTestId("input_lastName")).ToHaveValue("Alexander")
	expect.Locator(page.Locator("div#inline_validation")).Not().ToContainText("Last name must be at least 2 characters.")
	page.GetByTestId("input_lastName").Fill("")
	page.GetByTestId("input_lastName").PressSequentially("t")
	expect.Locator(page.GetByTestId("input_lastName")).ToHaveValue("t")
	expect.Locator(page.Locator("div#inline_validation")).ToContainText("Last name must be at least 2 characters.")
	return nil
}

func inlineValidationSuccess(page playwright.Page) error {
	page.GetByTestId("input_email").PressSequentially("test@test.com")
	expect.Locator(page.GetByTestId("input_email")).ToHaveValue("test@test.com")
	expect.Locator(page.Locator("div#inline_validation")).Not().ToContainText("Email '' is already taken or is invalid. Please enter another email.")

	page.GetByTestId("input_firstName").PressSequentially("Alexander")
	expect.Locator(page.GetByTestId("input_firstName")).ToHaveValue("Alexander")
	expect.Locator(page.Locator("div#inline_validation")).Not().ToContainText("First name must be at least 2 characters.")

	page.GetByTestId("input_lastName").PressSequentially("Alexander")
	expect.Locator(page.GetByTestId("input_lastName")).ToHaveValue("Alexander")
	expect.Locator(page.Locator("div#inline_validation")).Not().ToContainText("Last name must be at least 2 characters.")

	page.GetByTestId("submit_button").Click()
	expect.Locator(page.Locator("#inline_validation")).ToContainText("Thank you for signing up!")
	return nil
}
