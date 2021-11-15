package main

import (
	"flag"
	"fmt"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

var debugMode bool

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Run server in Debug mode.")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", showSource)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3001", recoveryMiddleware(mux, debugMode)))
}

func recoveryMiddleware(app http.Handler, debugMode bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				log.Error(stack)
				if debugMode {
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(writer, "<h1>Error!</h1>")
					stack = strings.ReplaceAll(stack, "\n", "<br>\n")
					pattern := regexp.MustCompile("((\\S+):(\\d+))")
					stack = pattern.ReplaceAllString(stack, "<a href=\"/debug/$2?lines=$3\">$1</a>")
					stack = fmt.Sprintf("<p>%s</p>", stack)
					fmt.Fprint(writer, stack)

				} else {
					http.Error(writer, "Server Error", http.StatusInternalServerError)
				}

			}
		}()
		app.ServeHTTP(writer, request)
	}
}

func showSource(w http.ResponseWriter, r *http.Request) {
	formatterArgs := []html.Option{
		html.Standalone(true),
		html.WithLineNumbers(true),
	}
	filePath := strings.Replace(r.URL.Path, "/debug", "", 1)
	source, err := os.ReadFile(filePath)
	if logErr(err, w) {
		return
	}

	// Line highlighting
	lineTxts, ok := r.URL.Query()["lines"]
	if ok {
		fmtdLines := make([][2]int, len(lineTxts))
		for i, txt := range lineTxts {
			lns := strings.SplitN(strings.ReplaceAll(txt, "L", ""), "-", 2)
			ln, err := strconv.Atoi(lns[0])
			if logErr(err, w) {
				return
			}
			if len(lns) == 1 {
				fmtdLines[i] = [2]int{ln, ln}
			} else {
				ln2, err := strconv.Atoi(lns[1])
				if logErr(err, w) {
					return
				}
				fmtdLines[i] = [2]int{ln, ln2}

			}
		}
		formatterArgs = append(formatterArgs, html.HighlightLines(fmtdLines))
	}

	lex := lexers.Get("go")
	style := styles.Get("monokai")
	formatter := html.New(formatterArgs...)
	iterator, _ := lex.Tokenise(nil, string(source))
	logErr(formatter.Format(w, style, iterator), w)

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

func logErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Error(err)
		http.Error(w, "Invalid Source File", http.StatusBadRequest)
		return true
	}
	return false
}
