package global

import (
	. "github.com/pspaces/gospace/space"
	"log"
	"strconv"
)

func updatePassTime(ts *Space, batiment int, porte int) {
	var x int
	tuple, err := ts.GetP("passTime", batiment, porte, &x)
	if err == nil {
		x = (tuple.GetFieldAt(3)).(int)
	} else {
		x = 1
	}
	ts.Put("passTime", batiment, porte, x)
}

func sensor(ts *Space, batiment int, porte int) bool {
	var passTime int

	_, errLaser := ts.QueryP("Laser detected", batiment, porte)
	_, errdoorClose := ts.QueryP("DoorClose", batiment, porte)
	tuplePassTime, errPassTime := ts.QueryP("passTime", batiment, porte, &passTime)

	if errPassTime == nil {
		passTime = (tuplePassTime.GetFieldAt(3)).(int)
	} else {
		passTime = 0
	}

	if errLaser == nil {
		passTime += 1
		if passTime > 1 {
			log.Printf("Déclenche Alarme")
			SendToConn("{\"message\": \"Alarme Door\",\"idPorte\": " + strconv.Itoa(porte) + ",\"idBatiment\": " + strconv.Itoa(batiment) + "}")
		}
		updatePassTime(ts, batiment, porte)
		ts.Get("Laser detected", batiment, porte)
	}

	if errdoorClose == nil {
		ts.Get("DoorClose", batiment, porte)
		var res bool = passTime >= 1
		ts.Put("Human Pass", batiment, porte, res)
		ts.GetP("passTime", batiment, porte, &passTime)
		return res
	}

	return sensor(ts, batiment, porte)
}
