package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	us "gopher.com/urlshortner/pkgs"
)

var (
	yamlFile *string
	jsonFile *string
)

func main() {

	parseFlags()

	pathsToUrls := map[string]string{
		"/stack/api": "https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/",
		"/fcc/api":   "https://www.freecodecamp.org/news/rest-api-best-practices-rest-endpoint-design-examples/",
	}

	mux := defaultMux()

	handler := us.MapHandler(pathsToUrls, mux)

	// Yaml Handler
	yamlBytes, err := readFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading file: %s \n", *yamlFile)
	} else {
		handler, err = us.YamlHandler(yamlBytes, handler)
		if err != nil {
			fmt.Println("Error parsing yaml file: ", err)
		}
	}

	// Json Handler
	jsonBytes, err := readFile(jsonFile)
	if err != nil {
		fmt.Printf("Error reading file: %s \n", *jsonFile)
	} else {
		handler, err = us.JsonHandler(jsonBytes, handler)
		if err != nil {
			fmt.Println("Error parsing json file: ", err)
		}
	}

	fmt.Println("Listening on Port :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "No URL mapped to given path '%s'", r.URL.Path)
			return
		}
		fmt.Fprintln(w, "Welcome to the URL shortener")
	})
	return mux
}

func parseFlags() {
	defer flag.Parse()
	yamlFile = flag.String("yaml", "urls.yaml", "Provide the .yaml file that contains list of paths and urls")
	jsonFile = flag.String("json", "urls.json", "Provide the .json file that contains list of paths and urls")
}

func readFile(file *string) ([]byte, error) {
	fileBytes, err := os.ReadFile(*file)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}
