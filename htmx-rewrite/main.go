package main

import (
	"loficode/templates/pages/home"
	"loficode/templates/pages/notfound"
	"loficode/templates/pages/privacypolicy"
	"loficode/templates/pages/tos"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			notfound.NotFound().Render(r.Context(), w)
			return
		}
		home.Home().Render(r.Context(), w)
	})
	http.Handle("/tos", templ.Handler(tos.TermsOfService()))
	http.Handle("/privacy-policy", templ.Handler(privacypolicy.PrivacyPolicy()))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
