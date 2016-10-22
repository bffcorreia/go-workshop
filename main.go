package main

import (
  //"fmt"
  "net/http"
  "html/template"
)

var templates = template.Must(template.ParseGlob("templates/*.tmpl"))

func handler(w http.ResponseWriter, r *http.Request) {
  type Index struct {
    Name string
    name string
  }

  i := Index {
    Name: "bffcorreia", // public
    name: "gdg",        // private
  }

  templates.ExecuteTemplate(w, "index", &i)
}

func main() { http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}

