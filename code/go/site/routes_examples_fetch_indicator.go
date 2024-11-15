package site

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

func setupExamplesFetchIndicator(examplesRouter chi.Router) error {

	examplesRouter.Get("/fetch_indicator/greet", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		sse.MergeFragmentTempl(fetchIndicatorEmpty())
		time.Sleep(2 * time.Second)
		sse.MergeFragmentTempl(fetchIndicatorGreeting())
	})

	return nil
}
