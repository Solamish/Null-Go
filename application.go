package nullgo

import (
	"fmt"
	"log"
	"net/http"
)

type App struct {
	router Router
}

func Default() (app *App) {
	app = &App{Router{}}
	app.router.Init()
	return app
}

func (app *App) Run(addr string) {
	http.Handle("/", app)
	Debug("Listening and serving HTTP on %s\n", addr)
	fmt.Println()
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Port is being occupied\n")
	}

}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.forward(w, r)
	//TODO
}

func (app *App) add(method int, uri string, h HandlerFunc) {
	length := len(uri)
	if uri[length-1] == '/' {
		uri = uri[:length-1]
	}
	app.router.add(method, uri, h)
}

func (app *App) Websocket(uri string, config WebSocketConfig) {
	app.GET(uri, enter)
	app.router.wsMap[uri] = config
}

func (app *App) GET(uri string, h HandlerFunc) {
	app.add(GET, uri, h)
}

func (app *App) POST(uri string, h HandlerFunc) {
	app.add(POST, uri, h)
}

func (app *App) PUT(uri string, h HandlerFunc) {
	app.add(PUT, uri, h)
}

func (app *App) DELETE(uri string, h HandlerFunc) {
	app.add(DELETE, uri, h)
}
