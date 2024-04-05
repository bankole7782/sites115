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

	type Context struct {
		Paths []string
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, Context{allPaths})
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

	tHTML := template.HTML(htmlStr)

	type Context struct {
		HTML template.HTML
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/blog_item.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, Context{tHTML})
}
