package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if newURL, ok := pathsToUrls[req.URL.Path]; ok {
			http.Redirect(w, req, newURL, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, req)
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type ymlRedirect []struct {
	Path string
	URL  string
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var yamlData ymlRedirect
	err := yaml.Unmarshal(yml, &yamlData)
	if err != nil {
		panic(err)
	}
	pathMap := make(map[string]string, len(yamlData))
	for _, data := range yamlData {
		pathMap[data.Path] = data.URL
	}

	return MapHandler(pathMap, fallback), nil
}
