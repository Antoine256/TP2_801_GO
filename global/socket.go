package global

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

import . "github.com/pspaces/gospace"

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
	ts := NewSpace("ts")
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

		if strings.Contains(string(p), "Laser detected") {
			type SM struct {
				Message    string
				IdPorte    int
				IdBatiment int
			}
			var variable SM
			err = json.Unmarshal(p, &variable)
			if err != nil {
				print("error")
			}
			ts.Put(variable.Message, variable.IdBatiment, variable.IdPorte)
		}
		//si on recois un message de type badge lu on lance une entry
		if strings.Contains(string(p), "Badge Lu") {
			type SM struct {
				Message    string
				IsInside   bool
				IdPorte    int
				IdBatiment int
				IdBadge    int
			}

			var variable SM
			err = json.Unmarshal(p, &variable)
			if err != nil {
				print("error")
			}
			//!\ faire la vérification qu'il n'y a pas d'entry en cours pour cette porte !
			go entry(variable.IdPorte, variable.IdBatiment, variable.IdBadge, variable.IsInside, &ts)
		}

		log.Printf("Received message: %s", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}
