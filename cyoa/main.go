package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type storyStruct map[string]struct {
	Title   string
	Story   []string
	Options []struct {
		Text string
		Arc  string
	}
}

func main() {
	story := importStoryJSON("gopher.json")
	storyHandler := makeHandler(story)
	fmt.Println("Starting the server on :8080")
	err := http.ListenAndServe(":8080", storyHandler)
	if err != nil {
		panic(err)
	}
}

func makeHandler(story storyStruct) http.Handler {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, found := story[strings.Trim(r.URL.Path, "/")]
		if !found {
			fmt.Fprintln(w, "404")
		}
		tmpl.Execute(w, page)
	})
}

func importStoryJSON(s string) storyStruct {
	r, err := ioutil.ReadFile(s)
	if err != nil {
		panic(err)
	}
	var story storyStruct
	err = json.Unmarshal(r, &story)
	if err != nil {
		panic(err)
	}
	return story
}
