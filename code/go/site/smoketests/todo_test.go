package smoketests

import (
	"context"
	"errors"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

func todoExampleTest(ctx context.Context, page playwright.Page) error {
	suite := []func(page playwright.Page) error{
		todoInitialState,
		todoModes,
		todoAdd,
		todoDelete,
		todoClear,
	}

	errs := []error{}

	for _, test := range suite {
		page.GetByTestId("reset_todos").Click()
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

func todoAdd(page playwright.Page) error {
	input := page.GetByTestId("todos_input")
	input.Click()
	input.Fill("test")
	input.Press("Enter")
	expect.Locator(page.Locator("#todo4")).ToContainText("test")
	expect.Locator(page.GetByTestId("todo_count")).ToContainText("4 items")
	return nil
}

func todoClear(page playwright.Page) error {
	page.GetByTestId("clear_todos").Click()
	expect.Locator(page.Locator("#todo0")).ToContainText("Learn Datastar")
	expect.Locator(page.GetByTestId("todo_count")).ToContainText("3 items")
	todosList := page.GetByTestId("todos_list")
	liElements, err := todosList.Locator("li").All()

	if err != nil {
		return fmt.Errorf("err locating elements: %w", err)
	}

	if len(liElements) != 3 {
		return fmt.Errorf("expected 3, received %d", len(liElements))
	}

	return nil
}

func todoDelete(page playwright.Page) error {
	page.Locator("#todo2").Hover()
	todo := page.GetByTestId("delete_todo2")
	expect.Locator(todo).ToBeVisible()
	todo.Click()
	expect.Locator(page.GetByTestId("todo_count")).ToContainText("2 items")
	todosList := page.GetByTestId("todos_list")
	liElements, err := todosList.Locator("li").All()
	if err != nil {
		return fmt.Errorf("err locating elements: %w", err)
	}

	if len(liElements) != 3 {
		return fmt.Errorf("expected 3, received %d", len(liElements))
	}
	return nil
}

func todoInitialState(page playwright.Page) error {
	expect.Locator(page.GetByPlaceholder("What needs to be done?")).ToBeEmpty()
	expect.Locator(page.GetByTestId("toggle_all_todos")).ToBeAttached()
	expect.Locator(page.Locator("#todo0")).ToContainText("Learn any backend language")
	expect.Locator(page.Locator("#todo1")).ToContainText("Learn Datastar")
	expect.Locator(page.Locator("#todo2")).ToContainText("???")
	expect.Locator(page.Locator("#todo3")).ToContainText("Profit")
	expect.Locator(page.GetByTestId("todo_count")).ToContainText("3 items")
	expect.Locator(page.GetByTestId("All_mode")).ToBeAttached()
	expect.Locator(page.GetByTestId("Active_mode")).ToBeAttached()
	expect.Locator(page.GetByTestId("Completed_mode")).ToBeAttached()
	return nil
}

func todoModes(page playwright.Page) error {
	all := page.GetByTestId("All_mode")
	active := page.GetByTestId("Active_mode")
	completed := page.GetByTestId("Completed_mode")
	todo0 := page.Locator("#todo0")
	todo1 := page.Locator("#todo1")
	todo2 := page.Locator("#todo2")
	todo3 := page.Locator("#todo3")

	active.Click()
	expect.Locator(todo0).Not().ToBeAttached()
	expect.Locator(todo1).ToContainText("Learn Datastar")
	expect.Locator(todo2).ToContainText("???")
	expect.Locator(todo3).ToContainText("Profit")

	completed.Click()
	expect.Locator(todo0).ToContainText("Learn any backend language")
	expect.Locator(todo1).Not().ToBeAttached()
	expect.Locator(todo2).Not().ToBeAttached()
	expect.Locator(todo3).Not().ToBeAttached()

	all.Click()
	expect.Locator(todo0).ToContainText("Learn any backend language")
	expect.Locator(todo1).ToContainText("Learn Datastar")
	expect.Locator(todo2).ToContainText("???")
	expect.Locator(todo3).ToContainText("Profit")
	return nil
}

func toggleAll(page playwright.Page) error {
	page.GetByTestId("toggle_all_todos").Click()
	page.GetByTestId("clear_todos").Click()
	expect.Locator(page.GetByTestId("todos_list")).Not().ToBeAttached()
	expect.Locator(page.GetByTestId("todo_count")).Not().ToBeAttached()
	return nil
}
