package site

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/benbjohnson/hashfs"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gorilla/sessions"
	datastar "github.com/starfederation/datastar/code/go/sdk"
)

//go:embed static/*
var staticFS embed.FS

var (
	staticSys    = hashfs.NewFS(staticFS)
	highlightCSS templ.Component
	mdRenderer   func() *mdhtml.Renderer
)

func staticPath(path string) string {
	return "/" + staticSys.HashName("static/"+path)
}

func RunBlocking(port int) toolbelt.CtxErrFunc {
	upsertBuildSize()

	return func(ctx context.Context) error {

		router := chi.NewRouter()

		router.Use(middleware.Recoverer)

		cleanup, err := setupRoutes(ctx, router)
		defer cleanup()
		if err != nil {
			return fmt.Errorf("error setting up routes: %w", err)
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}

		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		return srv.ListenAndServe()

	}
}

func setupRoutes(ctx context.Context, router chi.Router) (cleanup func() error, err error) {
	defer router.Handle("/static/*", hashfs.FileServer(staticSys))
	defer router.Get("/hotreload", func(w http.ResponseWriter, r *http.Request) {
		datastar.NewSSE(w, r)
		<-r.Context().Done()
	})

	ns, err := embeddednats.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating embedded nats server: %w", err)
	}
	ns.WaitForServer()

	cleanup = func() error {
		return errors.Join(
			ns.Close(),
		)
	}

	htmlFormatter := html.New(html.WithClasses(true), html.TabWidth(2))
	if htmlFormatter == nil {
		return cleanup, fmt.Errorf("couldn't create html formatter")
	}
	styleName := "nord"
	highlightStyle := styles.Get(styleName)
	if highlightStyle == nil {
		return cleanup, fmt.Errorf("couldn't find style %s", styleName)
	}
	highlightCSSBuffer := &bytes.Buffer{}
	if err := htmlFormatter.WriteCSS(highlightCSSBuffer, highlightStyle); err != nil {
		return cleanup, fmt.Errorf("error writing highlight css: %w", err)
	}
	highlightCSS = templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, fmt.Sprintf(`<style>%s</style>`, highlightCSSBuffer.String()))
		return err
	})

	mdRenderer = func() *mdhtml.Renderer {
		return mdhtml.NewRenderer(mdhtml.RendererOptions{
			Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank,
			RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {

				// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
				htmlHighlight := func(w io.Writer, source, lang, defaultLang string) error {
					if lang == "" {
						lang = defaultLang
					}
					l := lexers.Get(lang)
					if l == nil {
						l = lexers.Analyse(source)
					}
					if l == nil {
						l = lexers.Fallback
					}
					l = chroma.Coalesce(l)

					it, err := l.Tokenise(nil, source)
					if err != nil {
						return err
					}
					return htmlFormatter.Format(w, highlightStyle, it)
				}

				renderCode := func(w io.Writer, codeBlock *ast.CodeBlock) {
					defaultLang := ""
					lang := string(codeBlock.Info)
					htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
				}

				if code, ok := node.(*ast.CodeBlock); ok {
					renderCode(w, code)
					return ast.GoToNext, true
				}
				return ast.GoToNext, false
			},
		})
	}

	sessionStore := sessions.NewCookieStore([]byte("datastar-session-secret"))
	sessionStore.MaxAge(int(24 * time.Hour / time.Second))

	if err := errors.Join(
		setupHome(router, sessionStore, ns),
		setupGuide(router),
		setupReferenceRoutes(router),
		setupExamples(ctx, router, sessionStore, ns),
		setupEssays(router),
		setupMemes(router),
	); err != nil {
		return cleanup, fmt.Errorf("error setting up routes: %w", err)
	}

	return cleanup, nil
}
