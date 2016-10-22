package main

import (
  "net/http"
  "html/template"
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
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
    Name   string
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

var upgrader = websocket.Upgrader{}

type Hub struct {
  clients      map[*Client]bool
  broadcast    chan []byte
  addClient    chan *Client
  removeClient chan *Client
}

var hub = Hub {
  clients:      make(map[*Client]bool),
  broadcast:    make(chan []byte),
  addClient:    make(chan *Client),
  removeClient: make(chan *Client),
}

func (hub *Hub) start() {
  for {
    select {
    case client := <-hub.addClient:
      hub.clients[client] = true
    case client := <-hub.removeClient:
      if _, ok := hub.clients[client]; ok {
        delete(hub.clients, client)
        close(client.send)
      }
    case msg := <-hub.broadcast:
      for client := range hub.clients {
        client.send <- msg
      }
    }
  }
}

type Client struct {
  ws *websocket.Conn
  send chan []byte
}

func (c *Client) write() {
  defer func() {
    hub.removeClient <- c
    c.ws.Close()
  }()

  for {
    select {
    case msg, ok := <-c.send:
      if !ok {
        c.ws.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      c.ws.WriteMessage(websocket.TextMessage, msg)
    }
  }
}

func (c *Client) read() {
  defer func() {
    hub.removeClient <- c
    c.ws.Close()
  }()

  for {
    _, msg, err := c.ws.ReadMessage()
    if err != nil {
      return
    }

    hub.broadcast <- msg
  }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    http.NotFound(w, r)
  }

  client := &Client {
    ws: conn,
    send: make(chan []byte),
  }

  hub.addClient <- client

  go client.write()
  go client.read()
}

func main() { http.HandleFunc("/", handler)
  go hub.start()

  r := mux.NewRouter()
  r.HandleFunc("/", handler)
  r.HandleFunc("/users", getUsersHandler).Methods("GET")
  r.HandleFunc("/users/new/{username:[a-zA-Z]+}", newUsersHandler).Methods("GET")
  r.HandleFunc("/ws", wsHandler)
  http.ListenAndServe(":8080", r)
}

