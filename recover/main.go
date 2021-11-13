package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

var debugMode bool

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Run server in Debug mode.")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3001", recoveryMiddleware(mux, debugMode)))
}

func recoveryMiddleware(app http.Handler, debugMode bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				log.Error(stack)
				if debugMode {
					http.Error(writer, "Server Error", http.StatusInternalServerError)
				} else {
					http.Error(writer, string(stack), http.StatusInternalServerError)
				}

			}
		}()
		app.ServeHTTP(writer, request)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
