package main

import (
	"loficode/templates"
	"net/http"

	"github.com/a-h/templ"
)

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))

	http.Handle("/", templ.Handler(templates.HomePage()))

	http.ListenAndServe(":8080", nil)
}
