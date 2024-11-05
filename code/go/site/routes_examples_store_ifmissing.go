package site

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

func setupExamplesStoreIfMissing(examplesRouter chi.Router) error {

	examplesRouter.Get("/store_ifmissing/updates", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		t := time.NewTicker(1 * time.Second)
		defer t.Stop()

		i := 1234
		for {
			select {
			case <-r.Context().Done():
				return
			case <-t.C:
				store := fmt.Sprintf("{id:%d}", i)

				switch i % 2 {
				case 0:
					fragment := fmt.Sprintf(`<div id="placeholder" data-store.ifmissing="%s" data-text="$id"></div>`, store)
					sse.RenderFragment(fragment, datastar.WithMergeUpsertAttributes())
				case 1:
					sse.MarshalAndPatchStoreIfMissing(store)
				}
				i++
			}
		}

	})

	return nil
}
