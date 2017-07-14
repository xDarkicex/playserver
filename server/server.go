package server

import (
	"fmt"
	"net/http"
	"regexp"

	"log"

	"github.com/mgutz/ansi"
	"github.com/pkg/errors"
)

type Params struct {
	Response http.ResponseWriter
	Request  *http.Request
}

// Route type for adding to the server
type Route struct {
	Regexp   *regexp.Regexp
	HasRegex bool
	Path     string
	Handler  func(*Params)
}

// Server struct for fancy methods
type Server struct {
	Mux    *http.ServeMux
	Routes []*Route
}

func (s *Server) AddRoute(route *Route) *Server {
	if route.HasRegex {
		var err error
		route.Regexp, err = regexp.Compile(route.Path)
		if err != nil {
			fmt.Println(errors.WithStack(err))
		}
	}
	s.Routes = append(s.Routes, route)
	return s
}
func (s *Server) Serve(port string) {
	lime := ansi.ColorFunc("green+h:black")
	red := ansi.ColorFunc("red+h:black")
	s.Mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		for _, route := range s.Routes {
			if route.HasRegex && route.Regexp.MatchString(request.URL.Path) || request.URL.Path == route.Path {
				route.Handler(&Params{
					Response: response,
					Request:  request,
				})
				log.Println(lime("200"), ":", "127.0.0.1", request.URL.String())
				return
			}
		}
		log.Println(red("404"), ":", "127.0.0.1", request.URL.String())
		http.NotFound(response, request)
	})
	log.Fatal(http.ListenAndServe(port, s.Mux))
}

// New will load all server and routes
func New() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}
