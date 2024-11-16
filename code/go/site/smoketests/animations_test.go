package smoketests

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	playwright "github.com/playwright-community/playwright-go"
)

func animationsExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		animationsColorThrob,
		animationsClickAndFadeOut,
		animationsClickAndFadeIn,
		animationsInFlightIndicator,
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

func animationsColorThrob(page playwright.Page) error {
	selector := "#color_throb"

	initialText, err := page.Locator(selector).TextContent()
	if err != nil {
		return fmt.Errorf("err while attempting to read TextContent: %w", err)
	}

	condition := fmt.Sprintf(`
			() => {
				const element = document.querySelector('%s');
				return element && element.textContent !== '%s';
			}
		`, selector, initialText)

	_, err = page.WaitForFunction(condition, nil)

	if err != nil {
		return fmt.Errorf("err encountered in WaitForFunction: %w", err)
	}

	newText, err := page.Locator(selector).TextContent()
	if err != nil {
		return fmt.Errorf("err while attempting to read TextContent: %w", err)
	}

	if newText == initialText {
		return errors.New("expected initialText to not equal newText")
	}

	return nil
}

func animationsClickAndFadeOut(page playwright.Page) error {
	button := page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Fade out then delete on click"})
	expect.Locator(button).ToBeAttached()
	button.Click()
	expect.Locator(button).Not().ToBeAttached()
	return nil
}

func animationsClickAndFadeIn(page playwright.Page) error {
	button := page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Fade me in on click"})
	button.Click()

	expect.Locator(button).ToHaveClass(regexp.MustCompile("/.*datastar-swapping.*datastar-settling.*/"))
	expect.Locator(button).Not().ToHaveClass(regexp.MustCompile("/.*datastar-swapping.*datastar-settling.*/"))
	return nil
}

func animationsInFlightIndicator(page playwright.Page) error {
	expect.Locator(page.GetByRole("textbox")).ToBeEmpty()
	page.GetByRole("textbox").Fill("test")

	expect.Locator(page.Locator("div#request_in_flight_indicator")).Not().ToHaveClass(regexp.MustCompile("/.*datastar-indicator-loading.*/"))
	expect.Locator(page.Locator("div#request_in_flight")).Not().ToContainText("Submitted!")

	page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Submit"}).Click()
	expect.Locator(page.Locator("div#request_in_flight_indicator")).ToHaveClass(regexp.MustCompile("/.*datastar-indicator-loading.*/"))

	expect.Locator(page.Locator("#request_in_flight")).ToContainText("Submitted!")
	return nil
}
