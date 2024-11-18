package smoketests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/starfederation/datastar/code/go/site"
	"github.com/stretchr/testify/require"
)

const port = 4321

type RodKey string

func rodFromContext(t *testing.T, ctx context.Context) *rod.Browser {
	b := ctx.Value(RodKey("rod")).(*rod.Browser)
	require.NotNil(t, b)
	return b
}

func pageFromURL(t *testing.T, ctx context.Context, url string) *rod.Page {
	b := rodFromContext(t, ctx)
	fullURL := fmt.Sprintf("http://localhost:%d/%s", port, url)
	p := b.MustPage(fullURL)
	require.NotNil(t, p)
	return p
}

func TestSite(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx = context.WithValue(ctx, RodKey("rod"), rod.New().MustConnect())

	log.Printf("running site on port %d", port)
	defer log.Print("closing site")
	go site.RunBlocking(port)(ctx)
	time.Sleep(1 * time.Second)

	t.Run("smoke tests", func(t *testing.T) {
		t.Run("active search example", func(t *testing.T) {
			activeSearchExampleTest(t, ctx)
		})

		t.Run("raf update example", func(t *testing.T) {
			rafUpdateExampleTest(t, ctx)
		})
	})
}
