package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"

	"strings"

	"github.com/oxtoacart/bpool"
	"github.com/xDarkicex/playserver/server"
)

var Templates map[string]*template.Template
var bufpool *bpool.BufferPool

func init() {
	bufpool = bpool.NewBufferPool(64)
	Templates = make(map[string]*template.Template)

	layouts, err := filepath.Glob("./layouts/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	includes, err := filepath.Glob("./includes/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		fmt.Println(strings.TrimRight(filepath.Base(layout), ".gohtml"))
		Templates[strings.TrimRight(filepath.Base(layout), ".gohtml")] = template.Must(template.ParseFiles(files...))
	}

}

func Render(p *server.Params, view string, data map[string]interface{}) error {
	data["view"] = view
	if pusher, ok := p.Response.(http.Pusher); ok {
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": p.Request.Header["Accept-Encoding"],
			},
		}
		if err := pusher.Push("./static/assets/stylesheets/application.css", options); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("./static/assets/javascript/application.js", options); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	tmpl, ok := Templates[view]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", view)
	}
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	err := tmpl.ExecuteTemplate(buf, "base.gohtml", data)
	if err != nil {
		return err
	}
	device := p.Request.UserAgent()
	expression := regexp.MustCompile("(Mobi(le|/xyz)|Tablet)")
	if !expression.MatchString(device) {
		p.Response.Header().Set("Connection", "keep-alive")
	}
	p.Response.Header().Set("Vary", "Accept-Encoding")
	p.Response.Header().Set("Cache-Control", "private, max-age=7776000")
	p.Response.Header().Set("Content-Type", "text/html; charset=utf-8")

	buf.WriteTo(p.Response)
	return nil
}
