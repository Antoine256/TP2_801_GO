package global

import (
	"encoding/json"
	. "github.com/pspaces/gospace/space"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func FireDetection(ts *Space, batiment int, chanceToFire int) {
	var alarmOn bool

	batimentTuple, err := ts.GetP("alarm", batiment, &alarmOn)
	if err == nil {
		alarmOn = (batimentTuple.GetFieldAt(2)).(bool)
	} else {
		alarmOn = false
	}

	if rand.Int()*100 > chanceToFire { // On simule un incendie
		SendToConn("{\"message\": \"Alarm on\",\"idBatiment\": " + strconv.Itoa(batiment) + "}")
		ts.Put("alarm", batiment, true)
	}
	if alarmOn && rand.Int()*100 < chanceToFire { // On simule la fin d'un incendie
		SendToConn("{\"message\": \"Alarm off\",\"idBatiment\": " + strconv.Itoa(batiment) + "}")
		ts.Put("alarm", batiment, false)
	}
	time.Sleep(5 * time.Second)
	FireDetection(ts, batiment, chanceToFire)
}

func AddDetectorOnAllBatiments() {
	var resp *http.Response
	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err = http.Get("http://127.0.0.1:8080/api/batiment")
	}()
	wg.Wait()
	if err != nil {
		log.Println("Error for get batiments", err)
		return
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	// Désérialiser la réponse JSON dans une slice de Batiment
	var batiments []BatimentType
	err = json.Unmarshal(body, &batiments)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	wg.Add(len(batiments))
	for _, batiment := range batiments {
		log.Print("Processus d'alarme lancé pour le batiment : ", batiment.ID, " ", batiment.Name, "\n")
		go FireDetection(&Ts, batiment.ID, 50)
		wg.Done()
	}
	wg.Wait()
}
