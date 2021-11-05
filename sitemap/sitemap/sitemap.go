package sitemap

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/gophercises/sitemap/link"
)

type XMLURL struct {
	Loc string
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []XMLURL `xml:"url"`
}

func WalkPage(url string, seen map[string]bool, sitemap Sitemap, depth int) Sitemap {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	links := link.ParseHTML(resp.Body)
	rootURL := "http://" + resp.Request.URL.Host
	for _, l := range links {
		newUrl := l.Href
		if newUrl[0:1] == "/" {
			newUrl = rootURL + newUrl
		}
		external := !strings.Contains(newUrl, rootURL)
		_, present := seen[newUrl]
		if external || present {
			continue
		}
		seen[newUrl] = true
		sitemap.URLs = append(sitemap.URLs, XMLURL{newUrl})
		if depth > 0 {
			sitemap = WalkPage(newUrl, seen, sitemap, depth-1)
		}
	}
	return sitemap
}

func WriteSiteMap(w *io.Writer, s Sitemap) error {
	en := xml.NewEncoder(*w)
	en.Indent("", " ")
	err := en.Encode(s)
	return err
}
