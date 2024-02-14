package japi

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//go:embed api.js
var japiJsScript string

type Handler func(string) (string, error)

type Api struct {
	Prefix   string
	Handlers map[string]Handler
	JsString string
}

func NewApi(prefix string) *Api {
	return &Api{
		Prefix:   prefix,
		Handlers: make(map[string]Handler),
		JsString: buildJsForPrefix(prefix),
	}
}

func (a *Api) HandleScript(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, a.JsString)
}

func buildJsForPrefix(prefix string) string {
	return strings.ReplaceAll(japiJsScript, "$$PREFIX$$", prefix)
}

func (a *Api) Register(name string, handler Handler) {
	a.Handlers[name] = handler
}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "No api name")
		return
	}
	handler, ok := a.Handlers[name]
	if !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, "No handler")
		return
	}
	bodyDat, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Can't read body: ", err)
		return
	}
	result, err := handler(string(bodyDat))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Api handler error ["+name+"]: ", err)
		return
	}

	// Success response
	w.WriteHeader(200)
	fmt.Fprint(w, result)
}
