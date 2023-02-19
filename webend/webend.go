package webend

import (
	"net/http"
	"text/template"
)

func webHandler(w http.ResponseWriter, req *http.Request) {
	// parse the index.html and serve
	// it to the client whenever
	// the endpoint is hit
	tmpl, err := template.ParseFiles("statictest/index.html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SetupWebend() {
	// define directories for css and js
	// setup file servers for both
	// so go can see them
	stylesDir := http.Dir("./statictest/css")
	scriptsDir := http.Dir("./statictest/js")
	styles := http.FileServer(stylesDir)
	scripts := http.FileServer(scriptsDir)

	// setup handlers for the routes
	http.HandleFunc("/", webHandler)
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
}
