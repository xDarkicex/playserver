package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/xDarkicex/playserver/neuron"
	"github.com/xDarkicex/playserver/render"
	"github.com/xDarkicex/playserver/server"
	"golang.org/x/net/websocket"
)

var httpAddr = flag.String("http", ":8080", "Listen address")
var data = make(map[string]interface{})

func main() {
	flag.Parse()
	s := server.New()
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./public/")))
	s.AddRoute(&server.Route{
		Path:     "^/static/",
		HasRegex: true,
		Handler: func(p *server.Params) {
			fileServer.ServeHTTP(p.Response, p.Request)
		},
	}).AddRoute(&server.Route{
		Path:    "/",
		Handler: index,
	}).AddRoute(&server.Route{
		Path: "/neuron-demo",
		Handler: func(p *server.Params) {
			if err := render.Render(p, "neuron", data); err != nil {
				log.Println(err)
			}
		},
	}).AddRoute(&server.Route{
		Path:    "/api/time",
		Handler: serverTime,
	}).AddRoute(&server.Route{
		Path: "/api/websocket",
		Handler: func(p *server.Params) {
			websocket.Handler(func(ws *websocket.Conn) {
				var msg string

				for {
					websocket.Message.Receive(ws, &msg)
					var data = make(map[string]interface{})
					json.Unmarshal([]byte(msg), &data)
					switch data["api"] {
					case "neuron":
						point := data["data"].(map[string]interface{})
						output := neuron.Ne.Process([]float64{point["x"].(float64), point["y"].(float64)})
						websocket.Message.Send(ws, `{
							"output": "`+strconv.FormatFloat(output, 'f', -1, 64)+`",
							"M": `+strconv.FormatFloat(neuron.M, 'f', -1, 64)+`, 
							"B": `+strconv.FormatFloat(neuron.B, 'f', -1, 64)+`}`)
					default:

						io.Copy(ws, ws)
					}
				}
			}).ServeHTTP(p.Response, p.Request)
		},
	}).Serve(*httpAddr)

}

func index(p *server.Params) {
	if err := render.Render(p, "index", data); err != nil {
		log.Println(err)
	}
}

func serverTime(p *server.Params) {
	t := time.Now()
	p.Response.Header().Set("Content-Type", "application/json")
	io.WriteString(p.Response, `{"time":`+`"`+string(t.Format("3:04"))+`"}`)
}
