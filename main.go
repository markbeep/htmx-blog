package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/markbeep/htmx-blog/internal/config"
	"github.com/markbeep/htmx-blog/internal/route"
	"github.com/markbeep/htmx-blog/internal/route/posts"
)

var (
	port = flag.String("port", os.Getenv("PORT"), "port to host the website at")
)

func main() {
	flag.Parse()
	if *port == "" {
		*port = "3000"
	}

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(config.Logger))

	postsHander := posts.PostsHandler{}
	postsHander.GenerateHTML("content", "generated")

	// Generate css file
	cmd := exec.Command("tailwindcss", "-i", "static/tw.css", "-o", "static/main.css", "--minify")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Tailwind failed", string(output), err)
	}
	log.Print("generated tailwind css files")

	r.Get("/", route.Index)
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
