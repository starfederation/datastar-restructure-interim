package smoketests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func rafUpdateExampleTest(t *testing.T, ctx context.Context) error {
	page := pageFromURL(t, ctx, "examples/raf_update")

	timeEl := page.MustElement("#time")

	last := ""
	for i := 0; i < 3; i++ {
		current, err := timeEl.Text()
		assert.NoError(t, err)

		assert.NotEqual(t, last, current)
		last = current
		time.Sleep(1 * time.Second)
	}

	return nil
}
