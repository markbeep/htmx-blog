package route

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/httplog"
	"github.com/markbeep/htmx-blog/internal/config"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type Post struct {
	Title         string   `yaml:"title"`
	RawDate       string   `yaml:"date"`
	Tags          []string `yaml:"tags"`
	ShowDate      bool     `yaml:"showDate"`
	Justify       bool     `yaml:"justify"`
	Draft         bool     `yaml:"draft"`
	Mathjax       bool     `yaml:"mathjax"`
	FormattedDate string
	ActualDate    time.Time
	XMLDate       string
	Path          string
	Words         int
	Description   string
	Buffer        []byte
	HtmlPath      string
}

type PostsHandler struct {
	posts       map[string]*Post
	sortedPosts []*Post
}

func (ph *PostsHandler) sortPosts() {
	ph.sortedPosts = []*Post{}
	for _, v := range ph.posts {
		ph.sortedPosts = append(ph.sortedPosts, v)
	}
	sort.Slice(ph.sortedPosts, func(i, j int) bool {
		return ph.sortedPosts[i].ActualDate.After(ph.sortedPosts[j].ActualDate)
	})
}

func (ph *PostsHandler) GetPosts() []*Post {
	return ph.sortedPosts
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicons/favicon.ico")
}

func Robots(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/robots.txt")
}

func Static(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo.IsDir() {
		config.Logger.Warn(err.Error())
		w.Write([]byte("404"))
		return
	}
	http.ServeFile(w, r, path)
}

func Content(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "..") {
		w.Write([]byte("invalid request"))
		return
	}
	// validate path so we only give out images
	valid := false
	for _, ending := range config.ImageFormats {
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

}

var XMLBuffer []byte

// First time it's called, the new index.xml will be generated with the host path
func XML(posts []*Post) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if XMLBuffer == nil {
			tmpl := template.Must(template.ParseFiles("components/index.xml"))
			var tmpBuffer bytes.Buffer
			tmpl.Execute(&tmpBuffer, map[string]any{
				"Posts":    posts,
				"FullPath": r.Host,
			})
			XMLBuffer = tmpBuffer.Bytes()
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(XMLBuffer)
	}
}

func MiddlewareLogging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" { // don't log /health
			next.ServeHTTP(w, r)
		} else {
			logger := httplog.NewLogger("htmx-blog", httplog.Options{
				LogLevel: "warn",
			})
			httplog.RequestLogger(logger)(next).ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func getMarkdownMetadata(post *Post, ctx *parser.Context) error {
	d := frontmatter.Get(*ctx)
	if err := d.Decode(&post); err != nil {
		return err
	}
	datetime, err := time.Parse("2006-01-02T15:04:05-07:00", post.RawDate)
	if err != nil {
		return err
	}
	post.FormattedDate = datetime.Format("Mon Jan 02, 2006")
	post.XMLDate = datetime.Format("Mon, 02 Jan 2006 15:04:05 -0700")
	post.ActualDate = datetime
	return nil
}

func (ph *PostsHandler) ConvertMarkdown(inPath string) error {
	ph.posts = map[string]*Post{}
	ph.sortedPosts = []*Post{}

	inPath = strings.Trim(inPath, "/")

	// Search all markdown files
	err := filepath.Walk(inPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".md") || err != nil {
			return nil
		}

		post := Post{}
		post.Path = strings.Replace(strings.TrimPrefix(path, inPath), ".md", "", 1)
		if err != nil {
			return err
		}
		ph.posts[post.Path] = &post

		log.Printf("Generated post: %s", post.Path)

		// Turn converted file into a template
		var postHtmlBuffer bytes.Buffer
		ctx := parser.NewContext()
		mdFile, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		config.Markdown.Convert(mdFile, &postHtmlBuffer, parser.WithContext(ctx))
		post.Buffer = postHtmlBuffer.Bytes()

		// parse post metadata
		getMarkdownMetadata(&post, &ctx)
		post.Words = len(strings.Split(string(mdFile), " "))

		return nil
	})
	if err != nil {
		return err
	}

	ph.sortPosts()

	return nil
}
