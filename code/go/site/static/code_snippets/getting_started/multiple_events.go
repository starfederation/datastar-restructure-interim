import (
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

sse := datastar.NewSSE(w,r)
sse.MergeFragments(`<div id="question">...</div>`)
sse.MergeFragments(`<div id="instructions">...</div>`)
sse.MergeSignals([]byte(`{answer: '...'}`)
sse.MergeSignals([]byte(`{prize: '...'}`)
