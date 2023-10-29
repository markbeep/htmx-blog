package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

var (
	port = flag.String("port", os.Getenv("PORT"), "port to host the website at")
)

var imageFormats = []string{"jpg", "jpeg", "gif", "png", "gif"}

type Post struct {
	Title      string
	Date       string
	ActualDate time.Time
	Path       string
	HTMLPath   string
	Tags       []string
	ShowDate   bool
	Justify    bool
	Draft      bool
	Mathjax    bool
	Words      int
}

// Takes the metadata tags located at the top of a markdown file and returns a Post struct
func getMarkdownMetadata(path string) (*Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return nil, errors.New("empty markdown file")
	}
	if lines[0] != "---" {
		return nil, errors.New("missing metadata or invalid format")
	}
	post := Post{}
	post.Words = len(strings.Split(string(data), " "))
	for i, l := range lines {
		if i == 0 {
			continue
		}
		if l == "---" {
			break
		}
		lineValues := strings.SplitN(l, ": ", 2)
		if len(lineValues) < 2 {
			continue
		}
		key := lineValues[0]
		value := lineValues[1]
		switch key {
		case "title":
			post.Title = strings.Trim(value, "\"")
		case "date":
			datetime, err := time.Parse("2006-01-02T15:04:05-07:00", value)
			if err != nil {
				post.Date = ""
				continue
			}
			post.Date = datetime.Format("Mon Jan 02, 2006")
			post.ActualDate = datetime
		case "tags":
			tags := []string{}
			json.Unmarshal([]byte(value), &tags)
			post.Tags = tags
		case "showDate":
			post.ShowDate = value == "true"
		case "justify":
			post.Justify = value == "true"
		case "draft":
			post.Draft = value == "true"
		case "mathjax":
			post.Mathjax = value == "true"
		}
	}
	return &post, nil
}

func main() {
	flag.Parse()
	if *port == "" {
		*port = "3000"
	}

	r := chi.NewRouter()
	logger := httplog.NewLogger("htmx-blog", httplog.Options{
		LogLevel: slog.LevelDebug,
	})
	r.Use(httplog.RequestLogger(logger))

	posts := map[string]Post{}
	// Find all generated files
	err := filepath.Walk("generated/content/posts", func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".html") {
			return nil
		}
		// Get the equivalent markdown file
		markdownPath := strings.TrimLeft(path, "generated/")
		markdownPath = strings.Replace(markdownPath, ".html", ".md", 1)
		post, err := getMarkdownMetadata(markdownPath)
		if err != nil {
			logger.Warn(err.Error())
			return nil
		}
		post.HTMLPath = path
		websitePath := strings.TrimLeft(path, "generated/content")
		websitePath = "/" + strings.Replace(websitePath, ".html", "", 1)
		post.Path = websitePath
		posts[websitePath] = *post
		log.Printf("added post: %s", websitePath)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// static
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/m.ico")
	})
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		http.ServeFile(w, r, path)
	})

	// index page
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			logger.Warn(err.Error())
		}
	})
	// about me
	r.Get("/about", func(w http.ResponseWriter, _ *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/about.html"))
		if err := tmpl.Execute(w, map[string]any{
			"Path": "about",
		}); err != nil {
			logger.Warn(err.Error())
		}
	})
	// polyring
	r.Get("/polyring", func(w http.ResponseWriter, _ *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/polyring.html"))
		if err := tmpl.Execute(w, map[string]any{
			"Path": "polyring",
		}); err != nil {
			logger.Warn(err.Error())
		}
	})

	// posts category
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/posts/posts.html"))
			sortedPosts := []Post{}
			for _, p := range posts {
				sortedPosts = append(sortedPosts, p)
			}
			sort.Slice(sortedPosts, func(i, j int) bool {
				return sortedPosts[i].ActualDate.After(sortedPosts[j].ActualDate)
			})
			tmpl.Execute(w, map[string]any{
				"Posts": sortedPosts,
				"Path":  "posts",
			})
		})
		// handles all the blog posts
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			post, ok := posts[r.URL.Path]
			if !ok {
				w.Write([]byte("404"))
				return
			}
			tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/posts/post-base.html", post.HTMLPath))
			tmpl.Execute(w, map[string]any{
				"Post": post,
			})
		})
	})

	// images
	r.Get("/content/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "..") {
			w.Write([]byte("invalid request"))
			return
		}
		// validate path so we only give out images
		valid := false
		for _, ending := range imageFormats {
			if strings.HasSuffix(r.URL.Path, ending) {
				valid = true
				break
			}
		}
		if !valid {
			w.Write([]byte("invalid request"))
			return
		}
		path := strings.Trim(r.URL.Path, "/")
		http.ServeFile(w, r, path)

	})

	host := fmt.Sprintf(":%s", *port)
	log.Printf("listening on %s", host)
	log.Fatal(http.ListenAndServe(host, r))
}
