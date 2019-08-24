package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("./templates/edit.html", "./templates/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	log.Printf("Loading %s\n", filename)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, temp string, p *Page) {
	log.Printf("Rendering template %s", temp)
	err := templates.ExecuteTemplate(w, temp, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		log.Printf("Warning, not found! Redirecting to edit page for %s.txt\n", title)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	log.Printf("Viewing %s.txt\n", title)
	renderTemplate(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Editing %s.txt\n", title)
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("Trouble loading page :(")}
	}
	renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Saving %s", title)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	port := ":8080"
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Printf("Server started! Listening on %s", port)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
