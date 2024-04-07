package global

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Autoriser toutes les origines
			return true
		},
	}
	conn *websocket.Conn
)

func SendToConn(message string) {
	if conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Error writing to WebSocket:", err)
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		//si on recois un message de type badge lu on lance une entry
		//!\ faire la v"érification qu'il n'y a pas d'entry en cours pour cette porte !
		
		log.Printf("Received message: %s", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}
