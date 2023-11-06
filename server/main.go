/**
* Starts server
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/randomness"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type ControlsRequest struct {
	Code string `json:"code"`
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SessionCodeMessage struct {
	Code string `json:"code"`
}

var upgrader = websocket.Upgrader{}
var connections map[string]*websocket.Conn = make(map[string]*websocket.Conn)

func generateSessionCode() (string, error) {
	token, err := randomness.GenerateRandomString(6)
	return token, err
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Println("Connection!")

	code, err := generateSessionCode()
	if err != nil {
		log.Printf("[ERROR] failed to generate code: %v", err)
		return
	}
	log.Println(code)
	connections[code] = conn
	defer func() {
		delete(connections, code)
	}()

	err = conn.WriteJSON(Message{Type: "session_code", Data: SessionCodeMessage{Code: code}})
	if err != nil {
		log.Printf("[ERROR] failed to send message to client")
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Print the message to the console
		log.Printf("%s sent(%d): %s\n", conn.RemoteAddr(), msgType, string(msg))

		// Write message back to browser
		// if err = conn.WriteMessage(msgType, msg); err != nil {
		// 	log.Println(err)
		// 	break
		// }
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func findConnection(r *http.Request) (*websocket.Conn, error) {
	var req ControlsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	log.Printf("ControlsRequest: %v", req)

	conn, ok := connections[req.Code]
	if !ok {
		return nil, fmt.Errorf("session for '%s' not found", req.Code)
	}

	return conn, nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.HandleFunc("/ws", wsHandler)

	FileServer(r, "/media", http.Dir("./media"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})

	r.Post("/api/next", func(w http.ResponseWriter, r *http.Request) {
		conn, err := findConnection(r)
		if err != nil {
			log.Printf("[ERROR] session not found %v", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		msg := Message{Type: "next"}
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Printf("[ERROR] failed to send message '%v' to client: %v", msg, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})
	r.Post("/api/prev", func(w http.ResponseWriter, r *http.Request) {
		conn, err := findConnection(r)
		if err != nil {
			log.Printf("[ERROR] session not found %v", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		msg := Message{Type: "prev"}
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Printf("[ERROR] failed to send message '%v' to client: %v", msg, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})

	r.Post("/api/start", func(w http.ResponseWriter, r *http.Request) {
		conn, err := findConnection(r)
		if err != nil {
			log.Printf("[ERROR] session not found %v", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		err = conn.WriteJSON(Message{Type: "start"})
		if err != nil {
			log.Printf("[ERROR] failed to send message to client: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("{}"))
	})

	log.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
