package smoketests

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
)

func rafUpdateExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		rafUpdateTimeUpdate,
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

func rafUpdateTimeUpdate(page playwright.Page) error {

	last := ""
	for i := 0; i < 3; i++ {
		current, err := page.Locator("#time").InnerText()
		if err != nil {
			return fmt.Errorf("could not read InnerText: %w", err)
		}

		if current == last {
			return errors.New("expected current time to not equal last time")
		}

		last = current
		time.Sleep(1 * time.Second)
	}

	return nil
}
