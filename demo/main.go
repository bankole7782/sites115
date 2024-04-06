package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/bankole7782/sites115"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/b/{object1}", blogHandler)
	r.HandleFunc("/b/{object1}/{object2}", blogHandler)
	r.HandleFunc("/search", searchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	s1o, err := sites115.Init("markdowns_md.tar.gz", "markdowns_idx.tar.gz")
	if err != nil {
		panic(err)
	}

	allPaths, err := s1o.ReadAllMD()
	if err != nil {
		panic(err)
	}

	pathsToTitle := make(map[string]string)
	for _, p := range allPaths {
		pathsToTitle[p], _ = s1o.ReadMDTitle(p)
	}

	type Context struct {
		PathsToTitle map[string]string
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, Context{pathsToTitle})
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	object1 := vars["object1"]
	object2 := vars["object2"]

	toFind := object1
	if object2 != "" {
		toFind += "/" + object2
	}

	s1o, err := sites115.Init("markdowns_md.tar.gz", "markdowns_idx.tar.gz")
	if err != nil {
		panic(err)
	}

	htmlStr, err := s1o.ReadMDAsHTML(toFind)
	if err != nil {
		panic(err)
	}

	mdTitle, _ := s1o.ReadMDTitle(toFind)

	tHTML := template.HTML(htmlStr)

	type Context struct {
		HTML    template.HTML
		MDTitle string
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/blog_item.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, Context{tHTML, mdTitle})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	s1o, err := sites115.Init("markdowns_md.tar.gz", "markdowns_idx.tar.gz")
	if err != nil {
		panic(err)
	}

	query := r.FormValue("q")
	if query == "" {
		fmt.Fprint(w, "Search by appending '/?q=search+terms' to the address bar")
		return
	}

	results, err := s1o.Search(query)
	if err != nil {
		panic(err)
	}

	pathsToTitle := make(map[string]string)
	for _, p := range results {
		pathsToTitle[p], _ = s1o.ReadMDTitle(p)
	}

	type Context struct {
		Query        string
		PathsToTitle map[string]string
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/search.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, Context{query, pathsToTitle})
}
