package main

import (
  //"fmt"
  "net/http"
  "html/template"
)

var templates = template.Must(template.ParseGlob("templates/*.tmpl"))

func handler(w http.ResponseWriter, r *http.Request) {
  type Header struct {
    Title string
  }

  type Index struct {
    Name string
    Header Header
  }

  i := Index {
    Name: "bffcorreia",
    Header: Header {
      Title: "I am a title",
    },
  }

  templates.ExecuteTemplate(w, "index", &i)
}

func main() { http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}

