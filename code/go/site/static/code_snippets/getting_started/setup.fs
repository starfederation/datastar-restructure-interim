// This is not valid F# code yet!!!
import (
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

// Creates a new `ServerSentEventGenerator` instance.
sse := datastar.NewSSE(w,r)

// Merges the HTML fragment into the DOM.
sse.MergeFragments([]byte("<div id="question">What do you put in a toaster?</div>"))

// Merges the `answer` value into the signals.
sse.MergeSignals([]byte("['answer' => 'bread']"))