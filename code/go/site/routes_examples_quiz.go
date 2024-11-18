package site

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

type QA struct {
	Question string
	Answer   string
}

var qaList = []QA{
	{"What do you put in a toaster?", "Bread"},
	{"How many months have 28 days?", "Twelve"},
	{"If you’re running in a race and pass the person in second place, what place are you in?", "Second"},
	{"What do you get if you divide 30 by half and add 10.", "Seventy"},
	{"What gets wetter the more it dries?", "Towel"},
}

func randomQA() QA {
	rand.Seed(time.Now().UnixNano())
	return qaList[rand.Intn(len(qaList))]
}

func setupExamplesQuiz(examplesRouter chi.Router) error {

	examplesRouter.Get("/quiz/data", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		QA := randomQA()
		sse.MergeFragments(fmt.Sprintf(`<div id="question2">%s</div>`, QA.Question))
		sse.MarshalAndMergeSignals(map[string]any{
			"answer": QA.Answer,
		})
	})

	return nil
}
