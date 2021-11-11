package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jamesjeffryes/gophercises/quiet_hn/hn"
)

var storyCache []item

func main() {
	// parse flags
	var port, numStories, workers int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.IntVar(&workers, "workers", 5, "the number stories to load concurrently")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))
	err := updateCache(numStories, 5)
	if err != nil {
		log.Fatalln(err)
	}
	go runCache(numStories, 1*time.Minute, 5)
	http.HandleFunc("/", handler(tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func runCache(numStories int, maxCacheAge time.Duration, workers int) {
	for _ = range time.Tick(maxCacheAge) {
		err := updateCache(numStories, workers)
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func updateCache(numStories int, workers int) error {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return err
	}
	storyCache = GetStories(numStories, ids, client, workers)
	return nil
}

func handler(tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		data := templateData{
			Stories: storyCache,
			Time:    time.Now().Sub(start),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func GetStories(numStories int, ids []int, client hn.Client, workers int) []item {
	toFetch := make(chan orderedID, numStories)
	fetched := make(chan item, numStories)

	for i, id := range ids[:numStories] {
		toFetch <- orderedID{i, id}
	}
	numFetching := numStories
	for w := 0; w < workers; w++ {
		go func() {
			for ord := range toFetch {
				hnItem, err := client.GetItem(ord.id)
				if err != nil {
					continue
				}
				item := parseHNItem(hnItem)
				if isStoryLink(item) {
					item.index = ord.index
					fetched <- item
				} else {
					toFetch <- orderedID{numFetching, ids[numFetching]}
					numFetching++
				}
			}
		}()
	}
	var stories []item
	for i := 0; len(stories) < numStories; i++ {
		stories = append(stories, <-fetched)
	}
	sort.Slice(stories, func(i, j int) bool {
		return stories[i].index < stories[j].index
	})
	log.Printf("Fetched %d items to find %d stories\n", numFetching, numStories)
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

type orderedID struct {
	index int
	id    int
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host  string
	index int
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
