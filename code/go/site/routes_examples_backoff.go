package site

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

func setupExamplesBackoff(examplesRouter chi.Router) error {

	examplesRouter.Put("/backoff/notValid", func(w http.ResponseWriter, r *http.Request) {

		sse := datastar.NewSSE(w, r)
		sse.Send("datastar-foo", []string{"bar"})
	})

	return nil
}
