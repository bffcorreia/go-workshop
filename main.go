package main

import (
  //"fmt"
  "net/http"
  "html/template"
)

var templates = template.Must(template.ParseGlob("templates/*.tmpl"))

func handler(w http.ResponseWriter, r *http.Request) {
  templates.ExecuteTemplate(w, "index", nil)
}

func main() { http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}

