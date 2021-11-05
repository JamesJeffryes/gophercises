package link

import (
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseFile(filePath string) []Link {
	r, err := os.Open(filePath)
	quitOnError(err)
	return ParseHTML(r)
}

func ParseHTML(r io.Reader) []Link {
	doc, err := html.Parse(r)
	quitOnError(err)

	links := make([]Link, 0)
	var walkNodes func(*html.Node)
	walkNodes = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, Link{attr.Val, grabText(n)})
				}
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walkNodes(c)
			}
		}
	}
	walkNodes(doc)
	return links
}

func grabText(n *html.Node) string {
	var sb strings.Builder
	var rec func(*html.Node)
	rec = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rec(c)
		}
	}
	rec(n)

	return strings.Join(strings.Fields(sb.String()), " ")
}

func quitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
