package route

import (
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/markbeep/htmx-blog/internal/config"
)

func Index(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		config.Logger.Warn(err.Error())
	}
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/m.ico")
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

func About(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/about.html"))
	if err := tmpl.Execute(w, map[string]any{
		"Path": "about",
	}); err != nil {
		config.Logger.Warn(err.Error())
	}
}

func Polyring(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/bottom-bar.html", "templates/polyring.html"))
	if err := tmpl.Execute(w, map[string]any{
		"Path": "polyring",
	}); err != nil {
		config.Logger.Warn(err.Error())
	}
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