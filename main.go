package main

import (
  //"fmt"
  "net/http"
  "html/template"
  "github.com/gorilla/mux"
  "encoding/json"
)

type User struct {
  Username string `json:"user"`
}

var users = []User{User{"admin"}}
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

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application.json")

  j, _ := json.Marshal(users)
  w.Write(j)
}

func newUsersHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application.json")

  username := mux.Vars(r)["username"]
  users = append(users, User{ Username: username })

  j, _ := json.Marshal(users)
  w.Write(j)
}

func main() { http.HandleFunc("/", handler)
  r := mux.NewRouter()
  r.HandleFunc("/", handler)
  r.HandleFunc("/users", getUsersHandler).Methods("GET")
  r.HandleFunc("/users/new/{username:[a-zA-Z]+}", newUsersHandler).Methods("GET")
  http.ListenAndServe(":8080", r)
}

