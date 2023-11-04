package config

import (
	"log/slog"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/go-chi/httplog/v2"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

var ImageFormats = []string{"jpg", "jpeg", "gif", "png", "gif"}

var Markdown = goldmark.New(
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("github-dark"),
			highlighting.WithFormatOptions(
				chromahtml.WithCustomCSS(map[chroma.TokenType]string{
					chroma.PreWrapper: `padding: 10px;
						margin: 20px 0 20px 0;
						border-radius: 10px;
						box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3);
						overflow: auto;`,
				}),
			),
		),
	),
)

var Logger = httplog.NewLogger("htmx-blog", httplog.Options{
	LogLevel: slog.LevelDebug,
})
