package site

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

func setupExamplesShow(examplesRouter chi.Router) error {
	examplesRouter.Get("/show/data", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		store := &ShowSignals{
			BindBool: false,
		}

		sse.MergeFragmentTempl(ShowView(store))
	})

	return nil
}
