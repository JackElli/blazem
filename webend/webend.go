package webend

import (
	"net/http"
	"text/template"
)

func webHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("statictest/index.html")

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func SetupWebend() {
	//web stuff
	styles := http.FileServer(http.Dir("./statictest/css"))
	scripts := http.FileServer(http.Dir("./statictest/js"))
	http.HandleFunc("/", webHandler)
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
}
