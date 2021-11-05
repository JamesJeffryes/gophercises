package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/jamesjeffryes/gophercises/sitemap/sitemap"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please supply a URL to parse")
	}
	depth := *flag.Int("d", 3, "Max search depth")
	seen := make(map[string]bool, 0)
	site := sitemap.Sitemap{}

	site = sitemap.WalkPage(os.Args[1], seen, site, depth)
	writer := io.Writer(os.Stdout)
	_ = sitemap.WriteSiteMap(&writer, site)

}
