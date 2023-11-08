package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/markbeep/htmx-blog/components"
	"github.com/markbeep/htmx-blog/internal/route"
)

var (
	port = flag.String("port", os.Getenv("PORT"), "port to host the website at")
)

func main() {
	flag.Parse()
	if *port == "" {
		*port = "3000"
	}

	postsHander := route.PostsHandler{}
	postsHander.ConvertMarkdown("content")

	r := chi.NewRouter()
	r.Use(route.MiddlewareLogging)

	r.Get("/", templ.Handler(components.Index()).ServeHTTP)
	r.Get("/health", templ.Handler(components.Health()).ServeHTTP)
	r.Get("/about", templ.Handler(components.About()).ServeHTTP)
	r.Get("/polyring", templ.Handler(components.Polyring()).ServeHTTP)
	r.Get("/*", templ.Handler(components.Error404()).ServeHTTP)

	r.Get("/favicon.ico", route.Favicon)
	r.Get("/static/*", route.Static)
	r.Get("/content/*", route.Content)
	r.Get("/robots.txt", route.Robots)
	r.Get("/posts/index.xml", route.XML(postsHander.GetPosts()))

	r.Get("/posts", templ.Handler(components.Posts(postsHander.GetPosts())).ServeHTTP)

	// Add all of the posts
	for _, p := range postsHander.GetPosts() {
		post := *p
		t := templ.ComponentFunc(func(_ context.Context, w io.Writer) (err error) {
			_, err = w.Write(post.Buffer)
			return
		})
		r.Get(post.Path, templ.Handler(components.Post(post, t)).ServeHTTP)
	}

	host := fmt.Sprintf(":%s", *port)
	log.Printf("listening on %s", host)
	log.Fatal(http.ListenAndServe(host, r))
}
