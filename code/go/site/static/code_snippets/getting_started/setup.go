import (
	"fmt"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

// Get a random question and answer from somewhere in your code.
QA := getRandomQuestionAnswer()

// Creates a new `ServerSentEventGenerator` instance.
sse := datastar.NewSSE(w,r)

// Merges HTML fragments into the DOM.
sse.MergeFragments(
    fmt.Sprintf(
		`<div id="question">%s</div>`,
		QA.question,
	)
)

// Merges signals into the store.
sse.MergeSignals(
	fmt.Sprintf(`{response: '', answer: '%s'}`,
	QA.answer),
)
