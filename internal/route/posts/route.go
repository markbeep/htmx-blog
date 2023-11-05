package posts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/markbeep/htmx-blog/internal/config"
	"github.com/markbeep/htmx-blog/internal/route"
)

type post struct {
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
	Buffer      []byte
	HtmlPath    string
}

type PostsHandler struct {
	posts       map[string]*post
	sortedPosts []*post
}

func (ph *PostsHandler) sortPosts() {
	ph.sortedPosts = []*post{}
	for _, v := range ph.posts {
		ph.sortedPosts = append(ph.sortedPosts, v)
	}
	sort.Slice(ph.sortedPosts, func(i, j int) bool {
		return ph.sortedPosts[i].ActualDate.After(ph.sortedPosts[j].ActualDate)
	})
}

func (p *PostsHandler) Posts(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/posts/posts.html"))
		tmpl.Execute(w, map[string]any{
			"Posts": p.sortedPosts,
			"Path":  "posts",
		})
	})

	r.Get("/index.xml", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.xml"))
		output, err := os.Create("generated/index.xml")
		if err != nil {
			config.Logger.Warn(err.Error())
			w.Write([]byte("404"))
			return
		}
		defer output.Close()
		tmpl.Execute(output, map[string]any{
			"Posts":    p.sortedPosts,
			"FullPath": r.Host,
		})
		http.ServeFile(w, r, "generated/index.xml")
	})

	// handles all the blog posts
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		post, ok := p.posts[r.URL.Path]
		if !ok {
			config.Logger.Warn(fmt.Sprintf("Unknown posts path: %s", r.URL.Path))
			route.Error404(w, r)
			return
		}
		w.Write(post.Buffer)
	})
}

// Takes the metadata tags located at the top of a markdown file and returns a Post struct.
// The returned data variable allows for the
func getMarkdownMetadata(path string, data *[]byte) (*post, error) {
	lines := strings.Split(string(*data), "\n")
	if len(lines) == 0 {
		return nil, errors.New("empty markdown file")
	}
	if lines[0] != "---" {
		return nil, errors.New("missing metadata or invalid format")
	}
	post := post{}
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

	// removes all the metadata lines from the top of the markdown file
	removedMetadata := strings.Join(strings.Split(string(*data), "\n")[metadataEndLine+1:], "\n")
	*data = []byte(removedMetadata)

	return &post, nil
}

func (ph *PostsHandler) GenerateHTML(inPath, outPath string, generateFiles bool) error {
	ph.posts = map[string]*post{}
	ph.sortedPosts = []*post{}

	inPath = strings.Trim(inPath, "/")
	outPath = strings.Trim(outPath, "/")

	// Search all markdown files
	err := filepath.Walk(inPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Makes sure all directories exist
		if info.IsDir() {
			path, err := url.JoinPath(outPath, path)
			if err != nil {
				return err
			}
			os.MkdirAll(path, 0700)
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
			return err
		}

		post.Path = strings.TrimPrefix(path, inPath)
		post.Path = strings.Replace(post.Path, ".md", "", 1)
		post.HtmlPath, err = url.JoinPath(outPath, strings.Replace(path, ".md", ".html", 1))
		if err != nil {
			return err
		}
		ph.posts[post.Path] = post

		if generateFiles {
			log.Printf("Generated post: %s", post.Path)

			// Turn converted file into a template
			var postHtmlBuffer bytes.Buffer
			postHtmlBuffer.WriteString(`{{ define "post" }}`)
			config.Markdown.Convert(mdFile, &postHtmlBuffer)
			postHtmlBuffer.WriteString(`{{ end }}`)

			tmpl := template.Must(template.ParseFiles(
				"templates/base.html",
				"templates/bottom-bar.html",
				"templates/posts/post-base.html",
			))
			tmpl = template.Must(tmpl.Parse(postHtmlBuffer.String()))

			var templatedBuffer bytes.Buffer
			tmpl.Execute(&templatedBuffer, map[string]any{
				"Post": post,
			})
			if err = os.WriteFile(post.HtmlPath, templatedBuffer.Bytes(), 0644); err != nil {
				config.Logger.Warn(err.Error())
				return nil
			}
			post.Buffer = templatedBuffer.Bytes()
		} else {
			post.Buffer, err = os.ReadFile(post.HtmlPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	ph.sortPosts()

	return nil
}
