import (
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

// Get a random question and answer from somewhere in your code.
QA := getRandomQuestionAnswer()

// Creates a new `ServerSentEventGenerator` instance.
sse := datastar.NewSSE(w,r)

// Merges the HTML fragment into the DOM.
sse.MergeFragments(
    fmt.Sprintf(`<div id="question">%s</div>`, QA.question)
)

// Merges the `answer` signal into the store.
sse.MergeSignals(fmt.Sprintf(`{answer: '%s'}`, QA.answer))
