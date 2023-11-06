/**
* Starts server
 */
package main

import (
	"log"
	"net/http"
	"server/randomness"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SessionCodeMessage struct {
	Code string `json:"code"`
}

var upgrader = websocket.Upgrader{}

func generateSessionCode() (string, error) {
	token, err := randomness.GenerateRandomString(6)
	return token, err
}

func sendSessionCode() {

}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fs := http.FileServer(http.Dir("./media"))
	http.Handle("/media/", http.StripPrefix("/media/", fs))

	log.Println("Server is running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
