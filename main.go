package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/markbeep/htmx-blog/components"
	"github.com/markbeep/htmx-blog/internal/route"
)

var (
	port     = flag.String("port", os.Getenv("PORT"), "port to host the website at")
	comments = []route.Comment{
		{
			Name:      "Alice Johnson",
			Content:   "This article really opened my eyes to new perspectives. Thanks for sharing!",
			CreatedAt: time.Date(2023, 3, 12, 8, 34, 0, 0, time.UTC),
		},
		{
			Name:      "Marcus Reed",
			Content:   "<h1>I disagree with the points made about renewable energy. I think there's more to consider.</h1>",
			CreatedAt: time.Date(2022, 11, 5, 14, 22, 0, 0, time.UTC),
		},
		{
			Name:      "Jasmine Yu",
			Content:   "Great tutorial! Followed it step by step and got my app running.",
			CreatedAt: time.Date(2023, 5, 20, 19, 47, 0, 0, time.UTC),
		},
		{
			Name:      "Erik Smith",
			Content:   "Could you elaborate on the statistics from the last section? They don't seem to add up.",
			CreatedAt: time.Date(2021, 8, 30, 22, 15, 0, 0, time.UTC),
		},
	}
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
		post.Comments = comments
		r.Get(post.Path, templ.Handler(components.Post(post, t)).ServeHTTP)
		r.Post(post.Path+"/comments", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			comment := route.Comment{CreatedAt: time.Now()}
			comment.Name = strings.TrimSpace(r.PostFormValue("name"))
			comment.Content = strings.TrimSpace(r.PostFormValue("content"))
			if comment.Name == "" || comment.Content == "" {
				// TODO: return comments with error
				log.Fatal("Empty name. Unimplemented")
			}
			comments = append(comments, comment)
			post.Comments = comments
			templ.Handler(components.Comments(post.Comments)).ServeHTTP(w, r)
		})
	}

	host := fmt.Sprintf(":%s", *port)
	log.Printf("listening on %s", host)
	log.Fatal(http.ListenAndServe(host, r))
}
