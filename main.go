package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
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
	"github.com/yuin/goldmark"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

var (
	port = flag.String("port", os.Getenv("PORT"), "port to host the website at")
)

var imageFormats = []string{"jpg", "jpeg", "gif", "png", "gif"}

type Post struct {
	Title       string
	Date        string
	XMLDate     string
	ActualDate  time.Time
	Path        string
	Tags        []string
	ShowDate    bool
	Justify     bool
	Draft       bool
	Mathjax     bool
	Words       int
	Description string
	buffer      bytes.Buffer // TODO: doesn't seem to be reading this in another thread
	htmlPath    string
}

// Takes the metadata tags located at the top of a markdown file and returns a Post struct
func getMarkdownMetadata(path string, data *[]byte) (*Post, error) {
	lines := strings.Split(string(*data), "\n")
	if len(lines) == 0 {
		return nil, errors.New("empty markdown file")
	}
	if lines[0] != "---" {
		return nil, errors.New("missing metadata or invalid format")
	}
	post := Post{}
	post.Words = len(strings.Split(string(*data), " "))
	post.Description = "" // TODO: add a way to get a post description
	metadataEndLine := 0
	for i, l := range lines {
		if i == 0 {
			continue
		}
		if l == "---" {
			metadataEndLine = i
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
			post.XMLDate = datetime.Format("Mon, 02 Jan 2006 15:04:05 -0700")
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

	removedMetadata := strings.Join(strings.Split(string(*data), "\n")[metadataEndLine+1:], "\n")
	*data = []byte(removedMetadata)

	return &post, nil
}

func main() {
	flag.Parse()
	if *port == "" {
		*port = "3000"
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("github-dark"),
				highlighting.WithFormatOptions(
					chromahtml.ClassPrefix("codeblock"),
					chromahtml.WithCustomCSS(map[chroma.TokenType]string{
						chroma.PreWrapper: "padding: 10px; margin: 20px 0 20px 0; border-radius: 10px; box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3);",
					}),
				),
			),
		),
	)

	r := chi.NewRouter()
	logger := httplog.NewLogger("htmx-blog", httplog.Options{
		LogLevel: slog.LevelDebug,
	})
	r.Use(httplog.RequestLogger(logger))

	posts := map[string]Post{}
	// Find all generated files
	err := filepath.Walk("content/posts", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			os.MkdirAll("generated/"+path, 0700)
			return nil
		}
		if !strings.HasSuffix(path, ".md") || err != nil {
			return nil
		}

		mdFile, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		post, err := getMarkdownMetadata(path, &mdFile)
		if err != nil {
			logger.Warn(err.Error())
			return nil
		}

		post.Path = strings.TrimPrefix(path, "content/")
		post.Path = "/" + strings.Replace(post.Path, ".md", "", 1)
		post.htmlPath = "generated/" + strings.Replace(path, ".md", ".html", 1)
		posts[post.Path] = *post
		log.Printf("added post: %s", post.Path)

		var postHtmlBuffer bytes.Buffer
		postHtmlBuffer.WriteString(`{{ define "post" }}`)
		markdown.Convert(mdFile, &postHtmlBuffer)
		postHtmlBuffer.WriteString(`{{ end }}`)

		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/posts/post-base.html"))
		tmpl = template.Must(tmpl.Parse(postHtmlBuffer.String()))
		log.Print("templates", tmpl.DefinedTemplates())

		var templatedBuffer bytes.Buffer
		tmpl.Execute(&templatedBuffer, map[string]any{
			"Post": post,
		})
		if err = os.WriteFile(post.htmlPath, templatedBuffer.Bytes(), 0644); err != nil {
			logger.Warn(err.Error())
			return nil
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Generate new css file
	cmd := exec.Command("tailwindcss", "-i", "tw.css", "-o", "static/main.css", "--minify")
	if err = cmd.Run(); err != nil {
		log.Fatalln("Tailwind failed", err)
	}
	if err = cmd.Wait(); err != nil {
		log.Fatalln("Tailwind failed", err)
	}

	sortedPosts := []Post{}
	for _, p := range posts {
		sortedPosts = append(sortedPosts, p)
	}
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].ActualDate.After(sortedPosts[j].ActualDate)
	})

	// static
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/m.ico")
	})
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		fileInfo, err := os.Stat(path)
		if err != nil || fileInfo.IsDir() {
			w.Write([]byte("404"))
			return
		}
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
			tmpl.Execute(w, map[string]any{
				"Posts": sortedPosts,
				"Path":  "posts",
			})
		})
		r.Get("/index.xml", func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.Must(template.ParseFiles("templates/index.xml"))
			output, err := os.Create("generated/index.xml")
			if err != nil {
				logger.Warn(err.Error())
				w.Write([]byte("404"))
				return
			}
			defer output.Close()
			tmpl.Execute(output, map[string]any{
				"Posts":    sortedPosts,
				"FullPath": r.Host,
			})
			http.ServeFile(w, r, "generated/index.xml")
		})
		// handles all the blog posts
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			post, ok := posts[r.URL.Path]
			if !ok {
				w.Write([]byte("404"))
				return
			}
			http.ServeFile(w, r, post.htmlPath)
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
