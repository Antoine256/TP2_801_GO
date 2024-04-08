package global

import (
	"bytes"
	_ "bytes"
	_ "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//On définit ici tout ce qui compose une entrée (une porte, deux reader, l'alarme de la porte, et le sensor)

//on crée une entry à chaque fois que l'on reçois la socket (badge lu)
//l'entry va faire la requete api ("/canAccess") avec les paramètres idbadge et idbatiment
//on lance la fonction qui dirige la porte avec le résultat de la requete (on bloque ou débloque la porte (+ socket pour le SW)) au bout de 30 seconde ferme la porte et écrit porte fermer dans l'espace tuple
//(on lance la fonction qui gere les reader (la couleur) elle envoie une socket avec la couleur, et au bout de 5 secondes pour vert et 10 pour rouge une autre pour eteindre)
//on lance la function sensor qui va nous dire ecrire dans l'espace tuple. entry écoute, et si + de deux recu, alarme immédiate (si alarme, on ecrit dans espace tuple pour fermer la porte et on ferme l'entry et on enregistre l'incident)
//quand la porte est refermée (recu dans l'espace tuple), l'entry est supprimée, et elle est consignée dans le journal de bord (requete api pour l'enregistrer) faire attention a si la personne est rentrée ou pas

func entry(idPorte int, idBatiment int, idBadge int) {
	//params := url.Values{}
	//params.Add("idBadge", strconv.Itoa(idBadge))
	//params.Add("idBatiment", strconv.Itoa(idBatiment))
	//u, _ := url.ParseRequestURI("http://127.0.0.1:8080/api")
	//u.Path = "/canAccess"
	//u.RawQuery = params.Encode()
	//urlStr := fmt.Sprintf("%v", u)
	//resp, err := http.Get(urlStr)
	var resp *http.Response
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err = http.Get("http://127.0.0.1:8080/api/canAccess?idBadge=" + strconv.Itoa(idBadge) + "&idBatiment=" + strconv.Itoa(idBatiment))
	}()
	wg.Wait()
	if err != nil {
		fmt.Println("Error for user entry:", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	res, err := strconv.ParseBool(string(body))
	var WgPorte sync.WaitGroup
	WgPorte.Add(1)
	if res == true {
		go openDoor(idPorte, idBatiment, &WgPorte)
	} else {
		go closeDoor(idPorte, idBatiment, &WgPorte)
	}
	WgPorte.Wait()
	var wg3 sync.WaitGroup
	wg3.Add(1)
	go func() {
		body := []byte(`{
		"idBatiment": ` + strconv.Itoa(idBatiment) + `,
		"idBadge": ` + strconv.Itoa(idBadge) + `,
		"goIn": ` + strconv.FormatBool(res) + `
	}`)
		log.Printf("send event")
		_, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/event/create", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		wg3.Done()
	}()
	wg3.Wait()
	/// attention au sensor si deux requete et requete api pour enregistrer l'event !
}
