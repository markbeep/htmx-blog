package main

import (
	"flag"
	"fmt"
	"os"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbeep/htmx-blog/internal/route"
	"github.com/markbeep/htmx-blog/internal/route/posts"
)

var (
	port         = flag.String("port", os.Getenv("PORT"), "port to host the website at")
	generate     = flag.Bool("generate", false, "if the html files should be generated")
	onlyGenerate = flag.Bool("only-generate", false, "if the program should end after generating files")
)

func main() {
	flag.Parse()
	if *port == "" {
		*port = "3000"
	}

	postsHander := posts.PostsHandler{}
	postsHander.GenerateHTML("content", "generated", *generate)
	if *onlyGenerate {
		return
	}

	r := chi.NewRouter()
	r.Use(route.MiddlewareLogging)

	r.Get("/", route.Index)
	r.Get("/health", route.Health)
	r.Get("/favicon.ico", route.Favicon)
	r.Get("/static/*", route.Static)
	r.Get("/about", route.About)
	r.Get("/polyring", route.Polyring)
	r.Get("/content/*", route.Content)
	r.Get("/*", route.Error404)
	r.Route("/posts", postsHander.Posts)

	host := fmt.Sprintf(":%s", *port)
	log.Printf("listening on %s", host)
	log.Fatal(http.ListenAndServe(host, r))
}
