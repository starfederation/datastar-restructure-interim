package smoketests

import (
	"context"
	"errors"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

func valueSelectExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		valueSelectInitialState,
		valueSelectFirstSelect,
		valueSelectSecondSelect,
		valueSelectTest,
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

func valueSelectInitialState(page playwright.Page) error {
	expect.Locator(page.GetByTestId("make_select")).ToBeAttached()
	expect.Locator(page.GetByTestId("make_select")).ToContainText("Select a Make")
	expect.Locator(page.GetByTestId("model_select")).Not().ToBeAttached()

	return nil
}

func valueSelectFirstSelect(page playwright.Page) error {
	// get audi value
	selectionAudi, err := page.GetByTestId("make_option_Audi").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if selectionAudi == "" {
		return errors.New("expected selectionAudi to not be empty")
	}

	// get toyota value
	selectionToyota, err := page.GetByTestId("make_option_Toyota").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if selectionToyota == "" {
		return errors.New("expected selectionToyota to not be empty")
	}

	// get ford value
	selectionFord, err := page.GetByTestId("make_option_Ford").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if selectionFord == "" {
		return errors.New("expected selectionFord to not be empty")
	}

	// set select to audi value
	page.GetByTestId("make_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(selectionAudi)})
	err = expect.Locator(page.GetByTestId("make_select")).ToContainText("Audi")
	if err != nil {
		return fmt.Errorf("'make_select' did not contain text 'Audi': %w", err)
	}
	err = expect.Locator(page.GetByTestId("make_select")).ToHaveValue(selectionAudi)
	if err != nil {
		return fmt.Errorf("'make_select' did not have value 'Audi': %w", err)
	}
	err = expect.Locator(page.GetByTestId("model_select")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("'make_select' was not attached: %w", err)
	}

	// set select to toyota value
	page.GetByTestId("make_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(selectionToyota)})
	err = expect.Locator(page.GetByTestId("make_select")).ToContainText("Toyota")
	if err != nil {
		return fmt.Errorf("'make_select' did not contain text 'Toyota': %w", err)
	}
	err = expect.Locator(page.GetByTestId("make_select")).ToHaveValue(selectionToyota)
	if err != nil {
		return fmt.Errorf("'make_select' did not have value 'Toyota': %w", err)
	}
	err = expect.Locator(page.GetByTestId("model_select")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("'make_select' was not attached: %w", err)
	}

	// set select to ford value
	page.GetByTestId("make_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(selectionFord)})
	err = expect.Locator(page.GetByTestId("make_select")).ToContainText("Ford")
	if err != nil {
		return fmt.Errorf("'make_select' did not contain text 'Ford': %w", err)
	}
	err = expect.Locator(page.GetByTestId("make_select")).ToHaveValue(selectionFord)
	if err != nil {
		return fmt.Errorf("'make_select' did not have value 'Ford': %w", err)
	}
	err = expect.Locator(page.GetByTestId("model_select")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("'make_select' was not attached: %w", err)
	}

	return nil
}

func valueSelectSecondSelect(page playwright.Page) error {
	// get value for option 'Audi'
	makeOptionAudiValue, err := page.GetByTestId("make_option_Audi").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if makeOptionAudiValue == "" {
		return errors.New("expected makeOptionAudiValue to not be empty")
	}

	// set make_select to the 'Audi' option value
	page.GetByTestId("make_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(makeOptionAudiValue)})
	err = expect.Locator(page.GetByTestId("model_select")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("'model_select' was not attached: %w", err)
	}

	// Get the model select values
	modelOptionA1, err := page.GetByTestId("model_option_A1").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if modelOptionA1 == "" {
		return errors.New("expected modelOptionA1 to not be empty")
	}

	modelOptionA3, err := page.GetByTestId("model_option_A3").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if modelOptionA3 == "" {
		return errors.New("expected modelOptionA3 to not be empty")
	}

	modelOptionA6, err := page.GetByTestId("model_option_A6").GetAttribute("value")
	if err != nil {
		return fmt.Errorf("could not get attribute: %w", err)
	}

	if modelOptionA6 == "" {
		return errors.New("expected modelOptionA6 to not be empty")
	}

	page.GetByTestId("model_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(modelOptionA1)})

	err = expect.Locator(page.GetByTestId("model_select")).ToContainText("A1")
	if err != nil {
		return fmt.Errorf("model_select did not contain text 'A1': %w", err)
	}

	err = expect.Locator(page.GetByTestId("model_select")).ToHaveValue(modelOptionA1)
	if err != nil {
		return fmt.Errorf("model_select did not have value '%s': %w", modelOptionA1, err)
	}

	err = expect.Locator(page.GetByTestId("select_button")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("select_button was not attached: %w", err)
	}

	page.GetByTestId("model_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(modelOptionA3)})

	err = expect.Locator(page.GetByTestId("model_select")).ToContainText("A3")
	if err != nil {
		return fmt.Errorf("model_select did not contain text 'A3': %w", err)
	}

	err = expect.Locator(page.GetByTestId("model_select")).ToHaveValue(modelOptionA3)
	if err != nil {
		return fmt.Errorf("model_select did not have value '%s': %w", modelOptionA3, err)
	}

	err = expect.Locator(page.GetByTestId("select_button")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("select_button was not attached: %w", err)
	}

	page.GetByTestId("model_select").SelectOption(playwright.SelectOptionValues{Values: playwright.StringSlice(modelOptionA6)})

	err = expect.Locator(page.GetByTestId("model_select")).ToContainText("A6")
	if err != nil {
		return fmt.Errorf("model_select did not contain text 'A6': %w", err)
	}

	err = expect.Locator(page.GetByTestId("model_select")).ToHaveValue(modelOptionA6)
	if err != nil {
		return fmt.Errorf("model_select did not have value '%s': %w", modelOptionA6, err)
	}

	err = expect.Locator(page.GetByTestId("select_button")).ToBeAttached()
	if err != nil {
		return fmt.Errorf("select_button was not attached: %w", err)
	}

	return nil
}

func valueSelectTest(page playwright.Page) error {
	return nil
}
