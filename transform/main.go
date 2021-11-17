package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const UPLOAD_DIR = "./imgs"

var shapes = []string{"Combo", "Triangles", "Rectangles", "Ellipses", "Circles",
	"Rotated Rectangles", "Beziers", "Rotated Ellipses", "Polygons"}

var (
	err  error
	tmpl *template.Template
)

func main() {
	tmpl, err = template.ParseFiles("main.gohtml")
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/imgs/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/select/", paramSelect)
	mux.HandleFunc("/", upload)
	log.Fatal(http.ListenAndServe(":3002", mux))
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err = tmpl.Execute(w, nil)
		logErr(err, w)
	} else {
		err = r.ParseMultipartForm(10 << 20)
		if logErr(err, w) {
			return
		}
		file, handler, err := r.FormFile("imgFile")
		if logErr(err, w) {
			return
		}
		defer file.Close()

		log.Printf("Uploaded File: %+v\n", handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)

		dir, err := os.MkdirTemp(UPLOAD_DIR, "upload")
		fPath := fmt.Sprintf("%s/%s", dir, handler.Filename)
		log.Println(fPath)
		f, err := os.Create(fPath)
		if logErr(err, w) {
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if logErr(err, w) {
			return
		}
		err = tmpl.Execute(w, htmlContext{Options: []imgOption{{Url: strings.Trim(fPath, "./")}}})
		logErr(err, w)

	}

}

func logErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Error(err)
		http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
		return true
	}
	return false
}

type imgParm struct {
	name      string
	shape     string
	imgDir    string
	numShapes string
	imgPath   string
}

type imgOption struct {
	Url  string
	Name string
}

type htmlContext struct {
	Options []imgOption
	Step    string
}

func paramSelect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	err := r.ParseForm()
	if logErr(err, w) {
		return
	}
	log.Print()
	imgPath := r.FormValue("img")
	log.Print(imgPath)
	if imgPath == "" {
		http.Redirect(w, r, "/", 301)
		return
	}
	shape := r.FormValue("shape")
	numShapes := r.FormValue("number")
	imgDir := path.Dir(imgPath)

	if shape == "" {
		numShapes = "50"
		imgParams := make(chan imgParm, len(shapes))
		imgOutput := make(chan imgOption, len(shapes))
		var options []imgOption
		for _, s := range shapes {
			imgParams <- imgParm{s, s, imgDir, numShapes, imgPath}
			go makePrimitive(w, imgParams, imgOutput)
		}
		for i := range shapes {
			options = append(options, <-imgOutput)
			log.Printf("%d of %d complete", i+1, len(shapes))
		}
		err = tmpl.Execute(w, htmlContext{options[1:], "shape"})
		logErr(err, w)
	} else if numShapes == "" {
		numOptions := []string{"50", "100", "150", "200"}
		imgParams := make(chan imgParm, len(numOptions))
		imgOutput := make(chan imgOption, len(numOptions))
		var options []imgOption
		for _, num := range numOptions {
			imgParams <- imgParm{num, shape, imgDir, num, imgPath}
			go makePrimitive(w, imgParams, imgOutput)
		}
		for i := range numOptions {
			options = append(options, <-imgOutput)
			sort.Slice(options, func(i, j int) bool {
				I, _ := strconv.Atoi(options[i].Name)
				J, _ := strconv.Atoi(options[j].Name)
				return I < J
			})
			log.Printf("%d of %d complete", i+1, len(numOptions))
		}
		err = tmpl.Execute(w, htmlContext{options, "number"})
		logErr(err, w)

	} else {
		name := fmt.Sprintf("%s %s", numShapes, shape)
		url := fmt.Sprintf("/%s/%s.jpeg", imgDir, name)
		err = tmpl.Execute(w, htmlContext{[]imgOption{{url, name}}, "selected"})
		logErr(err, w)
	}
	log.Print(time.Since(start))
}

func makePrimitive(w http.ResponseWriter, paramChan <-chan imgParm, outputChan chan<- imgOption) {
	params := <-paramChan
	shapeIdx := strconv.Itoa(indexOf(params.shape, shapes))
	outPath := fmt.Sprintf("%s/%s %s.jpeg", params.imgDir, params.numShapes, params.shape)
	// Skip existing files
	if _, err := os.Stat(outPath); err == nil {
		outputChan <- imgOption{"/" + outPath, params.name}
		return
	}
	cmd := exec.Command("primitive", "-i", params.imgPath, "-o", outPath,
		"-m", shapeIdx, "-n", params.numShapes)
	text, err := cmd.CombinedOutput()
	if logErr(err, w) {
		log.Error(string(text))
	} else {
		outputChan <- imgOption{"/" + outPath, params.name}
	}
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
