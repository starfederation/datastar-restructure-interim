package smoketests

import (
	"context"
	"log"

	"github.com/playwright-community/playwright-go"
)

func animations(ctx context.Context, page playwright.Page) error {
	colerThrob := page.Locator("#coler_throb")
	page.WaitForFunction(`
	({ selector, initialText }) => {
          const element = document.querySelector(selector);
          return element && element.textContent !== initialText;
        },
	`, map[string]any{
		"selector":    "#coler_throb",
		"initialText": "Coler Throb",
	})
	for i := 0; i < 3; i++ {
		text, err := colerThrob.
		if err != nil {
			return err
		}

		log.Print(text)
	}

	return nil
}
