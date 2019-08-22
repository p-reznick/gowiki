package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	log.Printf("Saving %s", filename)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	log.Printf("Loading %s\n", filename)
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, temp string, p *Page) {
	log.Printf("Rendering template %s", temp)
	t, _ := template.ParseFiles(temp)
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		log.Printf("Warning, not found! Redirecting to edit page for %s.txt\n", title)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	log.Printf("Viewing %s.txt\n", title)
	renderTemplate(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("Trouble loading page :(")}
	}
	log.Printf("Editing %s.txt\n", title)
	renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	port := ":8080"
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Printf("Server started! Listening on %s", port)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
