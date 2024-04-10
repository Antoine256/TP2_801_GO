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

func FireDetection(ts *Space, batiments []BatimentType, chanceToFire int) { // Pour la simulation
	// Générer un indice aléatoire pour sélectionner un bâtiment
	randomIndex := rand.Intn(len(batiments))
	batiment := batiments[randomIndex]

	var alarmOn bool

	batimentTuple, err := ts.GetP("alarm", batiment.ID, &alarmOn)
	if err == nil {
		alarmOn = (batimentTuple.GetFieldAt(2)).(bool)
	} else {
		alarmOn = false
	}

	println(batiment.ID, alarmOn)

	if rand.Int()*100 > chanceToFire {
		SendToConn("{\"message\": \"Alarm on\",\"idBatiment\": " + strconv.Itoa(batiment.ID) + "}")
		ts.Put("alarm", batiment.ID, true)
	}
	if alarmOn && rand.Int()*100 < chanceToFire {
		SendToConn("{\"message\": \"Alarm off\",\"idBatiment\": " + strconv.Itoa(batiment.ID) + "}")
		ts.Put("alarm", batiment.ID, false)
	}
	time.Sleep(5 * time.Second)
	FireDetection(ts, batiments, chanceToFire)
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

	wg.Add(1)
	go FireDetection(&Ts, batiments, 50)
	wg.Done()
	wg.Wait()
}
